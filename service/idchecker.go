package service

import (
	"context"

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
	switch id.(type) {
	case m.ThreatModelID, *m.ThreatModelID:
		return true
	}
	return false
}

func (c *ServiceThreatModelIDChecker) CheckID(ctx context.Context, id any) (bool, error) {
	var idStruct m.ThreatModelID
	switch id.(type) {
	case m.ThreatModelID:
		idStruct = id.(m.ThreatModelID)
	case *m.ThreatModelID:
		idStruct = *(id.(*m.ThreatModelID))
	}

	_, err := c.service.GetThreatModel(ctx, idStruct)
	if err == nil {
		return true, nil
	} else if err == ErrNoSuchThreatModel {
		return false, nil
	} else {
		return false, err
	}
}
