package service

import (
	"context"
	"fmt"

	m "github.com/jtyers/tmaas-model"
	"github.com/jtyers/tmaas-service-util/idchecker"
)

type ServiceThreatModelIDChecker struct {
	service ThreatModelService
}

func NewServiceThreatModelIDChecker(service ThreatModelService) *ServiceThreatModelIDChecker {
	return &ServiceThreatModelIDChecker{service}
}

var _ idchecker.IDCheckerForType = (*ServiceThreatModelIDChecker)(nil)

func (c *ServiceThreatModelIDChecker) CanHandle(id any) bool {
	_, ok := id.(m.ThreatModelID)
	return ok
}

func (c *ServiceThreatModelIDChecker) CheckID(ctx context.Context, id any) (bool, error) {
	threatModelID, ok := id.(m.ThreatModelID)
	if !ok {
		return false, fmt.Errorf("not a ThreatModelID")
	}

	_, err := c.service.GetThreatModel(ctx, threatModelID)
	if err == nil {
		return true, nil
	} else if err == ErrNoSuchThreatModel {
		return false, nil
	} else {
		return false, err
	}
}
