package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	jwt "github.com/jtyers/gin-jwt/v2"

	"github.com/jtyers/tmaas-api-util/combo"
	"github.com/jtyers/tmaas-api-util/errors"
	corsconfig "github.com/jtyers/tmaas-cors-config"
	m "github.com/jtyers/tmaas-model"
	"github.com/jtyers/tmaas-service-util/log"
	"github.com/jtyers/tmaas-threat-model-api/service"
)

func NewRouter(handlers *ThreatModelHandlers, comboFactory combo.ComboMiddlewareFactory, errorsMiddlewareFactory errors.ErrorsMiddlewareFactory, corsMiddleware corsconfig.CorsMiddleware) http.Handler {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware.Handler())

	r.Use(comboFactory.ExtractTokensToContext())

	r.Use(errorsMiddlewareFactory.NewErrorMiddleware([]errors.ErrorConfig{
		errors.NewErrorConfig(errors.ForExact(service.ErrNoSuchThreatModel), errors.StatusCode(http.StatusNotFound)),
		errors.NewErrorConfig(errors.ForExact(errors.ErrUnauthorized), errors.StatusCode(http.StatusUnauthorized)),
		errors.NewErrorConfig(errors.ForValidationErrors(), errors.ConvertValidationErrors()),
	}))

	r.NoRoute(comboFactory.StrictPermission(m.PermissionReadOwnThreatModels), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Infof("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r.PUT("/api/v1/threatModel",
		comboFactory.StrictUserPermission(m.PermissionReadOwnThreatModels),
		handlers.PutThreatModelHandler,
	)

	r.GET("/api/v1/threatModel",
		comboFactory.StrictUserPermission(m.PermissionReadOwnThreatModels),
		handlers.GetThreatModelsHandler,
	)
	r.GET("/api/v1/threatModel/:threatModelID",
		comboFactory.StrictPermission(m.PermissionReadOwnThreatModels), // Permit service accounts to access this
		handlers.GetThreatModelHandler,
	)

	r.DELETE("/api/v1/threatModel/:threatModelID",
		comboFactory.StrictUserPermission(m.PermissionEditOwnThreatModels),
		handlers.DeleteThreatModelHandler,
	)
	r.PATCH("/api/v1/threatModel/:threatModelID",
		comboFactory.StrictPermission(m.PermissionReadOwnThreatModels),
		handlers.PatchThreatModelHandler,
	)

	return r
}
