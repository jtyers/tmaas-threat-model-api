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
		BaseURL: serviceutil.EnsureSuffix(serviceutil.GetEnvWithDefault("THREATPLANE_THREAT_MODEL_API_URL", "https://threatmodel.api.threatplane.io/"), "/"),
	}
}

var ThreatModelServiceClientProviderSet = wire.NewSet(
	requestor.RequestorWithContextProviderSet,
	ThreatModelServiceClientMinimalProviderSet,
)

// The minimal provider set provides just the types for the
// client itself, not the types the client depends on. This pattern
// is useful where a downstream app needs multiple clients.
// Only one 'full' client provider set should be used (doesn't
// matter which generally) and other clients should use the
// minimal sets. This avoids duplicate provides for core dependencies
// that all clients use, like RequestorWithContext.
var ThreatModelServiceClientMinimalProviderSet = wire.NewSet(
	NewThreatModelServiceClientConfig,

	wire.Bind(new(service.ThreatModelService), new(*ThreatModelServiceClient)),
	NewThreatModelServiceClient,

	NewClientThreatModelIDChecker,
)
