package dao

//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mocks.go -package $GOPACKAGE

import (
	m "github.com/jtyers/tmaas-model"
	servicedao "github.com/jtyers/tmaas-service-dao"
	"github.com/jtyers/tmaas-service-dao/clover"
)

// ThreatDao is needed because wire does not directly support
// generics, but does appear to support injecting non-generic types
// that contain embedded types
type ThreatDao interface {
	// TODO You will hit a mockgen error here relating to *ast.IndexExpr
	// if you use vanilla mockgen. Sadly, gomock has not had much maintenance
	// and its generics support only exists in an un-merged (but tested) PR.
	// To get the working version of mockgen:
	//  1. clone https://github.com/bradleygore/gomock
	//  2. checkout task_HOSPENG-4373-gomock-generics
	//  3. run `go install ./...`
	servicedao.TypedDao[m.Threat]
}

func NewThreatCloverDao(db clover.CloverDbWrapper, config ThreatCloverCollectionConfig) *clover.CloverDao[ThreatCloverCollectionConfig] {
	return clover.NewCloverDao(db, config)
}

func NewThreatDao(dao *clover.CloverDao[ThreatCloverCollectionConfig]) ThreatDao {
	return servicedao.NewDefaultTypedDao[m.Threat](dao, nil)
}
