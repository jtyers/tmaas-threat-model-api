package client

import (
	"context"
	"fmt"

	m "github.com/jtyers/tmaas-model"
	"github.com/jtyers/tmaas-service-util/idchecker"
	"github.com/jtyers/tmaas-service-util/requestor"
	"github.com/jtyers/tmaas-threat-model-api/service"
)

var (
	URLPrefix = "api/v1/threatmodel"
)

type ThreatModelServiceClientConfig struct {
	// The base URL for API requests.
	BaseURL string
}

// A client for DeviceService that makes calls over HTTPS.
type ThreatModelServiceClient struct {
	config    ThreatModelServiceClientConfig
	requestor requestor.RequestorWithContext
}

var _ service.ThreatModelService = (*ThreatModelServiceClient)(nil)

func NewThreatModelServiceClient(
	config ThreatModelServiceClientConfig,
	requestor requestor.RequestorWithContext,
	idChecker idchecker.IDChecker,
) *ThreatModelServiceClient {
	return &ThreatModelServiceClient{config, requestor}
}

// Retrieve a threatModel by threatModelID.
func (s *ThreatModelServiceClient) GetThreatModel(ctx context.Context, id m.ThreatModelID) (*m.ThreatModel, error) {
	result := m.ThreatModel{}
	err := s.requestor.GetInto(ctx, s.config.BaseURL+URLPrefix+"/"+id.String(), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Retrieve all threatModel.
func (s *ThreatModelServiceClient) GetThreatModels(ctx context.Context) ([]*m.ThreatModel, error) {
	result := []*m.ThreatModel{}
	err := s.requestor.GetInto(ctx, s.config.BaseURL+URLPrefix, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *ThreatModelServiceClient) GetThreats(ctx context.Context, id m.ThreatModelID) ([]*m.Threat, error) {
	// not yet implemented by the ThreatModelService's web router
	return nil, fmt.Errorf("not yet implemented")

	//result := []*m.Threat{}
	//err := s.requestor.GetInto(ctx, s.config.BaseURL+URLPrefix+"/", &result)
	//if err != nil {
	//	return nil, err
	//}

	// return result, nil
}

// Creates a threatModel. `threatModel` should not have ID or threatModelID set.
func (s *ThreatModelServiceClient) CreateThreatModel(ctx context.Context, threatModel m.ThreatModel) (*m.ThreatModel, error) {
	body, err := requestor.StructReader(threatModel)
	if err != nil {
		return nil, err
	}

	result := m.ThreatModel{}
	err = s.requestor.PutInto(ctx, s.config.BaseURL+URLPrefix, body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Updates a threatModel
func (s *ThreatModelServiceClient) UpdateThreatModel(ctx context.Context, threatModelID m.ThreatModelID, threatModel m.ThreatModel) error {
	body, err := requestor.StructReader(threatModel)
	if err != nil {
		return err
	}

	_, err = s.requestor.Patch(ctx, s.config.BaseURL+URLPrefix+"/"+threatModelID.String(), body)
	if err != nil {
		return err
	}

	return nil
}

// Delete a threatModel by threatModelID.
func (s *ThreatModelServiceClient) DeleteThreatModel(ctx context.Context, id m.ThreatModelID) error {
	_, err := s.requestor.Delete(ctx, s.config.BaseURL+URLPrefix+"/"+id.String())
	if err != nil {
		return err
	}

	return nil
}
