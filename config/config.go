package config

import (
	"errors"
	"time"

	enc "github.com/zjkmxy/go-ndn/pkg/encoding"
)

const CostInfinity = uint64(16)
const MulticastStrategy = "/localhost/nfd/strategy/multicast"
const NlsrOrigin = uint64(128)

var Localhop = enc.Name{enc.NewStringComponent(enc.TypeGenericNameComponent, "localhop")}
var Localhost = enc.Name{enc.NewStringComponent(enc.TypeGenericNameComponent, "localhost")}

type Config struct {
	// NetworkName should be the same for all routers in the network.
	NetworkName string
	// RouterName should be unique for each router in the network.
	RouterName string
	// Period of sending Advertisement Sync Interests.
	AdvertisementSyncInterval time.Duration
	// Time after which a neighbor is considered dead.
	RouterDeadInterval time.Duration

	// Parsed Global Prefix
	NetworkNameN enc.Name
	// Parsed Router Prefix
	RouterNameN enc.Name
	// Advertisement Sync Prefix
	AdvSyncPfxN enc.Name
	// Advertisement Data Prefix
	AdvDataPfxN enc.Name
	// Prefix Table Sync Prefix
	PfxSyncPfxN enc.Name
	// Prefix Table Data Prefix
	PfxDataPfxN enc.Name
	// NLSR readvertise prefix
	ReadvertisePfxN enc.Name
}

func (c *Config) Parse() (err error) {
	// Validate prefixes not empty
	if c.NetworkName == "" || c.RouterName == "" {
		return errors.New("NetworkName and RouterName must be set")
	}

	// Parse prefixes
	c.NetworkNameN, err = enc.NameFromStr(c.NetworkName)
	if err != nil {
		return err
	}

	c.RouterNameN, err = enc.NameFromStr(c.RouterName)
	if err != nil {
		return err
	}

	// Validate intervals are not too short
	if c.AdvertisementSyncInterval < 1*time.Second {
		return errors.New("AdvertisementSyncInterval must be at least 1 second")
	}

	// Dead interval at least 2 sync intervals
	if c.RouterDeadInterval < 2*c.AdvertisementSyncInterval {
		return errors.New("RouterDeadInterval must be at least 2*AdvertisementSyncInterval")
	}

	// Create name table
	c.AdvSyncPfxN = append(Localhop, append(c.NetworkNameN,
		enc.NewStringComponent(enc.TypeKeywordNameComponent, "DV"),
		enc.NewStringComponent(enc.TypeKeywordNameComponent, "ADS"),
	)...)
	c.AdvDataPfxN = append(Localhop, append(c.RouterNameN,
		enc.NewStringComponent(enc.TypeKeywordNameComponent, "DV"),
		enc.NewStringComponent(enc.TypeKeywordNameComponent, "ADV"),
	)...)
	c.PfxSyncPfxN = append(c.NetworkNameN,
		enc.NewStringComponent(enc.TypeKeywordNameComponent, "DV"),
		enc.NewStringComponent(enc.TypeKeywordNameComponent, "PFS"),
	)
	c.PfxDataPfxN = append(c.RouterNameN,
		enc.NewStringComponent(enc.TypeKeywordNameComponent, "DV"),
		enc.NewStringComponent(enc.TypeKeywordNameComponent, "PFX"),
	)
	c.ReadvertisePfxN = append(Localhost,
		enc.NewStringComponent(enc.TypeGenericNameComponent, "nlsr"),
	)

	return nil
}
