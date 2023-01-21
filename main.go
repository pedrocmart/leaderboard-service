package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pedrocmart/leaderboard-service/coreservices"
	httpHandlers "github.com/pedrocmart/leaderboard-service/handlers"
	"github.com/pedrocmart/leaderboard-service/models"
	"github.com/pedrocmart/leaderboard-service/utils"

	_ "modernc.org/ql/driver"
)

func main() {
	//first we will initialize core that needs a response writer attatched to it
	core = coreservices.InitCore()
	core.ConnectResponseWriter()

	coreservices.NewCoreService(core)

	createsInMemoryDB()
	prepareConnectHTTP()
}

func prepareConnectHTTP() {
	httpHandlers.ConnectBasic(router, core)

	fmt.Printf("Listening and serving on Host: %s, Port: %s\n", host, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func createsInMemoryDB() {
	//defining in memory database
	mdb, err := sql.Open("ql-mem", "memory://mem.db")
	if err != nil {
		log.Fatal(err)
	}
	coreservices.NewStoreService(core, mdb)
	tx, err := mdb.Begin()
	if err != nil {
		return
	}

	if _, err := tx.Exec("CREATE TABLE users (id INT, score INT);"); err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}
}

var core *models.Core
var router = mux.NewRouter()
var host = utils.GetEnvOrDefault("HOST", "0.0.0.0")
var port = utils.GetEnvOrDefault("PORT", "8894")
