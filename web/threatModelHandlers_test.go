package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/jtyers/tmaas-api-util/combo"
	"github.com/jtyers/tmaas-api-util/errors"
	cmocks "github.com/jtyers/tmaas-cors-config/mocks"
	m "github.com/jtyers/tmaas-model"
	"github.com/jtyers/tmaas-model/structs"
	"github.com/jtyers/tmaas-threat-model-api/service"
)

type msi map[string]interface{}

func createServer(comboFactory combo.ComboMiddlewareFactory, ts service.ThreatModelService) (*httptest.Server, func()) {
	errors := errors.NewDefaultErrorsMiddlewareFactory() // use real middleware to check error handling

	// use dummy CORS middleware
	corsMiddlware := cmocks.NewMockCorsMiddleware()

	// generate a test server so we can capture and inspect the request
	handlers := NewThreatModelHandlers(ts)
	testServer := httptest.NewServer(NewRouter(handlers, comboFactory, errors, corsMiddlware))

	gin.SetMode(gin.TestMode)
	closer := func() { testServer.Close() }
	return testServer, closer
}

func toJsonString(data any) string {
	result := &bytes.Buffer{}
	json.NewEncoder(result).Encode(data)

	return result.String()
}

func readToBytes(r io.Reader) []byte {
	responseBuf := &bytes.Buffer{}
	responseBuf.ReadFrom(r)

	return responseBuf.Bytes()
}

func readToString(r io.Reader) string {
	responseBuf := &bytes.Buffer{}
	responseBuf.ReadFrom(r)

	return responseBuf.String()
}

