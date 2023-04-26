//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/google/wire"
	"github.com/jtyers/tmaas-api/dao"
	"github.com/jtyers/tmaas-api/service"
	"github.com/jtyers/tmaas-api/web"
)

func InitialiseRouter() (http.Handler, error) {
	wire.Build(
		dao.ThreatModelDaoProviderSet,
		service.ThreatModelServiceProviderSet,
		web.ThreatModelWebProviderSet,
	)
	return nil, nil
}
