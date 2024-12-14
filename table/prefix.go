package table

import (
	"sync"
	"time"

	"github.com/pulsejet/go-ndn-dv/config"
	"github.com/pulsejet/go-ndn-dv/tlv"
	enc "github.com/zjkmxy/go-ndn/pkg/encoding"
	ndn_sync "github.com/zjkmxy/go-ndn/pkg/engine/sync"
	"github.com/zjkmxy/go-ndn/pkg/log"
	"github.com/zjkmxy/go-ndn/pkg/ndn"
	"github.com/zjkmxy/go-ndn/pkg/security"
	"github.com/zjkmxy/go-ndn/pkg/utils"
)

type PrefixTable struct {
	config *config.Config
	engine ndn.Engine
	svs    *ndn_sync.SvSync

	routers map[uint64]*PrefixTableRouter
	me      *PrefixTableRouter

	repo       map[uint64][]byte
	repoMutex  sync.RWMutex
	snapshotAt uint64
}

type PrefixTableRouter struct {
	Name     enc.Name
	Fetching bool
	Known    uint64
	Latest   uint64
	Prefixes map[uint64]*PrefixEntry
}

type PrefixEntry struct {
	Name enc.Name
}

func NewPrefixTable(
	config *config.Config,
	engine ndn.Engine,
	svs *ndn_sync.SvSync,
) *PrefixTable {
	pt := &PrefixTable{
		config: config,
		engine: engine,
		svs:    svs,

		routers: make(map[uint64]*PrefixTableRouter),
		me:      nil,

		repo:      make(map[uint64][]byte),
		repoMutex: sync.RWMutex{},
	}

	pt.me = pt.GetRouter(config.RouterPfxN)
	pt.me.Known = svs.GetSeqNo(config.RouterPfxN)
	pt.me.Latest = pt.me.Known
	pt.publishSnap()

	return pt
}

func (pt *PrefixTable) GetRouter(name enc.Name) *PrefixTableRouter {
	hash := name.Hash()
	router := pt.routers[hash]
	if router == nil {
		router = &PrefixTableRouter{
			Name:     name,
			Prefixes: make(map[uint64]*PrefixEntry),
		}
		pt.routers[hash] = router
	}
	return router
}

func (pt *PrefixTable) Announce(name enc.Name) {
	log.Infof("Announcing prefix %s", name)
	pt.me.Prefixes[name.Hash()] = &PrefixEntry{Name: name}

	op := tlv.PrefixOpList{
		ExitRouter: &tlv.Destination{Name: pt.config.RouterPfxN},
		PrefixOpAdds: []*tlv.PrefixOpAdd{{
			Name: name,
			Cost: 1,
		}},
	}
	pt.publishOp(op.Encode())
}

func (pt *PrefixTable) Withdraw(name enc.Name) {
	log.Infof("Withdrawing prefix %s", name)
	delete(pt.me.Prefixes, name.Hash())

	op := tlv.PrefixOpList{
		ExitRouter:      &tlv.Destination{Name: pt.config.RouterPfxN},
		PrefixOpRemoves: []*tlv.PrefixOpRemove{{Name: name}},
	}
	pt.publishOp(op.Encode())
}

func (pt *PrefixTable) Apply(ops *tlv.PrefixOpList) {
	if ops.ExitRouter == nil || len(ops.ExitRouter.Name) == 0 {
		log.Warn("PrefixOpList has no ExitRouter")
		return
	}

	router := pt.GetRouter(ops.ExitRouter.Name)

	if ops.PrefixOpReset {
		log.Infof("Reset prefix table for %s", ops.ExitRouter.Name)
		router.Prefixes = make(map[uint64]*PrefixEntry)
	}

	for _, add := range ops.PrefixOpAdds {
		log.Infof("Added prefix for %s: %s", ops.ExitRouter.Name, add.Name)
		router.Prefixes[add.Name.Hash()] = &PrefixEntry{Name: add.Name}
	}

	for _, remove := range ops.PrefixOpRemoves {
		log.Infof("Removed prefix for %s: %s", ops.ExitRouter.Name, remove.Name)
		delete(router.Prefixes, remove.Name.Hash())
	}
}

func (pt *PrefixTable) publishOp(content enc.Wire) {
	// Increment our sequence number
	seq := pt.svs.IncrSeqNo(pt.config.RouterPfxN)
	pt.me.Known = seq
	pt.me.Latest = seq

	// Create the new data
	name := append(pt.config.PfxDataPfxN, enc.NewSequenceNumComponent(seq))
	pt.publish(name, content)

	// Create snapshot if needed
	if pt.snapshotAt-seq >= 100 {
		pt.publishSnap()
	}
}

func (pt *PrefixTable) publishSnap() {
	snap := tlv.PrefixOpList{
		ExitRouter:    &tlv.Destination{Name: pt.config.RouterPfxN},
		PrefixOpReset: true,
		PrefixOpAdds:  make([]*tlv.PrefixOpAdd, 0, len(pt.me.Prefixes)),
	}

	for _, entry := range pt.me.Prefixes {
		snap.PrefixOpAdds = append(snap.PrefixOpAdds, &tlv.PrefixOpAdd{
			Name: entry.Name,
			Cost: 1,
		})
	}

	// Store snapshot in repo
	// TODO: this can be a segmented object
	pt.snapshotAt = pt.me.Latest
	snapPfx := append(pt.config.PfxDataPfxN,
		enc.NewStringComponent(enc.TypeKeywordNameComponent, "SNAP"))
	snapName := append(snapPfx, enc.NewSequenceNumComponent(pt.snapshotAt))
	pt.publish(snapName, snap.Encode())

	// Point prefix to the snapshot
	pt.repoMutex.Lock()
	defer pt.repoMutex.Unlock()
	pt.repo[snapPfx.Hash()] = pt.repo[snapName.Hash()]
}

func (pt *PrefixTable) publish(name enc.Name, content enc.Wire) {
	// TODO: sign the prefix table data
	signer := security.NewSha256Signer()

	wire, _, err := pt.engine.Spec().MakeData(
		name,
		&ndn.DataConfig{
			ContentType: utils.IdPtr(ndn.ContentTypeBlob),
			Freshness:   utils.IdPtr(1 * time.Second),
		},
		content,
		signer)
	if err != nil {
		log.Warnf("advertDataOnInterest: Failed to make Data: %+v", err)
		return
	}

	// Store the data packet in our mem repo
	pt.repoMutex.Lock()
	defer pt.repoMutex.Unlock()
	pt.repo[name.Hash()] = wire.Join()
}

func (pt *PrefixTable) OnDataInterestAsync(
	interest ndn.Interest,
	reply ndn.ReplyFunc,
	extra ndn.InterestHandlerExtra,
) {
	go pt.onDataInterest(interest, reply, extra)
}

// Received prefix data Interest
func (pt *PrefixTable) onDataInterest(
	interest ndn.Interest,
	reply ndn.ReplyFunc,
	extra ndn.InterestHandlerExtra,
) {
	// TODO: remove old entries from repo

	pt.repoMutex.RLock()
	defer pt.repoMutex.RUnlock()

	// Find exact match in repo
	name := interest.Name()
	if data := pt.repo[name.Hash()]; data != nil {
		err := reply(enc.Wire{data})
		if err != nil {
			log.Warnf("advertDataOnInterest: Failed to reply: %+v", err)
		}
		return
	}

	log.Warnf("Failed to find data for for %s", name)
}
