package dao

import "github.com/jtyers/tmaas-service-dao/clover"

// A custom CollectionConfig that configures the collection for storing threats.
// We define as its own type so we can inject it as a dependency.
type ThreatControlCloverCollectionConfig struct {
	*clover.DefaultCollectionConfig
}

func NewThreatControlCloverCollectionConfig() ThreatControlCloverCollectionConfig {
	return ThreatControlCloverCollectionConfig{
		clover.NewDefaultCollectionConfig("threat_controls"),
	}
}
