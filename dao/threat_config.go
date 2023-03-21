package dao

import "github.com/jtyers/tmaas-service-dao/clover"

// A custom CollectionConfig that configures the collection for storing threats.
// We define as its own type so we can inject it as a dependency.
type ThreatCloverCollectionConfig struct {
	*clover.DefaultCollectionConfig
}

func NewThreatCloverCollectionConfig() ThreatCloverCollectionConfig {
	return ThreatCloverCollectionConfig{
		clover.NewDefaultCollectionConfig("threats"),
	}
}
