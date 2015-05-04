package controllers

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/nickjanus/ProteinGraphQuery/app/models"
	"github.com/revel/revel"
	"log"
	"strconv"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Query(query, cutoff string) revel.Result {
	var queryResult []byte
	var results []models.GraphEntry
	resultsAboveCutoff := make([]models.GraphEntry, 0)
	n, err := strconv.ParseFloat(cutoff, 64)
	if err != nil {
		log.Fatal("Could not parse cutoff argument")
	}

	db, err := bolt.Open(models.DatabaseName, 0600, nil)
	defer db.Close()
	if err != nil {
		log.Fatal("Could not open database!")
	}

	db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(models.EntryBucketName))
		queryResult = b.Get([]byte(query))
		return nil
	})

	if err := json.Unmarshal(queryResult, &results); err != nil {
		log.Println("Could not decode database entry!")
	}

	for _, entry := range results {
		if entry.MakesCutoff(n) {
			resultsAboveCutoff = append(resultsAboveCutoff, entry)
		}
	}
	result, err := json.Marshal(resultsAboveCutoff)
	if err != nil {
		log.Fatal("Cannot encode resulting graph")
	}
	graph := resultsAboveCutoff
	//graph := string(result)
	log.Println(string(result))

	return c.Render(graph)
}
