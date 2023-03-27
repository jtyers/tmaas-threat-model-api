package dao

// does not generate any code, but does validate our configuration
// // go:generate wire check
// cannot use this as wire stumbles on the generics

import (
	"github.com/google/wire"
	"github.com/jtyers/tmaas-service-dao/clover"
	util "github.com/jtyers/tmaas-service-util"
)

func NewCloverDatabaseConfig() clover.DatabaseConfig {
	return clover.NewDefaultDatabaseConfig(
		util.GetEnvWithDefault("THREAT_DB_PATH", "threat-db"),
	)
}

var DaoDepsProviderSet = wire.NewSet(
	NewCloverDatabaseConfig,
	clover.DaoProviderSet,
)

var ThreatDaoProviderSet = wire.NewSet(
	NewThreatCloverCollectionConfig,
	NewThreatCloverDao,
	NewThreatDao,
)

var ThreatModelDaoProviderSet = wire.NewSet(
	NewThreatModelCloverCollectionConfig,
	NewThreatModelCloverDao,
	NewThreatModelDao,
)

//wire.Bind(new(dao.Dao), new(*CloverDao)),
//NewCloverDao,
