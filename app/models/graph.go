package models

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
	"strconv"
)

const EntryBucketName string = "graph"
const DatabaseName string = "db/HumanPredictions.db"

type GraphEntry struct {
	Target string
	Score  float64
}

func (e GraphEntry) makesCutoff(cutoff float64) bool {
	return e.Score >= cutoff
}

//Assume entries are large ~1-2MB
func retrieveEntries(db *bolt.DB, key string) *[]GraphEntry {
	var results []GraphEntry
	var queryResult []byte

	db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(EntryBucketName))
		queryResult = b.Get([]byte(key))
		return nil
	})

	if err := json.Unmarshal(queryResult, &results); err != nil {
		log.Println("Could not decode database entry for: " + key)
	}

	return &results
}

func genEntryMap(entries *[]GraphEntry, cutoff float64) map[string]float64 {
	m := make(map[string]float64)
	for _, entry := range *entries {
		if entry.makesCutoff(cutoff) {
			m[entry.Target] = entry.Score
		}
	}
	return m
}

//For visualizing graph
type Node struct {
	Name  string `json:"id"`
	Label string `json:"label"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Size  int    `json:"size"`
}

type Edge struct {
	Name        string  `json:"id"`
	Origin      string  `json:"source"`
	Destination string  `json:"target"`
	Score       float64 `json:"weight"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

func RetrieveSubgraph(db *bolt.DB, searchNodes []string, cutoff float64) Graph {
	edges := make([]Edge, 0)
	nodes := make([]Node, 0)
	nodeSet := make(map[string]bool)
	x := 0
	y := 0
	edgeID := 0
	for _, n := range searchNodes {
		if _, present := nodeSet[n]; !present {
			nodes = append(nodes, Node{n, n, x, y, 5})
			nodeSet[n] = true
		}
		refNodes := genEntryMap(retrieveEntries(db, n), cutoff)
		for destName, score := range refNodes {
			edgeID++
			x = (x + 1) % 3
			y = edgeID / 3

			if _, present := nodeSet[destName]; !present {
				nodes = append(nodes, Node{destName, destName, x, y, 3})
				nodeSet[destName] = true
			}
			edges = append(edges, Edge{strconv.Itoa(edgeID) + "e", n, destName, score})
		}

		for refName, _ := range refNodes {
			otherNodes := genEntryMap(retrieveEntries(db, refName), cutoff)

			for otherNodeName, score := range otherNodes {
				_, present := refNodes[otherNodeName]
				if present {
					edges = append(edges,
						Edge{strconv.Itoa(edgeID), refName, otherNodeName, score})
					edgeID++
				}
			}
		}
	}

	return Graph{nodes, edges}
}
