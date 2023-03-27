package service

// does not generate any code, but does validate our configuration
// // go:generate wire check
// cannot use this as wire stumbles on the generics

import (
	"github.com/google/wire"
	"github.com/jtyers/tmaas-model/validator"
	"github.com/jtyers/tmaas-service-util/id"
	"github.com/jtyers/tmaas-threat-service/dao"
)

var serviceDepsProviderSet = wire.NewSet(
	id.NewDefaultRandomIDProvider,
	wire.Bind(new(id.RandomIDProvider), new(*id.DefaultRandomIDProvider)),

	validator.StructValidatorProviderSet,
)

var ThreatServiceProviderSet = wire.NewSet(
	serviceDepsProviderSet,

	dao.ThreatDaoProviderSet,
	wire.Bind(new(ThreatService), new(*DefaultThreatService)),
	NewDefaultThreatService,
)

var ThreatModelServiceProviderSet = wire.NewSet(
	serviceDepsProviderSet,

	dao.ThreatModelDaoProviderSet,
	wire.Bind(new(ThreatModelService), new(*DefaultThreatModelService)),
	NewDefaultThreatModelService,
)
