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
