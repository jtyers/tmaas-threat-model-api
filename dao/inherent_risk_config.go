package dao

import "github.com/jtyers/tmaas-service-dao/clover"

// A custom CollectionConfig that configures the collection for storing threats.
// We define as its own type so we can inject it as a dependency.
type InherentRiskCloverCollectionConfig struct {
	*clover.DefaultCollectionConfig
}

func NewInherentRiskCloverCollectionConfig() InherentRiskCloverCollectionConfig {
	return InherentRiskCloverCollectionConfig{
		clover.NewDefaultCollectionConfig("inherent_risks"),
	}
}
