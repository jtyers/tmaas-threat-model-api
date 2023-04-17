// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/jtyers/tmaas-api-util/combo"
	"github.com/jtyers/tmaas-api-util/errors"
	extractor2 "github.com/jtyers/tmaas-api-util/firebase/extractor"
	"github.com/jtyers/tmaas-api-util/serviceaccount/extractor"
	"github.com/jtyers/tmaas-api/dao"
	"github.com/jtyers/tmaas-api/service"
	"github.com/jtyers/tmaas-api/web"
	"github.com/jtyers/tmaas-cors-config"
	"github.com/jtyers/tmaas-model/validator"
	"github.com/jtyers/tmaas-service-dao/clover"
	"github.com/jtyers/tmaas-service-util/id"
	dao2 "github.com/jtyers/tmaas-threat-service/dao"
	"net/http"
)

// Injectors from wire.go:

func InitialiseRouter() (http.Handler, error) {
	databaseConfig := dao.NewCloverDatabaseConfig()
	db, err := clover.NewCloverDB(databaseConfig)
	if err != nil {
		return nil, err
	}
	threatModelCloverCollectionConfig := dao2.NewThreatModelCloverCollectionConfig()
	cloverDao := dao2.NewThreatModelCloverDao(db, threatModelCloverCollectionConfig)
	threatModelDao := dao2.NewThreatModelDao(cloverDao)
	defaultRandomIDProvider := id.NewDefaultRandomIDProvider()
	defaultStructValidator, err := validator.NewDefaultStructValidator()
	if err != nil {
		return nil, err
	}
	defaultThreatModelService := service.NewDefaultThreatModelService(threatModelDao, defaultRandomIDProvider, defaultStructValidator)
	threatModelHandlers := web.NewThreatModelHandlers(defaultThreatModelService)
	context := clover.NewContext()
	iamClient, err := extractor.NewIamClient(context)
	if err != nil {
		return nil, err
	}
	defaultVerifier := extractor.NewDefaultVerifier(iamClient)
	defaultExtractor := extractor.NewDefaultExtractor(defaultVerifier)
	app, err := extractor2.NewFirebaseApp(context)
	if err != nil {
		return nil, err
	}
	client, err := extractor2.NewFirebaseAuthClient(context, app)
	if err != nil {
		return nil, err
	}
	defaultFirebaseVerifier := extractor2.NewDefaultFirebaseVerifier(client)
	defaultFirebaseExtractor := extractor2.NewDefaultFirebaseExtractor(defaultFirebaseVerifier)
	serviceAccountPermissionsJson := combo.NewServiceAccountPermissionsJson()
	defaultComboMiddlewareFactory, err := combo.NewDefaultComboMiddlewareFactory(defaultExtractor, defaultFirebaseExtractor, serviceAccountPermissionsJson)
	if err != nil {
		return nil, err
	}
	defaultErrorsMiddlewareFactory := errors.NewDefaultErrorsMiddlewareFactory()
	corsMiddleware := corsconfig.FromEnv()
	handler := web.NewRouter(threatModelHandlers, defaultComboMiddlewareFactory, defaultErrorsMiddlewareFactory, corsMiddleware)
	return handler, nil
}
