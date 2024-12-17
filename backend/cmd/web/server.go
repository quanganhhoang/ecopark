package main

import (
	"fmt"
	"log"
	"net/http"
)

func (app *App) RunServer() {
	app.Logger.Info("Starting server", "app.Port", app.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", app.Port), app.Router))
}