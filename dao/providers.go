package dao

// does not generate any code, but does validate our configuration
// // go:generate wire check
// cannot use this as wire stumbles on the generics

import (
	"github.com/google/wire"
	"github.com/jtyers/tmaas-service-dao/clover"
)

func NewCloverConfiguration() clover.CloverConfiguration {
	return clover.CloverConfiguration{
		DatabaseName:   "db",
		CollectionName: "threats",
	}
}

var ThreatDaoProviderSet = wire.NewSet(
	NewCloverConfiguration,
	clover.DaoProviderSet,
	NewThreatDao,
)
