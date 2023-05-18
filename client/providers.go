package client

//go:generate wire check

import (
	"github.com/google/wire"
	serviceutil "github.com/jtyers/tmaas-service-util"
	"github.com/jtyers/tmaas-service-util/requestor"
	"github.com/jtyers/tmaas-threat-model-api/service"
)

func NewThreatModelServiceClientConfig() ThreatModelServiceClientConfig {
	return ThreatModelServiceClientConfig{
		BaseURL: serviceutil.GetEnvWithDefault("TMAAS_DEVICE_API_URL", "https://device.api.tmas.io/"),
	}
}

var ThreatModelServiceClientProviderSet = wire.NewSet(
	requestor.RequestorWithContextProviderSet,
	NewThreatModelServiceClientConfig,

	wire.Bind(new(service.ThreatModelService), new(*ThreatModelServiceClient)),
	NewThreatModelServiceClient,

	NewClientThreatModelIDChecker,
)
