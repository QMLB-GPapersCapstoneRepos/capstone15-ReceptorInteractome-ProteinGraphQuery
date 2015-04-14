package controllers

import (
	"archive/tar"
	"compress/gzip"
	"github.com/boltdb/bolt"
	//"github.com/revel/revel"
	"bytes"
	"io"
	"io/ioutil"
	"json"
	"log"
	"os"
	"strings"
)

var graphBucketName []byte

type GraphEntry struct {
  Target byte[]
  Score	byte[]
}

func (entry *GraphEntry) Encode() error {
  result, err = json.Marshal(entry)
  Check(err)
  return result
}

func InitDB() {
	dbName := "db/HumanPredictions.db"
	graphBucketName := []byte("graph")

	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		log.Println("Generating new database...")
		db, err := bolt.Open(dbName, 0600, nil)
		Check(err)
		defer db.Close()
		db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucket(graphBucketName)
			Check(err)
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
				ImportFile(file, db)
			}
		}
	}
	log.Println("Database is setup")
}

func ImportArchivedFile(f *os.File, db *bolt.DB) {
	var contents []byte
	var lines [][]byte
	var entries []byte

	reader, err = gzip.NewReader(f)
	Check(err)
	reader := tar.NewReader(reader)
	header, err := reader.Next()
	Check(err)

	if header.Typeflag == tar.TypeReg {
		if contents, err = ioutil.ReadAll(reader); err != nil {
			Check(err)
		}
		lines := bytes.Split(contents, '\n')
	} else {
		log.Printf("Unexpected file type in archive: $s!", header.Name)
		return
	}

	db.Update(func(tx *bolt.Tx) error {
		graphBucket := tx.Bucket(graphBucketName)
		for _, line := range lines {
			items := bytes.Split(line, '\t')
			dump := GraphEntry{,""}
			err = graphBucket.put(items[1], dump)
			Check(err)
		}
	})
}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
