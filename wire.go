//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/google/wire"
	"github.com/jtyers/tmaas-api/dao"
	"github.com/jtyers/tmaas-api/service"
	"github.com/jtyers/tmaas-api/web"

	"github.com/jtyers/tmaas-api-util/combo"
	"github.com/jtyers/tmaas-api-util/errors"
	"github.com/jtyers/tmaas-cors-config"
)

func InitialiseRouter() (http.Handler, error) {
	wire.Build(
		dao.DaoDepsProviderSet,
		service.ServiceDepsProviderSet,
		//service.ThreatServiceProviderSet,
		service.ThreatModelServiceProviderSet,

		corsconfig.CorsConfigProviderSet,

		combo.ComboMiddlewareFactoryProviderSet,
		errors.ErrorsMiddlewareFactoryProviderSet,

		web.NewRouter,
		web.NewThreatModelHandlers,
	)
	return nil, nil
}
