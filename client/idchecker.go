package client

import (
	"context"
	"fmt"

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

// ClientThreatModelIDCheckerIniter is a crap name for an 'initer', which
// registers the checker. It exists because we cannot call initialisation
// code as part of wire directly, so we use wire to pull this in as call it
// as part of main().
type ClientThreatModelIDCheckerIniter struct {
	idChecker idchecker.IDChecker
	checker   *ClientThreatModelIDChecker
}

func NewClientThreatModelIDCheckerIniter(
	idChecker idchecker.IDChecker,
	checker *ClientThreatModelIDChecker,
) *ClientThreatModelIDCheckerIniter {
	return &ClientThreatModelIDCheckerIniter{idChecker, checker}
}

// Register registers the ClientThreatModelIDChecker with the IDChecker
func (i *ClientThreatModelIDCheckerIniter) Register() error {
	return i.idChecker.RegisterIDChecker(i.checker)
}
