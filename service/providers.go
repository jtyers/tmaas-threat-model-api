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

var ThreatServiceProviderSet = wire.NewSet(
	dao.ThreatDaoProviderSet,

	id.NewDefaultRandomIdProvider,
	wire.Bind(new(id.RandomIdProvider), new(*id.DefaultRandomIdProvider)),

	validator.StructValidatorProviderSet,

	wire.Bind(new(ThreatService), new(*DefaultThreatService)),
	NewDefaultThreatService,
)
