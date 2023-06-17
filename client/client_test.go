package client

// Basically a copy of the handlers tests, but here we use the API client implementation
// to perform the calls, and thus are testing the end-to-end client/handler/service interaction.

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/jtyers/tmaas-api-util/combo"
	"github.com/jtyers/tmaas-api-util/errors"
	comocks "github.com/jtyers/tmaas-cors-config/mocks"
	m "github.com/jtyers/tmaas-model"
	"github.com/jtyers/tmaas-service-util/log"
	"github.com/jtyers/tmaas-service-util/requestor"
	"github.com/jtyers/tmaas-threat-model-api/service"
	"github.com/jtyers/tmaas-threat-model-api/web"
)

func createServer(comboFactory combo.ComboMiddlewareFactory, svc *service.MockThreatModelService) (*httptest.Server, func()) {
	log.InitialiseLogging()

	errors := errors.NewDefaultErrorsMiddlewareFactory() // use real middleware to check error handling

	// use dummy CORS middleware
	corsMiddlware := comocks.NewMockCorsMiddleware()

	// generate a test server so we can capture and inspect the request
	handlers := web.NewThreatModelHandlers(svc)
	testServer := httptest.NewServer(web.NewRouter(handlers, comboFactory, errors, corsMiddlware))

	gin.SetMode(gin.TestMode)
	closer := func() { testServer.Close() }
	return testServer, closer
}

func createClient(server *httptest.Server) *ThreatModelServiceClient {
	return &ThreatModelServiceClient{
		config:    ThreatModelServiceClientConfig{BaseURL: server.URL + "/"},
		requestor: requestor.NewDefaultRequestorWithContext(),
	}
}

func TestGetThreatModelHandler(t *testing.T) {
	authorisedServiceAccount := "lookup-service-go"
	serviceAccountPermissionsJson := combo.ServiceAccountPermissionsJson(`{"` + authorisedServiceAccount + `": ["readOwnThreatModels"]}`)

	threatModel := m.ThreatModel{
		ThreatModelID: m.NewThreatModelIDP("1234-1234-1234-1234"),
		Title:         "my-first-threatModel",
	}

	var tests = []struct {
		name          string
		token         m.AuthenticationToken
		dsReturnValue *m.ThreatModel
		dsReturnError error
		expectedBody  *m.ThreatModel // not checked if nil
		expectedError error
	}{
		{
			"should get existing threatModel",
			&m.AuthenticationInfo{UserID: "u-12345678", Roles: []m.Role{m.RoleUser}},
			&threatModel,
			nil,
			&threatModel,
			nil,
		},
		{
			"should return ErrNoSuchThreatModel for non-existent threatModel",
			&m.AuthenticationInfo{UserID: "u-12345678", Roles: []m.Role{m.RoleUser}},
			nil,
			service.ErrNoSuchThreatModel,
			nil,
			service.ErrNoSuchThreatModel,
		},
		{
			"service token: should get existing threatModel",
			&m.ServiceAccountToken{Name: authorisedServiceAccount},
			&threatModel,
			nil,
			&threatModel,
			nil,
		},
		{
			"service token: should return ErrNoSuchThreatModel for non-existent threatModel",
			&m.ServiceAccountToken{Name: authorisedServiceAccount},
			nil,
			service.ErrNoSuchThreatModel,
			nil,
			service.ErrNoSuchThreatModel,
		},
		{
			"should return 401 for unauthenticated users",
			nil, // no token
			nil, // <- both of these being nil means
			nil, // ThreatModelService call is not expected
			nil,
			requestor.ErrRequestFailed{http.StatusUnauthorized, ``},
		},
	}

	for _, test := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		t.Run(test.name, func(t *testing.T) {
			// given
			mockThreatModelService := service.NewMockThreatModelService(ctrl)

			if test.dsReturnValue != nil || test.dsReturnError != nil {
				mockThreatModelService.EXPECT().Get(gomock.AssignableToTypeOf(&gin.Context{}), threatModel.ThreatModelID).Return(
					test.dsReturnValue, test.dsReturnError)
			}

			comboFactory := combo.NewMockComboMiddlewareFactoryWithTokensAndPermissions(ctrl, test.token, serviceAccountPermissionsJson)
			server, closeServer := createServer(comboFactory, mockThreatModelService)
			defer closeServer()

			client := createClient(server)

			// when
			ctx := context.Background()
			response, err := client.Get(ctx, threatModel.ThreatModelID)

			// then
			require.Equal(t, test.expectedError, err)
			require.Equal(t, test.expectedBody, response)
		})
	}
}

