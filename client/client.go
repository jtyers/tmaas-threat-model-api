package client

import (
	"context"
	"fmt"

	m "github.com/jtyers/tmaas-model"
	"github.com/jtyers/tmaas-service-util/requestor"
	"github.com/jtyers/tmaas-threat-model-api/service"
)

var (
	URLPrefix       = "%sapi/v1/threatmodel"
	URLPrefixWithID = URLPrefix + "/%s"
)

type ThreatModelServiceClientConfig struct {
	// The base URL for API requests.
	BaseURL string
}

// A client for ThreatModelService that makes calls over HTTPS.
type ThreatModelServiceClient struct {
	config    ThreatModelServiceClientConfig
	requestor requestor.RequestorWithContext
}

var _ service.ThreatModelService = (*ThreatModelServiceClient)(nil)

func NewThreatModelServiceClient(
	config ThreatModelServiceClientConfig,
	requestor requestor.RequestorWithContext,
) *ThreatModelServiceClient {
	return &ThreatModelServiceClient{config, requestor}
}

// Retrieve a ThreatModel by ID.
func (s *ThreatModelServiceClient) Get(ctx context.Context, id m.ThreatModelID) (*m.ThreatModel, error) {
	result := m.ThreatModel{}
	err := s.requestor.GetInto(ctx, fmt.Sprintf(URLPrefixWithID, s.config.BaseURL, id.String()), &result)
	if err != nil {
		if reqErr, ok := err.(requestor.ErrRequestFailed); ok && reqErr.StatusCode == 404 {
			return nil, service.ErrNoSuchThreatModel
		}
		return nil, err
	}

	return &result, nil
}

// Retrieve all ThreatModels.
func (s *ThreatModelServiceClient) GetAll(ctx context.Context) ([]*m.ThreatModel, error) {
	result := []*m.ThreatModel{}
	err := s.requestor.GetInto(ctx, fmt.Sprintf(URLPrefix, s.config.BaseURL), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *ThreatModelServiceClient) Query(ctx context.Context, q *m.ThreatModelQuery) ([]*m.ThreatModel, error) {
	// not yet implemented by the ThreatModelService's web router
	return nil, fmt.Errorf("not yet implemented")
}

func (s *ThreatModelServiceClient) QuerySingle(ctx context.Context, q *m.ThreatModelQuery) (*m.ThreatModel, error) {
	// not yet implemented by the ThreatModelService's web router
	return nil, fmt.Errorf("not yet implemented")
}

func (s *ThreatModelServiceClient) GetThreats(ctx context.Context, id m.ThreatModelID) ([]*m.Threat, error) {
	// not yet implemented by the ThreatModelService's web router
	return nil, fmt.Errorf("not yet implemented")
}

// Creates a ThreatModel.
func (s *ThreatModelServiceClient) Create(ctx context.Context, params m.ThreatModelParams) (*m.ThreatModel, error) {
	body, err := requestor.StructReader(params)
	if err != nil {
		return nil, err
	}

	result := m.ThreatModel{}
	err = s.requestor.PutInto(ctx, fmt.Sprintf(URLPrefix, s.config.BaseURL), body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Updates a ThreatModel.
func (s *ThreatModelServiceClient) Update(ctx context.Context, id m.ThreatModelID, params m.ThreatModelParams) error {
	body, err := requestor.StructReader(params)
	if err != nil {
		if reqErr, ok := err.(requestor.ErrRequestFailed); ok && reqErr.StatusCode == 404 {
			return service.ErrNoSuchThreatModel
		}
		return err
	}

	_, err = s.requestor.Patch(ctx, fmt.Sprintf(URLPrefixWithID, s.config.BaseURL, id.String()), body)
	if err != nil {
		return err
	}

	return nil
}

// Delete a ThreatModel by ID..
func (s *ThreatModelServiceClient) Delete(ctx context.Context, id m.ThreatModelID) error {
	_, err := s.requestor.Delete(ctx, fmt.Sprintf(URLPrefixWithID, s.config.BaseURL, id.String()))
	if err != nil {
		if reqErr, ok := err.(requestor.ErrRequestFailed); ok && reqErr.StatusCode == 404 {
			return service.ErrNoSuchThreatModel
		}
		return err
	}

	return nil
}
