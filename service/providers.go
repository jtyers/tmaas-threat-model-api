package service

// does not generate any code, but does validate our configuration
// // go:generate wire check
// cannot use this as wire stumbles on the generics

import (
	"github.com/google/wire"
	"github.com/jtyers/tmaas-model/validator"
	"github.com/jtyers/tmaas-service-util/idchecker"
)

var ServiceDepsProviderSet = wire.NewSet(
	validator.StructValidatorProviderSet,
)

var ThreatModelServiceProviderSet = wire.NewSet(
	ServiceDepsProviderSet,

	wire.Bind(new(ThreatModelService), new(*DefaultThreatModelService)),
	NewDefaultThreatModelService,

	NewServiceThreatModelIDChecker,

	wire.Bind(new(idchecker.IDChecker), new(*idchecker.DefaultIDChecker)),
	idchecker.NewDefaultIDChecker,
)
