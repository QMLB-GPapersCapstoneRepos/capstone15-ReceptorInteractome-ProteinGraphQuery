package controllers

import (
	"archive/tar"
	"compress/gzip"
	"github.com/boltdb/bolt"
	//"github.com/revel/revel"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const graphBucketName string = "graph"

type GraphEntry struct {
	Target string
	Score  float64
}

func (entry *GraphEntry) Encode() []byte {
	result, err := json.Marshal(entry)
	Check(err)
	return result
}

func InitDB() {
	dbName := "db/HumanPredictions.db"

	//if _, err := os.Stat(dbName); os.IsNotExist(err) {
	log.Println("Generating new database...")
	db, err := bolt.Open(dbName, 0600, nil)
	Check(err)
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(graphBucketName))
		//Check(err)
		log.Println(err)
		return nil
	})

	entries, err := ioutil.ReadDir("public/predictions")
	Check(err)

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".tar.gz") {
			log.Println(name)
			file, err := os.Open("public/predictions/" + name)
			Check(err)
			defer file.Close()
			ImportArchivedFile(file, db)
		}
	}
	//}
	log.Println("Database is setup")
}

func ImportArchivedFile(f *os.File, db *bolt.DB) {
	var contents []byte

	gzipReader, err := gzip.NewReader(f)
	Check(err)
	reader := tar.NewReader(gzipReader)
	header, err := reader.Next() //error: truncated gzip input for human.hprdlabel.receptor.RFall.ScorLablFeaGeneInfo.allhuman.ggi.10886.addDisease.tar.gz
	Check(err)

	if header.Typeflag != tar.TypeReg {
		log.Printf("Unexpected file type in archive: $s!", header.Name)
		return
	}

	if contents, err = ioutil.ReadAll(reader); err != nil {
		log.Println("Had a problem reading this archive!")
		log.Println(err)
		contents = contents[:len(contents)-1] //assume last item is incomplete and remove
	}
	lines := bytes.Split(contents, []byte("\n"))

	db.Update(func(tx *bolt.Tx) error {
		graphBucket := tx.Bucket([]byte(graphBucketName))
		for _, line := range lines[1:] {
			if len(line) == 0 {
				continue
			}
			items := bytes.Split(line, []byte("\t"))
			score, err := strconv.ParseFloat(string(items[4]), 64)
			Check(err)
			entry := GraphEntry{string(items[3]), score}
			Check(err)
			if key := items[1]; len(key) > 0 {
				err = graphBucket.Put(key, entry.Encode())
				Check(err)
			}
		}
		return nil
	})
}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
