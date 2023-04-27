// @title ThreatPlane Threat Model API
// @version 1.0
// @description The API used to interact with Threat Models

// @contact.name ThreatPlane
// @contact.url http://www.threatplane.io
// @contact.email support@threatplane.io

// @license.name Commercial licence
// @license.url https://threatplane.io/terms

// @host localhost:8080
// @BasePath /
// @query.collection.format multi

package main

import (
	"net/http"

	util "github.com/jtyers/tmaas-service-util"
	log "github.com/jtyers/tmaas-service-util/log"
)

func main() {
	log.InitialiseLogging()

	port := util.GetEnv("PORT")
	r, err := InitialiseRouter()
	if err != nil {
		log.Fatalf("error while initialising router: %v", err)
	}

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("error while serving: %v", err)
	}
}
