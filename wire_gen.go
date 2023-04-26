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
	"github.com/jtyers/tmaas-service-dao/firestore"
	"github.com/jtyers/tmaas-service-util/id"
	"net/http"
)

// Injectors from wire.go:

func InitialiseRouter() (http.Handler, error) {
	configuration := dao.NewConfiguration()
	firestoreConfiguration := dao.NewFirestoreConfiguration(configuration)
	context := firestore.NewContext()
	client, err := firestore.NewFirestoreClient(context, firestoreConfiguration)
	if err != nil {
		return nil, err
	}
<<<<<<< HEAD
	threatModelCloverCollectionConfig := dao.NewThreatModelCloverCollectionConfig()
	cloverDao := dao.NewThreatModelCloverDao(db, threatModelCloverCollectionConfig)
	threatModelDao := dao.NewThreatModelDao(cloverDao)
=======
	collectionRef := firestore.NewFirestoreCollection(firestoreConfiguration, client)
	firestoreDao := firestore.NewFirestoreDao(collectionRef, client)
	threatModelDao := dao.NewThreatModelDao(firestoreDao)
>>>>>>> 97775591bb8c4e15c8beb586a8ebe9752be47237
	defaultRandomIDProvider := id.NewDefaultRandomIDProvider()
	defaultStructValidator, err := validator.NewDefaultStructValidator()
	if err != nil {
		return nil, err
	}
	defaultThreatModelService := service.NewDefaultThreatModelService(threatModelDao, defaultRandomIDProvider, defaultStructValidator)
	threatModelHandlers := web.NewThreatModelHandlers(defaultThreatModelService)
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
	authClient, err := extractor2.NewFirebaseAuthClient(context, app)
	if err != nil {
		return nil, err
	}
	defaultFirebaseVerifier := extractor2.NewDefaultFirebaseVerifier(authClient)
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
