package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	m "github.com/jtyers/tmaas-model"
	"github.com/jtyers/tmaas-threat-model-api/service"
)

type ThreatModelHandlers struct {
	threatModelService service.ThreatModelService
}

func NewThreatModelHandlers(ts service.ThreatModelService) *ThreatModelHandlers {
	return &ThreatModelHandlers{threatModelService: ts}
}

// @Summary Retrieves a threatModel by ID.
// @ID get-threat-model-by-id
// @Produce json
// @Param threatModelID path string true "threat model ID"
// @Success 200 {object} m.ThreatModel
// @Failure 404 {string} string "Not Found"
// @Router /api/v1/threatmodel/{threatModelID} [get]
func (th *ThreatModelHandlers) GetThreatModelHandler(c *gin.Context) {
	threatModelIDStr := c.Param("threatModelID")
	threatModelID := m.NewThreatModelIDP(threatModelIDStr)

	result, err := th.threatModelService.Get(c, threatModelID)
	if err != nil {
		c.Error(err)
	} else {
		c.PureJSON(http.StatusOK, result)
	}
}

// @Summary Get all Threat Models
// @ID get-all-threat-models
// @Produce json
// @Success 200 {object} []m.ThreatModel
// @Router /api/v1/threatmodel [get]
func (th *ThreatModelHandlers) GetThreatModelsHandler(c *gin.Context) {
	result, err := th.threatModelService.GetAll(c)
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// @Summary Create a new Threat Model
// @ID create-threat-model
// @Produce json
// @Param data body m.ThreatModel true "todo data"
// @Success 200 {object} m.ThreatModel
// @Failure 400 {string} string "Nope"
// @Router /api/v1/threatmodel [put]
func (th *ThreatModelHandlers) PutThreatModelHandler(c *gin.Context) {
	var t m.ThreatModelParams

	err := c.BindJSON(&t)
	if err != nil {
		c.Error(err)
		return
	}

	result, err := th.threatModelService.Create(c, t)
	if err != nil {
		c.Error(err)
		return
	}

	c.PureJSON(http.StatusOK, result)
}

// @Summary Update a new Threat Model
// @ID update-threat-model
// @Produce json
// @Param data body m.ThreatModel true "todo data"
// @Success 200 {object} m.ThreatModel
// @Failure 400 {string} string "Nope"
// @Router /api/v1/threatmodel/{threatModelID} [patch]
func (th *ThreatModelHandlers) PatchThreatModelHandler(c *gin.Context) {
	threatModelIDStr := c.Param("threatModelID")
	threatModelID := m.NewThreatModelIDP(threatModelIDStr)

	var t m.ThreatModelParams

	err := c.BindJSON(&t)
	if err != nil {
		c.Error(err)
		return
	}

	err = th.threatModelService.Update(c, threatModelID, t)
	if err != nil {
		c.Error(err)
		return
	}

	c.PureJSON(http.StatusOK, t)
}

// @Summary Delete a threat model
// @ID delete-threat-model
// @Produce json
// @Param threatModelID path string true "threat model ID"
// @Success 200 {object} m.ThreatModel
// @Failure 404 {string} string "Not Found"
// @Router /api/v1/threatmodel/{threatModelID} [delete]
func (th *ThreatModelHandlers) DeleteThreatModelHandler(c *gin.Context) {
	threatModelIDStr := c.Param("threatModelID")
	threatModelID := m.NewThreatModelIDP(threatModelIDStr)
	//
	err := th.threatModelService.Delete(c, threatModelID)
	if err != nil {
		c.Error(err)
		return
	}
}
