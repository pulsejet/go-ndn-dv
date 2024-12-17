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
	// Network should be the same for all routers in the network.
	Network string `json:"network"`
	// Router should be unique for each router in the network.
	Router string `json:"router"`
	// Period of sending Advertisement Sync Interests.
	AdvertisementSyncInterval_ms uint64 `json:"advertise_interval"`
	// Time after which a neighbor is considered dead.
	RouterDeadInterval_ms uint64 `json:"router_dead_interval"`

	// Parsed Global Prefix
	networkNameN enc.Name
	// Parsed Router Prefix
	routerNameN enc.Name
	// Advertisement Sync Prefix
	advSyncPfxN enc.Name
	// Advertisement Data Prefix
	advDataPfxN enc.Name
	// Prefix Table Sync Prefix
	pfxSyncPfxN enc.Name
	// Prefix Table Data Prefix
	pfxDataPfxN enc.Name
	// NLSR readvertise prefix
	readvertisePfxN enc.Name
}

func DefaultConfig() Config {
	return Config{
		Network:                      "", // invalid
		Router:                       "", // invalid
		AdvertisementSyncInterval_ms: 5000,
		RouterDeadInterval_ms:        30000,
	}
}

func (c *Config) Parse() (err error) {
	// Validate prefixes not empty
	if c.Network == "" || c.Router == "" {
		return errors.New("network and router must be set")
	}

	// Parse prefixes
	c.networkNameN, err = enc.NameFromStr(c.Network)
	if err != nil {
		return err
	}

	c.routerNameN, err = enc.NameFromStr(c.Router)
	if err != nil {
		return err
	}

	// Validate intervals are not too short
	if c.AdvertisementSyncInterval() < 1*time.Second {
		return errors.New("AdvertisementSyncInterval must be at least 1 second")
	}

	// Dead interval at least 2 sync intervals
	if c.RouterDeadInterval() < 2*c.AdvertisementSyncInterval() {
		return errors.New("RouterDeadInterval must be at least 2*AdvertisementSyncInterval")
	}

	// Create name table
	c.advSyncPfxN = append(Localhop, append(c.networkNameN,
		enc.NewStringComponent(enc.TypeKeywordNameComponent, "DV"),
		enc.NewStringComponent(enc.TypeKeywordNameComponent, "ADS"),
	)...)
	c.advDataPfxN = append(Localhop, append(c.routerNameN,
		enc.NewStringComponent(enc.TypeKeywordNameComponent, "DV"),
		enc.NewStringComponent(enc.TypeKeywordNameComponent, "ADV"),
	)...)
	c.pfxSyncPfxN = append(c.networkNameN,
		enc.NewStringComponent(enc.TypeKeywordNameComponent, "DV"),
		enc.NewStringComponent(enc.TypeKeywordNameComponent, "PFS"),
	)
	c.pfxDataPfxN = append(c.routerNameN,
		enc.NewStringComponent(enc.TypeKeywordNameComponent, "DV"),
		enc.NewStringComponent(enc.TypeKeywordNameComponent, "PFX"),
	)
	c.readvertisePfxN = append(Localhost,
		enc.NewStringComponent(enc.TypeGenericNameComponent, "nlsr"),
	)

	return nil
}

func (c *Config) NetworkName() enc.Name {
	return c.networkNameN
}

func (c *Config) RouterName() enc.Name {
	return c.routerNameN
}

func (c *Config) AdvertisementSyncPrefix() enc.Name {
	return c.advSyncPfxN
}

func (c *Config) AdvertisementDataPrefix() enc.Name {
	return c.advDataPfxN
}

func (c *Config) PrefixTableSyncPrefix() enc.Name {
	return c.pfxSyncPfxN
}

func (c *Config) PrefixTableDataPrefix() enc.Name {
	return c.pfxDataPfxN
}

func (c *Config) ReadvertisePrefix() enc.Name {
	return c.readvertisePfxN
}

func (c *Config) AdvertisementSyncInterval() time.Duration {
	return time.Duration(c.AdvertisementSyncInterval_ms) * time.Millisecond
}

func (c *Config) RouterDeadInterval() time.Duration {
	return time.Duration(c.RouterDeadInterval_ms) * time.Millisecond
}
