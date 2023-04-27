package dao

// does not generate any code, but does validate our configuration
// // go:generate wire check
// cannot use this as wire stumbles on the generics

import (
	"github.com/google/wire"
	"github.com/jtyers/tmaas-service-dao/datastore"
	util "github.com/jtyers/tmaas-service-util"
	"github.com/jtyers/tmaas-service-util/id"
)

func NewDatastoreConfig() datastore.DatastoreConfiguration {
	return datastore.DatastoreConfiguration{
		ProjectID:        util.GetEnv("PROJECT_ID"),
		DatastoreKeyKind: DatastoreKeyKind,
	}
}

func NewThreatModelRandomIDProviderPrefix() id.RandomIDProviderPrefix {
	return id.RandomIDProviderPrefix(ThreatModelIDPrefix)
}

var ThreatModelDaoProviderSet = wire.NewSet(
	datastore.DaoProviderSet,

	NewThreatModelRandomIDProviderPrefix,
	NewThreatModelDao,
	NewThreatModelIDCreator,
	NewDatastoreConfig,
)
