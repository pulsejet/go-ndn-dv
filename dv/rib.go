package dv

import (
	"github.com/pulsejet/go-ndn-dv/tlv"
	enc "github.com/zjkmxy/go-ndn/pkg/encoding"
)

// Routing Information Base (RIB)
type rib struct {
	// destination hash -> entry
	entries map[uint64]*rib_entry
	// neighbor hash -> neighbor name
	neighbors map[uint64]enc.Name
}

type rib_entry struct {
	// full name of destination router
	name enc.Name
	// neighbor hash -> cost
	costs map[uint64]uint64
	// next hop for lowest cost
	nextHop1 uint64
	// lowest cost in this entry
	lowest1 uint64
	// second lowest cost in this entry
	lowest2 uint64
}

func newRib() *rib {
	return &rib{
		entries:   make(map[uint64]*rib_entry),
		neighbors: make(map[uint64]enc.Name),
	}
}

func (r *rib) clear() {
	r.entries = make(map[uint64]*rib_entry)
	// keep neighbors; it's just a cache
}

// Set a destination in the RIB. Returns true if the Advertisement might change.
func (r *rib) set(destName enc.Name, nextHop enc.Name, cost uint64) bool {
	destHash := destName.Hash()
	nextHopHash := nextHop.Hash()

	// Create RIB entry if it doesn't exist
	entry, ok := r.entries[destHash]
	if !ok {
		entry = &rib_entry{
			name:  destName.Clone(),
			costs: make(map[uint64]uint64),
		}
		r.entries[destHash] = entry
	}

	// Create neighbor link if it doesn't exist
	if _, ok := r.neighbors[nextHopHash]; !ok {
		r.neighbors[nextHopHash] = nextHop.Clone()
	}

	return entry.set(nextHopHash, cost)
}

// Get all advertisement entries in the RIB.
func (r *rib) advert() *tlv.Advertisement {
	advert := &tlv.Advertisement{
		Entries: make([]*tlv.AdvEntry, 0, len(r.entries)),
	}

	for _, entry := range r.entries {
		advert.Entries = append(advert.Entries, &tlv.AdvEntry{
			Destination: &tlv.Destination{Name: entry.name},
			NextHop: &tlv.Destination{
				Name: r.neighbors[entry.nextHop1],
			},
			Cost:      entry.lowest1,
			OtherCost: entry.lowest2,
		})
	}

	return advert
}

func (e *rib_entry) set(nextHop uint64, cost uint64) bool {
	if known, ok := e.costs[nextHop]; ok && known == cost {
		// Prevent triggering an unnecessary refresh
		return false
	}

	e.costs[nextHop] = cost
	return e.refresh()
}

// Update lowest and second lowest costs for the entry.
func (e *rib_entry) refresh() bool {
	lowest1 := CostInfinity
	lowest2 := CostInfinity
	nextHop1 := uint64(0)

	for hop, cost := range e.costs {
		if cost < lowest1 {
			lowest2 = lowest1
			lowest1 = cost
			nextHop1 = hop
		} else if cost < lowest2 {
			lowest2 = cost
		}
	}

	if e.lowest1 != lowest1 || e.lowest2 != lowest2 || e.nextHop1 != nextHop1 {
		e.lowest1 = lowest1
		e.lowest2 = lowest2
		e.nextHop1 = nextHop1
		return true
	}

	return false
}
