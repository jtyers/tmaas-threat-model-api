package dao

//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mocks.go -package $GOPACKAGE

import (
	gdatastore "cloud.google.com/go/datastore"
	m "github.com/jtyers/tmaas-model"
	servicedao "github.com/jtyers/tmaas-service-dao"
	"github.com/jtyers/tmaas-service-dao/datastore"
)

// ThreatModelDao is needed because wire does not directly support
// generics, but does appear to support injecting non-generic types
// that contain embedded types
type ThreatModelDao interface {
	// TODO You will hit a mockgen error here relating to *ast.IndexExpr
	// if you use vanilla mockgen. Sadly, gomock has not had much maintenance
	// and its generics support only exists in an un-merged (but tested) PR.
	// To get the working version of mockgen:
	//  1. clone https://github.com/bradleygore/gomock
	//  2. checkout task_HOSPENG-4373-gomock-generics
	//  3. run `go install ./...`
	servicedao.IDTypedDao[m.ThreatModelID, m.ThreatModel]
}

type ThreatModelIDCreator struct{}

func NewThreatModelIDCreator() ThreatModelIDCreator {
	return ThreatModelIDCreator{}
}
func (ThreatModelIDCreator) Create(id string) m.ThreatModelID {
	return m.ThreatModelID(id)
}

func NewThreatModelDao(client *gdatastore.Client, config datastore.DatastoreConfiguration, idCreator ThreatModelIDCreator) (ThreatModelDao, error) {
	var errorMappings map[error]error

	return datastore.NewDatastoreDao[m.ThreatModelID, m.ThreatModel, *m.ThreatModel](client, config, idCreator, errorMappings)
}
