package service

// does not generate any code, but does validate our configuration
// // go:generate wire check
// cannot use this as wire stumbles on the generics

import (
	"github.com/google/wire"
	dfdclient "github.com/jtyers/tmaas-dfd-api/client"
	"github.com/jtyers/tmaas-model/validator"
	"github.com/jtyers/tmaas-service-util/idchecker"
	thclient "github.com/jtyers/tmaas-threat-api/client"
)

var ServiceDepsProviderSet = wire.NewSet(
	validator.StructValidatorProviderSet,
)

func NewIDCheckerForTypes(dfd *dfdclient.ClientDataFlowDiagramIDChecker) idchecker.IDCheckerForTypes {
	return idchecker.IDCheckerForTypes([]idchecker.IDCheckerForType{
		dfd,
	})
}

var ThreatModelServiceProviderSet = wire.NewSet(
	ServiceDepsProviderSet,

	wire.Bind(new(ThreatModelService), new(*DefaultThreatModelService)),
	NewDefaultThreatModelService,

	NewServiceThreatModelIDChecker,

	wire.Bind(new(idchecker.IDChecker), new(*idchecker.DefaultIDChecker)),
	idchecker.NewDefaultIDChecker,

	dfdclient.DataFlowDiagramServiceClientProviderSet,
	thclient.ThreatServiceClientMinimalProviderSet,

	NewIDCheckerForTypes,
)
