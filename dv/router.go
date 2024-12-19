package dv

import (
	"sync"
	"time"

	"github.com/pulsejet/go-ndn-dv/config"
	"github.com/pulsejet/go-ndn-dv/nfdc"
	"github.com/pulsejet/go-ndn-dv/table"
	enc "github.com/zjkmxy/go-ndn/pkg/encoding"
	"github.com/zjkmxy/go-ndn/pkg/ndn"
	mgmt "github.com/zjkmxy/go-ndn/pkg/ndn/mgmt_2022"
	ndn_sync "github.com/zjkmxy/go-ndn/pkg/sync"
	"github.com/zjkmxy/go-ndn/pkg/utils"
)

type Router struct {
	// go-ndn app that this router is attached to
	engine ndn.Engine
	// config for this router
	config *config.Config
	// nfd management thread
	nfdc *nfdc.NfdMgmtThread
	// single mutex for all operations
	mutex sync.Mutex

	// channel to stop the DV
	stop chan bool
	// heartbeat for outgoing Advertisements
	heartbeat *time.Ticker
	// deadcheck for neighbors
	deadcheck *time.Ticker

	// neighbor table
	neighbors *table.NeighborTable
	// routing information base
	rib *table.Rib
	// prefix table
	pfx *table.PrefixTable
	// forwarding table
	fib *table.Fib

	// advertisement sequence number for self
	advertSyncSeq uint64
	// prefix table svs instance
	pfxSvs *ndn_sync.SvSync
}

// Create a new DV router.
func NewRouter(config *config.Config, engine ndn.Engine) (*Router, error) {
	// Validate configuration
	err := config.Parse()
	if err != nil {
		return nil, err
	}

	// Create the DV router
	dv := &Router{
		engine: engine,
		config: config,
		nfdc:   nfdc.NewNfdMgmtThread(engine),
		mutex:  sync.Mutex{},
	}

	// Create sync groups
	dv.pfxSvs = ndn_sync.NewSvSync(engine, config.PrefixTableSyncPrefix(), dv.onPfxSyncUpdate)

	// Set initial sequence numbers
	now := uint64(time.Now().UnixMilli())
	dv.advertSyncSeq = now
	dv.pfxSvs.SetSeqNo(dv.config.RouterName(), now)

	// Create tables
	dv.neighbors = table.NewNeighborTable(config, dv.nfdc)
	dv.rib = table.NewRib(config)
	dv.pfx = table.NewPrefixTable(config, engine, dv.pfxSvs)
	dv.fib = table.NewFib(config, dv.nfdc)

	return dv, nil
}

// Start the DV router. Blocks until Stop() is called.
func (dv *Router) Start() (err error) {
	dv.stop = make(chan bool, 1)

	// Start timers
	dv.heartbeat = time.NewTicker(dv.config.AdvertisementSyncInterval())
	dv.deadcheck = time.NewTicker(dv.config.RouterDeadInterval())
	defer dv.heartbeat.Stop()
	defer dv.deadcheck.Stop()

	// Start management thread
	go dv.nfdc.Start()
	defer dv.nfdc.Stop()

	// Configure face
	err = dv.configureFace()
	if err != nil {
		return err
	}

	// Register interest handlers
	err = dv.register()
	if err != nil {
		return err
	}

	// Start sync groups
	dv.pfxSvs.Start()
	defer dv.pfxSvs.Stop()

	// Add self to the RIB
	dv.rib.Set(dv.config.RouterName(), dv.config.RouterName(), 0)

	for {
		select {
		case <-dv.heartbeat.C:
			dv.advertSyncSendInterest()
		case <-dv.deadcheck.C:
			dv.checkDeadNeighbors()
		case <-dv.stop:
			return nil
		}
	}
}

// Stop the DV router.
func (dv *Router) Stop() {
	dv.stop <- true
}

// Configure the face to forwarder.
func (dv *Router) configureFace() (err error) {
	// Enable local fields on face. This includes incoming face indication.
	dv.nfdc.Exec(nfdc.NfdMgmtCmd{
		Module: "faces",
		Cmd:    "update",
		Args: &mgmt.ControlArgs{
			Mask:  utils.IdPtr(uint64(0x01)),
			Flags: utils.IdPtr(uint64(0x01)),
		},
		Retries: -1,
	})

	return nil
}

// Register interest handlers for DV prefixes.
func (dv *Router) register() (err error) {
	// Advertisement Sync
	err = dv.engine.AttachHandler(dv.config.AdvertisementSyncPrefix(), func(args ndn.InterestHandlerArgs) {
		go dv.advertSyncOnInterest(args)
	})
	if err != nil {
		return err
	}

	// Advertisement Data
	err = dv.engine.AttachHandler(dv.config.AdvertisementDataPrefix(), func(args ndn.InterestHandlerArgs) {
		go dv.advertDataOnInterest(args)
	})
	if err != nil {
		return err
	}

	// Prefix Data
	err = dv.engine.AttachHandler(dv.config.PrefixTableDataPrefix(), func(args ndn.InterestHandlerArgs) {
		go dv.pfx.OnDataInterest(args)
	})
	if err != nil {
		return err
	}

	// Readvertise Data
	err = dv.engine.AttachHandler(dv.config.ReadvertisePrefix(), func(args ndn.InterestHandlerArgs) {
		go dv.readvertiseOnInterest(args)
	})
	if err != nil {
		return err
	}

	// Register routes to forwarder
	pfxs := []enc.Name{
		dv.config.AdvertisementSyncPrefix(),
		dv.config.AdvertisementDataPrefix(),
		dv.config.PrefixTableSyncPrefix(),
		dv.config.PrefixTableDataPrefix(),
		dv.config.ReadvertisePrefix(),
	}
	for _, prefix := range pfxs {
		dv.nfdc.Exec(nfdc.NfdMgmtCmd{
			Module: "rib",
			Cmd:    "register",
			Args: &mgmt.ControlArgs{
				Name:   prefix,
				Cost:   utils.IdPtr(uint64(0)),
				Origin: utils.IdPtr(config.NlsrOrigin),
			},
			Retries: -1,
		})
	}

	// Set strategy to multicast for sync prefixes
	mcast, _ := enc.NameFromStr(config.MulticastStrategy)
	pfxs = []enc.Name{
		dv.config.AdvertisementSyncPrefix(),
		dv.config.PrefixTableSyncPrefix(),
	}
	for _, prefix := range pfxs {
		dv.nfdc.Exec(nfdc.NfdMgmtCmd{
			Module: "strategy-choice",
			Cmd:    "set",
			Args: &mgmt.ControlArgs{
				Name: prefix,
				Strategy: &mgmt.Strategy{
					Name: mcast,
				},
			},
			Retries: -1,
		})
	}

	return nil
}