func TestGetThreatModelsHandler(t *testing.T) {
	serviceAccountPermissionsJson := combo.ServiceAccountPermissionsJson(`{}`)

	threatModel1 := m.ThreatModel{
		ThreatModelID: m.NewThreatModelIDP("1234-1234-1234-1234"),
		Title:         "my-first-threatModel",
	}

	threatModel2 := m.ThreatModel{
		ThreatModelID: m.NewThreatModelIDP("2345-2345-2345-2345"),
		Title:         "my-second-threatModel",
	}

	threatModel3 := m.ThreatModel{
		ThreatModelID: m.NewThreatModelIDP("3456-3456-3456-3456"),
		Title:         "my-third-threatModel",
	}

	var tests = []struct {
		name             string
		ai               *m.AuthenticationInfo
		dsReturnValue    []*m.ThreatModel
		dsReturnError    error
		expectedResponse []*m.ThreatModel
		expectedError    error
	}{
		{
			"should get threatModels",
			&m.AuthenticationInfo{UserID: m.UserID("u-12345678"), Roles: []m.Role{&m.RoleUser}},
			[]*m.ThreatModel{&threatModel1, &threatModel2, &threatModel3},
			nil,
			[]*m.ThreatModel{&threatModel1, &threatModel2, &threatModel3},
			nil,
		},
		{
			"no threatModels return should yield empty array",
			&m.AuthenticationInfo{UserID: m.UserID("u-12345678"), Roles: []m.Role{&m.RoleUser}},
			[]*m.ThreatModel{},
			nil,
			[]*m.ThreatModel{},
			nil,
		},
		{
			"should return 401 if no token supplied",
			nil,
			[]*m.ThreatModel{},
			nil,
			nil,
			requestor.ErrRequestFailed{http.StatusUnauthorized, ``},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// given
			mockThreatModelService := service.NewMockThreatModelService(ctrl)
			if test.ai != nil {
				mockThreatModelService.EXPECT().GetAll(gomock.AssignableToTypeOf(&gin.Context{})).Return(test.dsReturnValue, test.dsReturnError)
			}

			comboFactory := combo.NewMockComboMiddlewareFactoryWithTokensAndPermissions(ctrl, test.ai,
				serviceAccountPermissionsJson)
			server, closeServer := createServer(comboFactory, mockThreatModelService)
			defer closeServer()

			client := createClient(server)

			// when
			ctx := context.Background()
			response, err := client.GetAll(ctx)

			// then
			require.Equal(t, test.expectedError, err)
			require.Equal(t, test.expectedResponse, response)
		})

	}
}

func TestCreateThreatModelHandler(t *testing.T) {
	serviceAccountPermissionsJson := combo.ServiceAccountPermissionsJson(`{}`)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var tests = []struct {
		name                           string
		input                          m.ThreatModelParams
		ai                             *m.AuthenticationInfo
		expectedCreatedThreatModel     *m.ThreatModel
		threatModelToReturnFromService *m.ThreatModel
		errorToReturnFromService       error
		expectedResponse               *m.ThreatModel
		expectedError                  error
	}{
		//
		// UNAUTHENTICATED TESTS
		//
		{
			"with authentication: should create threatModel and pass AuthenticationInfo",
			m.ThreatModelParams{Title: m.String("my new threatModel")},
			&m.AuthenticationInfo{UserID: m.UserID("u-1234-1234"), Roles: []m.Role{&m.RoleUser}},
			&m.ThreatModel{Title: "my new threatModel"},
			&m.ThreatModel{Title: "my new threatModel"},
			nil,
			&m.ThreatModel{Title: "my new threatModel"},
			nil,
		},
		{
			"with authentication: should return an error without msg if ThreatModelService returns a non-public error",
			m.ThreatModelParams{Title: m.String("my new threatModel")},
			&m.AuthenticationInfo{UserID: m.UserID("u-1234-1234"), Roles: []m.Role{&m.RoleUser}},
			&m.ThreatModel{Title: "my new threatModel"},
			nil,
			fmt.Errorf("some random error"),
			nil,
			requestor.ErrRequestFailed{http.StatusInternalServerError, ``},
		},
		{
			"no authentication: should fail",
			m.ThreatModelParams{Title: m.String("my new threatModel")},
			nil,
			nil,
			nil,
			errors.ErrUnauthorized,
			nil,
			requestor.ErrRequestFailed{http.StatusUnauthorized, ``},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			mockThreatModelService := service.NewMockThreatModelService(ctrl)

			if test.expectedCreatedThreatModel != nil {
				mockThreatModelService.EXPECT().Create(gomock.Any(), test.input).Return(
					test.threatModelToReturnFromService, test.errorToReturnFromService)
			}

			comboFactory := combo.NewMockComboMiddlewareFactoryWithTokensAndPermissions(ctrl, test.ai,
				serviceAccountPermissionsJson)
			server, closeServer := createServer(comboFactory, mockThreatModelService)
			defer closeServer()

			client := createClient(server)

			// when
			ctx := context.Background()
			response, err := client.Create(ctx, test.input)

			// then
			require.Equal(t, test.expectedError, err)
			require.Equal(t, test.expectedResponse, response)
		})
	}
}

