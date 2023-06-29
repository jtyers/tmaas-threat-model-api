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

// @Summary Retrieves threat models by threat model ID
// @Produce json
// @Param id path string true "The threat model ID to retrieve data for"
// @Security firebase
// @Success 200 {object} m.ThreatModel "The threat model data"
// @Failure 401 {string} string "If the token supplied is invalid, expired or does not have access to call this API."
// @Failure 404 {string} string "If the threat model ID does not exist or is not visible to this user."
// @Router /api/v1/threatmodel/{id} [get]
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

// @Summary Retrieves all threat models visible to the user
// @Produce json
// @Security firebase
// @Success 200 {array} m.ThreatModel "The threat model data"
// @Failure 401 {string} string "If the token supplied is invalid, expired or does not have access to call this API."
// @Router /api/v1/threatmodel [get]
func (th *ThreatModelHandlers) GetThreatModelsHandler(c *gin.Context) {
	result, err := th.threatModelService.GetAll(c)
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// @Summary Create a new ThreatModel
// @Accept json
// @Produce json
// @Param data body m.ThreatModelParams true "Parameters for the threat model to create"
// @Security firebase
// @Success 200 {object} m.ThreatModel "The created ThreatModel"
// @Failure 400 {string} string "If the threat model data supplied was invalid or badly formed, or any field failed validation (such as a missing required field or a value out of range), or an invalid ID supplied for any fields that accept IDs"
// @Failure 401 {string} string "If the token supplied is invalid, expired or does not have access to call this API"
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

// @Summary Update a ThreatModel
// @Accept json
// @Produce json
// @Param id path string true "The threat model ID to update"
// @Security firebase
// @Param data body m.ThreatModelParams true "The parameters containing fields to update"
// @Success 200 {object} m.ThreatModel "The (full) updated threat model data"
// @Failure 400 {string} string "If the threat model data supplied was invalid or badly formed, or any field failed validation (such as a missing required field or a value out of range), or an invalid ID supplied for any fields that accept IDs"
// @Failure 401 {string} string "If the token supplied is invalid, expired or does not have access to call this API"
// @Router /api/v1/threatmodel/{id} [patch]
func (th *ThreatModelHandlers) PatchThreatModelHandler(c *gin.Context) {
	threatModelIDStr := c.Param("threatModelID")
	threatModelID := m.NewThreatModelIDP(threatModelIDStr)

	var t m.ThreatModelParams

	err := c.BindJSON(&t)
	if err != nil {
		c.Error(err)
		return
	}

	updated, err := th.threatModelService.Update(c, threatModelID, t)
	if err != nil {
		c.Error(err)
		return
	}

	c.PureJSON(http.StatusOK, updated)
}

// @Summary Delete a ThreatModel by ID
// @Produce json
// @Param id path string true "ThreatModel ID"
// @Security firebase
// @Success 200 {string} string "Returned when the delete succeeds."
// @Failure 401 {string} string "If the token supplied is invalid, expired or does not have access to call this API."
// @Failure 404 {string} string "If the supplied threat model ID does not exist or is not visible to this user."
// @Router /api/v1/threatmodel/{id} [delete]
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
