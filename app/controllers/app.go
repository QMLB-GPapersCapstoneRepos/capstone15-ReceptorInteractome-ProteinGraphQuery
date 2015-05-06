package controllers

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/nickjanus/ProteinGraphQuery/app/models"
	"github.com/revel/revel"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Query(query, cutoff string) revel.Result {
	query = strings.Trim(query, "\t ")
	c.Validation.Required(query).Message("Please enter a query.")
	c.Validation.Required(cutoff).Message("Please specify a cutoff")
	c.Validation.Match(query, regexp.MustCompile("^[[:word:]]+(,[[:space:]]*[[:word:]]+)*$")).
		Message("Please separate multiple items with commas.")
	c.Validation.Match(cutoff, regexp.MustCompile("^[-.0-9]+$")).
		Message("Cutoff must be specified as a number.")
	n, _ := strconv.ParseFloat(cutoff, 64)
	if n < -0.5 {
		c.Validation.Error("Cutoff must be greater than -0.5")
	}

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.Index)
	}
	query = strings.Replace(query, " ", "", -1)
	queries := strings.Split(query, ",")

	//TODO share db connection between requests
	db, err := bolt.Open(models.DatabaseName, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	graph := models.RetrieveSubgraph(db, queries, n)
	graphEnc, err := json.Marshal(graph)
	if err != nil {
		log.Fatal("Cannot encode resulting graph")
	}
	jsonGraph := string(graphEnc)

	return c.Render(graph, jsonGraph)
}
