package dao

import "github.com/jtyers/tmaas-service-dao/clover"

// A custom CollectionConfig that configures the collection for storing threats.
// We define as its own type so we can inject it as a dependency.
type ComponentCloverCollectionConfig struct {
	*clover.DefaultCollectionConfig
}

func NewComponentCloverCollectionConfig() ComponentCloverCollectionConfig {
	return ComponentCloverCollectionConfig{
		clover.NewDefaultCollectionConfig("components"),
	}
}
