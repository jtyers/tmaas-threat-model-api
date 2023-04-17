package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jtyers/tmaas-api/service"
	m "github.com/jtyers/tmaas-model"
	"github.com/jtyers/tmaas-model/structs"
)

type ThreatModelHandlers struct {
	threatModelService service.ThreatModelService
}

func NewThreatModelHandlers(ts service.ThreatModelService) *ThreatModelHandlers {
	return &ThreatModelHandlers{threatModelService: ts}
}

// Retrieves a threatModel by ID.
func (th *ThreatModelHandlers) GetThreatModelHandler(c *gin.Context) {
	threatModelIDStr := c.Param("threatModelID")
	threatModelID := m.ThreatModelID(threatModelIDStr)

	result, err := th.threatModelService.GetThreatModel(c, threatModelID)
	if err != nil {
		c.Error(err)
	} else {
		c.PureJSON(http.StatusOK, structs.StructToMap(result))
	}
}

func (th *ThreatModelHandlers) GetThreatModelsHandler(c *gin.Context) {
	result, err := th.threatModelService.GetThreatModels(c)
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func (th *ThreatModelHandlers) PutThreatModelHandler(c *gin.Context) {
	var t m.ThreatModel

	err := c.BindJSON(&t)
	if err != nil {
		c.Error(err)
		return
	}

	result, err := th.threatModelService.CreateThreatModel(c, t)
	if err != nil {
		c.Error(err)
		return
	}

	c.PureJSON(http.StatusOK, structs.StructToMap(result))
}

func (th *ThreatModelHandlers) PatchThreatModelHandler(c *gin.Context) {
	threatModelIDStr := c.Param("threatModelID")
	threatModelID := m.ThreatModelID(threatModelIDStr)

	var t m.ThreatModel

	err := c.BindJSON(&t)
	if err != nil {
		c.Error(err)
		return
	}

	err = th.threatModelService.UpdateThreatModel(c, threatModelID, t)
	if err != nil {
		c.Error(err)
		return
	}

	c.PureJSON(http.StatusOK, structs.StructToMap(t))
}

// func (th *ThreatModelHandlers) DeleteThreatModelHandler(c *gin.Context) {
// 	threatModelIDStr := c.Param("threatModelID")
// 	threatModelID := m.ThreatModelID(threatModelIDStr)
//
// 	err := th.threatModelService.DeleteThreatModel(c, threatModelID)
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}
// }
