package main

import (
	"log"

	"github.com/jack5341/giggle-hoster/internal/database"
	"github.com/jack5341/giggle-hoster/internal/types"
)

func main() {
	db, err := database.EstablishDBConnection()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(types.Node{}, types.Pod{}); err != nil {
		log.Fatal(err)
	}

	database.Db = db
}
