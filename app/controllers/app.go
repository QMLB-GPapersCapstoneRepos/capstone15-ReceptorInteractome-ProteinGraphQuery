package controllers

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/nickjanus/ProteinGraphQuery/app/models"
	"github.com/revel/revel"
	"log"
	"strconv"
	"time"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Query(query, cutoff string) revel.Result {
	n, err := strconv.ParseFloat(cutoff, 64)
	if err != nil {
		log.Fatal("Could not parse cutoff argument")
	}

	//TODO share db connection between requests
	db, err := bolt.Open(models.DatabaseName, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	graph := models.RetrieveSubgraph(db, query, n)
	result, err := json.Marshal(graph)
	if err != nil {
		log.Fatal("Cannot encode resulting graph")
	}
	log.Println(string(result))

	return c.Render(graph)
}
