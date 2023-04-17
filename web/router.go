package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	jwt "github.com/jtyers/gin-jwt/v2"
	"github.com/jtyers/tmaas-api-util/errors"
	"github.com/jtyers/tmaas-api/service"
	corsconfig "github.com/jtyers/tmaas-cors-config"
	m "github.com/jtyers/tmaas-model"

	"github.com/jtyers/tmaas-api-util/combo"
	"github.com/jtyers/tmaas-service-util/log"
)

func NewRouter(handlers *ThreatModelHandlers, comboFactory combo.ComboMiddlewareFactory, errorsMiddlewareFactory errors.ErrorsMiddlewareFactory, corsMiddleware corsconfig.CorsMiddleware) http.Handler {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware.Handler())

	r.Use(comboFactory.ExtractTokensToContext())

	r.Use(errorsMiddlewareFactory.NewErrorMiddleware(map[error]int{
		service.ErrNoSuchThreatModel: http.StatusNotFound,
		errors.ErrUnauthorized:       http.StatusUnauthorized,
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
	r.GET("/api/v1/threatModel/by-id/:threatModelID",
		comboFactory.StrictPermission(m.PermissionReadOwnThreatModels), // Permit service accounts to access this
		handlers.GetThreatModelHandler,
	)

	//r.DELETE("/api/v1/threatModel/by-id/:threatModelID",
	//	comboFactory.StrictUserPermission(e.PermissionWriteOwnThreatModel),
	//	handlers.DeleteThreatModelHandler,
	//)
	r.PATCH("/api/v1/threatModel/by-id/:threatModelID",
		comboFactory.StrictPermission(m.PermissionReadOwnThreatModels),
		handlers.PatchThreatModelHandler,
	)

	return r
}