func TestPatchThreatModelHandler(t *testing.T) {
	serviceAccountPermissionsJson := combo.ServiceAccountPermissionsJson(`{}`)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var tests = []struct {
		name               string
		ai                 *m.AuthenticationInfo
		inputThreatModelID m.ThreatModelID
		input              m.ThreatModelParams
		dsReturnError      error
		expectedError      error
	}{
		{
			"should update threatModel details",
			&m.AuthenticationInfo{UserID: "u-1234", Roles: []m.Role{&m.RoleUser}},
			m.NewThreatModelIDP("d-1234"),
			m.ThreatModelParams{Title: m.String("foo")},
			nil,
			nil,
		},
		{
			"should return 401 if no JWT supplied",
			nil,
			m.NewThreatModelIDP("d-1234"),
			m.ThreatModelParams{Title: m.String("foo")},
			nil,
			requestor.ErrRequestFailed{http.StatusUnauthorized, ``},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			mockThreatModelService := service.NewMockThreatModelService(ctrl)

			comboFactory := combo.NewMockComboMiddlewareFactoryWithTokensAndPermissions(ctrl, test.ai,
				serviceAccountPermissionsJson)
			server, closeServer := createServer(comboFactory, mockThreatModelService)
			defer closeServer()

			if test.ai != nil {
				mockThreatModelService.EXPECT().Update(gomock.Any(), test.inputThreatModelID,
					test.input).Return(test.dsReturnError)
			}

			client := createClient(server)

			// when
			ctx := context.Background()
			err := client.Update(ctx, test.inputThreatModelID, test.input)

			// then
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestDeleteThreatModelHandler(t *testing.T) {
	serviceAccountPermissionsJson := combo.ServiceAccountPermissionsJson(`{}`)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var tests = []struct {
		name               string
		ai                 *m.AuthenticationInfo
		inputThreatModelID m.ThreatModelID
		dsReturnError      error
		expectedError      error
	}{
		{
			"should delete threatModel",
			&m.AuthenticationInfo{UserID: m.UserID("u-1"), Roles: []m.Role{&m.RoleUser}},
			m.NewThreatModelIDP("d-12345678"),
			nil,
			nil,
		},
		{
			"should return 401 when no token passed",
			nil,
			m.NewThreatModelIDP("d-12345678"),
			nil,
			requestor.ErrRequestFailed{http.StatusUnauthorized, ""},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			mockThreatModelService := service.NewMockThreatModelService(ctrl)

			comboFactory := combo.NewMockComboMiddlewareFactoryWithTokensAndPermissions(ctrl, test.ai,
				serviceAccountPermissionsJson)
			server, closeServer := createServer(comboFactory, mockThreatModelService)
			defer closeServer()

			if test.ai != nil {
				mockThreatModelService.EXPECT().Delete(gomock.AssignableToTypeOf(&gin.Context{}), test.inputThreatModelID).Return(test.dsReturnError)
			}

			client := createClient(server)

			// when
			ctx := context.Background()
			err := client.Delete(ctx, test.inputThreatModelID)

			// then
			require.Equal(t, test.expectedError, err)
		})
	}
}
