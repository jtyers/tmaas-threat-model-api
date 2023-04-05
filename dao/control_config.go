package dao

import "github.com/jtyers/tmaas-service-dao/clover"

// A custom CollectionConfig that configures the collection for storing controls.
// We define as its own type so we can inject it as a dependency.
type ControlCloverCollectionConfig struct {
	*clover.DefaultCollectionConfig
}

func NewControlCloverCollectionConfig() ControlCloverCollectionConfig {
	return ControlCloverCollectionConfig{
		clover.NewDefaultCollectionConfig("controls"),
	}
}
