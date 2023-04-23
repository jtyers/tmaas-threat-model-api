package dao

// does not generate any code, but does validate our configuration
// // go:generate wire check
// cannot use this as wire stumbles on the generics

import (
	"github.com/google/wire"
	"github.com/jtyers/tmaas-service-dao/firestore"
	serviceutil "github.com/jtyers/tmaas-service-util"
)

var ThreatModelDaoProviderSet = wire.NewSet(
	firestore.DaoProviderSet,

	NewConfiguration,
	NewFirestoreConfiguration,
	NewThreatModelDao,
)

type Configuration struct {
	// Google Project ID, required by many Google Cloud APIs
	ProjectID string

	// Name of topic to publish User Changes to
	UserChangesTopicName string
}

func NewConfiguration() Configuration {
	return Configuration{
		ProjectID:            serviceutil.GetEnv("PROJECT_ID"),
		UserChangesTopicName: serviceutil.GetEnv("PUBSUB_USER_CHANGES_TOPIC"),
	}
}

func NewFirestoreConfiguration(c Configuration) firestore.FirestoreConfiguration {
	return firestore.FirestoreConfiguration{
		ProjectID:      c.ProjectID,
		CollectionName: "threat-models", // hard-coded :-)
	}
}
