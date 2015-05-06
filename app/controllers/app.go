package controllers

import (
	"encoding/json"
	"github.com/nickjanus/ProteinGraphQuery/app"
	"github.com/nickjanus/ProteinGraphQuery/app/models"
	"github.com/revel/revel"
	"log"
	"regexp"
	"strconv"
	"strings"
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

	graph := models.RetrieveSubgraph(app.DB, queries, n)
	graphEnc, err := json.Marshal(graph)
	if err != nil {
		log.Fatal("Cannot encode resulting graph")
	}
	jsonGraph := string(graphEnc)

	return c.Render(graph, jsonGraph)
}
