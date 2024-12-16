package dv

import (
	"github.com/pulsejet/go-ndn-dv/config"
	"github.com/pulsejet/go-ndn-dv/table"
	enc "github.com/zjkmxy/go-ndn/pkg/encoding"
	"github.com/zjkmxy/go-ndn/pkg/log"
)

// Compute the RIB chnages for this neighbor
func (dv *Router) ribUpdate(ns *table.NeighborState) {
	dv.mutex.Lock()
	defer dv.mutex.Unlock()

	if ns.Advert == nil {
		return
	}

	// TODO: use cost to neighbor
	localCost := uint64(1)

	// Trigger our own advertisement if needed
	var dirty bool = false

	// Reset destinations for this neighbor
	dv.rib.DirtyResetNextHop(ns.Name)

	for _, entry := range ns.Advert.Entries {
		// Use the advertised cost by default
		cost := entry.Cost + localCost

		// Poison reverse - try other cost if next hop is us
		if entry.NextHop.Name.Equal(dv.config.RouterPfxN) {
			if entry.OtherCost < config.CostInfinity {
				cost = entry.OtherCost + localCost
			} else {
				cost = config.CostInfinity
			}
		}

		// Skip unreachable destinations
		if cost >= config.CostInfinity {
			continue
		}

		// Check advertisement changes
		dirty = dv.rib.Set(entry.Destination.Name, ns.Name, cost) || dirty
	}

	// Drop dead entries
	dirty = dv.rib.Prune() || dirty

	// If advert changed, increment sequence number
	if dirty {
		go dv.fibUpdate()
		go dv.advertSyncNotifyNew()
		go dv.prefixDataFetchAll()
	}
}

// Check for dead neighbors
func (dv *Router) checkDeadNeighbors() {
	dv.mutex.Lock()
	defer dv.mutex.Unlock()

	dirty := false
	for _, ns := range dv.neighbors.GetAll() {
		// Check if the neighbor is entirely dead
		if ns.IsDead() {
			log.Infof("checkDeadNeighbors: Neighbor %s is dead", ns.Name.String())

			// This is the ONLY place that can remove neighbors
			dv.neighbors.Remove(ns.Name)

			// Remove neighbor from RIB and prune
			dirty = dv.rib.RemoveNextHop(ns.Name) || dirty
			dirty = dv.rib.Prune() || dirty
		}
	}

	if dirty {
		go dv.fibUpdate()
		go dv.advertSyncNotifyNew()
	}
}

// Update the FIB
func (dv *Router) fibUpdate() {
	dv.mutex.Lock()
	defer dv.mutex.Unlock()

	// Name prefixes from global prefix table as well as RIB
	names := make(map[uint64]enc.Name)
	fibEntries := make(map[uint64][]table.FibEntry)

	// Update paths to all routers from RIB
	for _, router := range dv.rib.Entries() {
		rNameH := router.Name().Hash()
		names[rNameH] = router.Name()

		// Skip if this is us
		if router.Name().Equal(dv.config.RouterPfxN) {
			continue
		}

		// Get FIB entry to reach this router
		fes := dv.rib.GetFibEntries(dv.neighbors, rNameH)

		// Add entry to the router itself
		fibEntries[rNameH] = append(fibEntries[rNameH], fes...)

		// Add entries to all prefixes announced by this router
		for _, pfx := range dv.pfx.GetRouter(router.Name()).Prefixes {
			pNameH := pfx.Name.Hash()
			names[pNameH] = pfx.Name

			// Use the same nexthop entries as the exit router itself
			// De-duplication is done by the fib table update function
			fibEntries[pNameH] = append(fibEntries[pNameH], fes...)
		}
	}

	// Update all FIB entries to NFD
	dv.fib.UnmarkAll()
	for nameH, fes := range fibEntries {
		if dv.fib.UpdateH(nameH, names[nameH], fes) {
			dv.fib.MarkH(nameH)
		}
	}
	dv.fib.RemoveUnmarked()
}
