package client

import (
	"context"

	m "github.com/jtyers/tmaas-model"
	"github.com/jtyers/tmaas-service-util/idchecker"
	"github.com/jtyers/tmaas-threat-model-api/service"
)

type ClientThreatModelIDChecker struct {
	client *ThreatModelServiceClient
}

func NewClientThreatModelIDChecker(client *ThreatModelServiceClient) *ClientThreatModelIDChecker {
	return &ClientThreatModelIDChecker{client}
}

var _ idchecker.IDCheckerForType = (*ClientThreatModelIDChecker)(nil)

func (c *ClientThreatModelIDChecker) CanHandle(id any) bool {
	switch id.(type) {
	case m.ThreatModelID, *m.ThreatModelID:
		return true
	}
	return false
}

func (c *ClientThreatModelIDChecker) CheckID(ctx context.Context, id any) (bool, error) {
	var idStruct m.ThreatModelID
	switch id.(type) {
	case m.ThreatModelID:
		idStruct = id.(m.ThreatModelID)
	case *m.ThreatModelID:
		idStruct = *(id.(*m.ThreatModelID))
	}

	_, err := c.client.Get(ctx, idStruct)
	if err == nil {
		return true, nil

	} else if err == service.ErrNoSuchThreatModel {
		// FIXME will never match, as it'll turn up as a RequestFailedError
		return false, nil

	} else {
		return false, err
	}
}
