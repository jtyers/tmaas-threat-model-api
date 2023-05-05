package client

import (
	"context"
	"fmt"

	m "github.com/jtyers/tmaas-model"
	"github.com/jtyers/tmaas-service-util/idchecker"
	"github.com/jtyers/tmaas-threat-model-api/service"
)

type ClientThreatModelIDChecker struct {
	client ThreatModelServiceClient
}

func NewClientThreatModelIDChecker(client ThreatModelServiceClient) *ClientThreatModelIDChecker {
	return &ClientThreatModelIDChecker{client}
}

var _ idchecker.IDCheckerForType = (*ClientThreatModelIDChecker)(nil)

func (c *ClientThreatModelIDChecker) CanHandle(id any) bool {
	_, ok := id.(m.ThreatModelID)
	return ok
}

func (c *ClientThreatModelIDChecker) CheckID(ctx context.Context, id any) (bool, error) {
	threatModelID, ok := id.(m.ThreatModelID)
	if !ok {
		return false, fmt.Errorf("not a ThreatModelID")
	}

	_, err := c.client.GetThreatModel(ctx, threatModelID)
	if err == nil {
		return true, nil

	} else if err == service.ErrNoSuchThreatModel {
		// FIXME will never match, as it'll turn up as a RequestFailedError
		return false, nil

	} else {
		return false, err
	}
}
