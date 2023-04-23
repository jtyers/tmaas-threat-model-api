package web

import (
	"github.com/google/wire"

	"github.com/jtyers/tmaas-api-util/combo"
	"github.com/jtyers/tmaas-api-util/errors"
	corsconfig "github.com/jtyers/tmaas-cors-config"
)

var ThreatModelWebProviderSet = wire.NewSet(
	corsconfig.CorsConfigProviderSet,

	combo.ComboMiddlewareFactoryProviderSet,
	errors.ErrorsMiddlewareFactoryProviderSet,

	NewRouter,
	NewThreatModelHandlers,
)
