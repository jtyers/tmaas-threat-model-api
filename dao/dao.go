package dao

//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mocks.go -package $GOPACKAGE

import (
	m "github.com/jtyers/tmaas-model"
	servicedao "github.com/jtyers/tmaas-service-dao"
)

// ThreatDao is needed because wire does not directly support
// generics, but does appear to support injecting non-generic types
// that contain embedded types
type ThreatDao interface {
	servicedao.TypedDao[m.Threat]
}

func NewThreatDao(dao servicedao.Dao) ThreatDao {
	return servicedao.NewDefaultTypedDao[m.Threat](dao)
}