func TestGetThreatModelHandler(t *testing.T) {
	authorisedServiceAccount := "lookup-service-go"
	serviceAccountPermissionsJson := combo.ServiceAccountPermissionsJson(`{"` + authorisedServiceAccount + `": ["readOwnThreatModels"]}`)

	threatModel := m.ThreatModel{
		ThreatModelID: m.NewThreatModelIDP("1234-1234-1234-1234"),
		Title:         "my-first-threatModel",
	}

	var tests = []struct {
		name             string
		token            m.AuthenticationToken
		dsReturnValue    *m.ThreatModel
		dsReturnError    error
		expectedResponse int
		expectedBody     *m.ThreatModel // not checked if nil
	}{
		{
			"should get existing threatModel",
			&m.AuthenticationInfo{UserID: "u-12345678", Roles: []m.Role{m.RoleUser}},
			&threatModel,
			nil,
			http.StatusOK,
			&threatModel,
		},
		{
			"service token: should get existing threatModel",
			&m.ServiceAccountToken{Name: authorisedServiceAccount},
			&threatModel,
			nil,
			http.StatusOK,
			&threatModel,
		},
		{
			"service token: should return 404 for non-existent threatModel",
			&m.ServiceAccountToken{Name: authorisedServiceAccount},
			nil,
			service.ErrNoSuchThreatModel,
			http.StatusNotFound,
			nil,
		},
		{
			"should return 401 for unauthenticated users",
			nil, // no token
			nil, // <- both of these being nil means
			nil, // ThreatModelService call is not expected
			http.StatusUnauthorized,
			nil,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// given
			svc := service.NewMockThreatModelService(ctrl)

			if test.dsReturnValue != nil || test.dsReturnError != nil {
				svc.EXPECT().GetThreatModel(gomock.AssignableToTypeOf(&gin.Context{}), threatModel.ThreatModelID).Return(
					test.dsReturnValue, test.dsReturnError)
			}

			comboFactory := combo.NewMockComboMiddlewareFactoryWithTokensAndPermissions(ctrl, test.token, serviceAccountPermissionsJson)
			server, closeServer := createServer(comboFactory, svc)
			defer closeServer()

			// when
			request, _ := http.NewRequest(http.MethodGet,
				server.URL+UrlPrefix+"/"+threatModel.ThreatModelID.String(), nil)
			response, err := http.DefaultClient.Do(request)

			// then
			require.Nil(t, err)
			require.Equal(t, test.expectedResponse, response.StatusCode)

			if test.expectedBody != nil {
				got := m.ThreatModel{}
				body := readToBytes(response.Body)
				err = structs.JSONToStruct(body, &got)
				require.Nil(t, err)

				require.Equal(t, &got, test.expectedBody)
			}
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
		name                 string
		ai                   *m.AuthenticationInfo
		dsReturnValue        []*m.ThreatModel
		dsReturnError        error
		expectedResponse     int
		expectedResponseBody []*m.ThreatModel
	}{
		{
			"should get threatModels",
			&m.AuthenticationInfo{UserID: m.UserID("u-12345678"), Roles: []m.Role{m.RoleUser}},
			[]*m.ThreatModel{&threatModel1, &threatModel2, &threatModel3},
			nil,
			http.StatusOK,
			[]*m.ThreatModel{&threatModel1, &threatModel2, &threatModel3},
		},
		{
			"no threatModels return should yield empty array",
			&m.AuthenticationInfo{UserID: m.UserID("u-12345678"), Roles: []m.Role{m.RoleUser}},
			[]*m.ThreatModel{},
			nil,
			http.StatusOK,
			[]*m.ThreatModel{},
		},
		{
			"should return 401 if no token supplied",
			nil,
			[]*m.ThreatModel{},
			nil,
			http.StatusUnauthorized,
			nil,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// given
			mockThreatModelService := service.NewMockThreatModelService(ctrl)
			if test.ai != nil {
				mockThreatModelService.EXPECT().GetThreatModels(gomock.AssignableToTypeOf(&gin.Context{})).Return(test.dsReturnValue, test.dsReturnError)
			}

			comboFactory := combo.NewMockComboMiddlewareFactoryWithTokensAndPermissions(ctrl, test.ai,
				serviceAccountPermissionsJson)
			server, closeServer := createServer(comboFactory, mockThreatModelService)
			defer closeServer()

			// when
			request, _ := http.NewRequest(http.MethodGet, server.URL+UrlPrefix, nil)
			response, err := http.DefaultClient.Do(request)

			// then
			require.Nil(t, err)
			require.Equal(t, test.expectedResponse, response.StatusCode)

			if test.expectedResponseBody != nil {
				got := []*m.ThreatModel{}
				body := readToBytes(response.Body)
				err = json.Unmarshal(body, &got)
				require.Nil(t, err)

				require.Equal(t, test.expectedResponseBody, got)
			}
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
		expectedHttpResponse           int
	}{
		//
		// UNAUTHENTICATED TESTS
		//
		{
			"no authentication: should deny access",
			m.ThreatModelParams{Title: m.String("my new threatModel")},
			nil,
			nil,
			nil,
			nil,
			http.StatusUnauthorized,
		},
		//
		// AUTHENTICATED TESTS
		//
		{
			"should create threatModel and pass AuthenticationInfo",
			m.ThreatModelParams{Title: m.String("my new threatModel")},
			&m.AuthenticationInfo{UserID: m.UserID("u-1234-1234"), Roles: []m.Role{m.RoleUser}},
			&m.ThreatModel{Title: "my new threatModel"},
			&m.ThreatModel{Title: "my new threatModel"},
			nil,
			http.StatusOK,
		},
		{
			"should return an error if ThreatModelService returns an error",
			m.ThreatModelParams{Title: m.String("my new threatModel")},
			&m.AuthenticationInfo{UserID: m.UserID("u-1234-1234"), Roles: []m.Role{m.RoleUser}},
			&m.ThreatModel{Title: "my new threatModel"},
			nil,
			fmt.Errorf("some random error"),
			http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			mockThreatModelService := service.NewMockThreatModelService(ctrl)

			if test.expectedCreatedThreatModel != nil {
				mockThreatModelService.EXPECT().CreateThreatModel(gomock.Any(), test.input).Return(
					test.threatModelToReturnFromService, test.errorToReturnFromService)
			}

			comboFactory := combo.NewMockComboMiddlewareFactoryWithTokensAndPermissions(ctrl, test.ai,
				serviceAccountPermissionsJson)
			server, closeServer := createServer(comboFactory, mockThreatModelService)
			defer closeServer()

			// when
			request, _ := http.NewRequest(http.MethodPut,
				server.URL+UrlPrefix, strings.NewReader(toJsonString(test.input)))
			response, err := http.DefaultClient.Do(request)

			// then
			require.Nil(t, err)
			require.Equal(t, test.expectedHttpResponse, response.StatusCode)

			if test.expectedHttpResponse == http.StatusOK {
				d := m.ThreatModel{}
				err = structs.JSONToStruct(readToBytes(response.Body), &d)
				require.Nil(t, err)

				require.Equal(t, test.threatModelToReturnFromService, &d)
			}
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
		expectedResponse   int
	}{
		{
			"should update threatModel details",
			&m.AuthenticationInfo{UserID: "u-1234", Roles: []m.Role{m.RoleUser}},
			m.NewThreatModelIDP("d-1234"),
			m.ThreatModelParams{Title: m.String("foo")},
			nil,
			http.StatusOK,
		},
		{
			"should return 401 if no JWT supplied",
			nil,
			m.NewThreatModelIDP("d-1234"),
			m.ThreatModelParams{Title: m.String("foo")},
			nil,
			http.StatusUnauthorized,
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

			var bodyReader io.Reader = nil
			if test.ai != nil {
				mockThreatModelService.EXPECT().UpdateThreatModel(gomock.Any(), test.inputThreatModelID,
					test.input).Return(test.dsReturnError)

				s, err := structs.StructToJSON(test.input)
				require.Nil(t, err)
				bodyReader = strings.NewReader(s)
			}

			// when
			request, _ := http.NewRequest(http.MethodPatch,
				server.URL+UrlPrefix+"/"+test.inputThreatModelID.String(), bodyReader)
			response, err := http.DefaultClient.Do(request)

			// then
			require.Nil(t, err)
			require.Equal(t, test.expectedResponse, response.StatusCode)
		})
	}
}

func TestDeleteThreatModelHandler(t *testing.T) {
	serviceAccountPermissionsJson := combo.ServiceAccountPermissionsJson(`{}`)
	//
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	//
	var tests = []struct {
		name               string
		ai                 *m.AuthenticationInfo
		inputThreatModelID m.ThreatModelID
		dsReturnError      error
		expectedResponse   int
	}{
		{
			"should delete threatModel",
			&m.AuthenticationInfo{UserID: m.UserID("u-1"), Roles: []m.Role{m.RoleUser}},
			m.NewThreatModelIDP("d-12345678"),
			nil,
			http.StatusOK,
		},
		{
			"should return 401 when no token passed",
			nil,
			m.NewThreatModelIDP("d-12345678"),
			nil,
			http.StatusUnauthorized,
		},
	}
	//
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			mockThreatModelService := service.NewMockThreatModelService(ctrl)
			//
			comboFactory := combo.NewMockComboMiddlewareFactoryWithTokensAndPermissions(ctrl, test.ai,
				serviceAccountPermissionsJson)
			server, closeServer := createServer(comboFactory, mockThreatModelService)
			defer closeServer()
			//
			if test.ai != nil {
				mockThreatModelService.EXPECT().DeleteThreatModel(gomock.AssignableToTypeOf(&gin.Context{}), test.inputThreatModelID).Return(test.dsReturnError)
			}
			//
			// when
			request, _ := http.NewRequest(http.MethodDelete,
				server.URL+UrlPrefix+"/"+test.inputThreatModelID.String(), nil)
			response, err := http.DefaultClient.Do(request)
			//
			// then
			require.Nil(t, err)
			require.Equal(t, test.expectedResponse, response.StatusCode)
		})
	}
}
