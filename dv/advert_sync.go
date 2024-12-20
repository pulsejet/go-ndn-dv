package dv

import (
	"time"

	enc "github.com/zjkmxy/go-ndn/pkg/encoding"
	"github.com/zjkmxy/go-ndn/pkg/log"
	"github.com/zjkmxy/go-ndn/pkg/ndn"
	ndn_sync "github.com/zjkmxy/go-ndn/pkg/sync"
	"github.com/zjkmxy/go-ndn/pkg/utils"
)

func (dv *Router) advertSyncSendInterest() (err error) {
	// SVS v2 Sync Interest
	syncName := append(dv.config.AdvertisementSyncPrefix(), enc.NewVersionComponent(2))

	// Sync Interest parameters for SVS
	cfg := &ndn.InterestConfig{
		MustBeFresh: true,
		Lifetime:    utils.IdPtr(1 * time.Millisecond),
		Nonce:       utils.ConvertNonce(dv.engine.Timer().Nonce()),
		HopLimit:    utils.IdPtr(uint(2)), // use localhop w/ this
	}

	// State Vector for our group
	// TODO: switch to new TLV types
	sv := &ndn_sync.StateVectorAppParam{
		StateVector: &ndn_sync.StateVector{
			Entries: []*ndn_sync.StateVectorEntry{{
				NodeId: dv.config.RouterName(),
				SeqNo:  dv.advertSyncSeq,
			}},
		},
	}

	// TODO: sign the sync interest

	interest, err := dv.engine.Spec().MakeInterest(syncName, cfg, sv.Encode(), nil)
	if err != nil {
		return err
	}

	// Sync Interest has no reply
	err = dv.engine.Express(interest, nil)
	if err != nil {
		return err
	}

	return nil
}

func (dv *Router) advertSyncOnInterest(args ndn.InterestHandlerArgs) {
	// Check if app param is present
	if args.Interest.AppParam() == nil {
		log.Warn("advertSyncOnInterest: received Sync Interest with no AppParam, ignoring")
		return
	}

	// If there is no incoming face ID, we can't use this
	if args.IncomingFaceId == nil {
		log.Warn("advertSyncOnInterest: received Sync Interest with no incoming face ID, ignoring")
		return
	}

	// TODO: verify signature on Sync Interest

	params, err := ndn_sync.ParseStateVectorAppParam(enc.NewWireReader(args.Interest.AppParam()), true)
	if err != nil || params.StateVector == nil {
		log.Warnf("advertSyncOnInterest: failed to parse StateVec: %+v", err)
		return
	}

	// Process each entry in the state vector
	dv.mutex.Lock()
	defer dv.mutex.Unlock()

	for _, entry := range params.StateVector.Entries {
		// Parse name from NodeId
		nodeId := entry.NodeId
		if nodeId == nil {
			log.Warnf("advertSyncOnInterest: failed to parse NodeId: %+v", err)
			continue
		}

		// Check if the entry is newer than what we know
		ns := dv.neighbors.Get(nodeId)
		if ns != nil {
			if ns.AdvertSeq >= entry.SeqNo {
				// Nothing has changed, skip
				ns.RecvPing(*args.IncomingFaceId)
				continue
			}
		} else {
			// Create new neighbor entry cause none found
			// This is the ONLY place where neighbors are created
			// In all other places, quit if not found
			ns = dv.neighbors.Add(nodeId)
		}

		ns.RecvPing(*args.IncomingFaceId)
		ns.AdvertSeq = entry.SeqNo

		go dv.advertDataFetch(nodeId, entry.SeqNo)
	}
}

func (dv *Router) advertSyncNotifyNew() {
	dv.mutex.Lock()
	defer dv.mutex.Unlock()

	dv.advertSyncSeq++
	go dv.advertSyncSendInterest()
}
