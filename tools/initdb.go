package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/nickjanus/ProteinGraphQuery/app/models"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	log.Println("Generating new database...")
	db, err := bolt.Open(models.DatabaseName, 0600, nil)
	Check(err)
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(models.EntryBucketName))
		Check(err)
		return nil
	})

	entries, err := ioutil.ReadDir("public/predictions")
	Check(err)

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".tar.gz") {
			file, err := os.Open("public/predictions/" + name)
			Check(err)
			defer file.Close()
			wg.Add(1)
			go func() {
				defer wg.Done()
				ImportArchivedFile(file, db)
			}()
		}
	}
	wg.Wait()
	log.Println("Database is setup")
}

func ImportArchivedFile(f *os.File, db *bolt.DB) {
	var contents []byte
	var target []byte
	entries := make([]models.GraphEntry, 0)

	gzipReader, err := gzip.NewReader(f)
	Check(err)
	reader := tar.NewReader(gzipReader)
	header, err := reader.Next()
	Check(err)

	if header.Typeflag != tar.TypeReg {
		log.Printf("Unexpected file type in archive: $s!", header.Name)
		return
	}

	if contents, err = ioutil.ReadAll(reader); err != nil {
		log.Println("Had a problem reading this archive: " + f.Name())
		log.Println(err)
		contents = contents[:len(contents)-1] //assume last item is incomplete and remove
	}

	lines := bytes.Split(contents, []byte("\n"))
	key := bytes.Split(lines[1], []byte("\t"))[3]
	for _, line := range lines[1:] {
		if len(line) == 0 {
			continue
		}

		items := bytes.Split(line, []byte("\t"))
		if target = items[1]; len(target) > 0 {
			score, err := strconv.ParseFloat(string(items[4]), 64)
			Check(err)
			entry := models.GraphEntry{string(target), score}
			Check(err)
			entries = append(entries, entry)
		}
	}
	value, err := json.Marshal(entries)
	Check(err)

	db.Batch(func(tx *bolt.Tx) error {
		graphBucket := tx.Bucket([]byte(models.EntryBucketName))

		err = graphBucket.Put(key, value)
		Check(err)
		return nil
	})
}

func Check(err error) {
	if err != nil {
		debug.PrintStack()
		log.Fatal(err)
		panic(err)
	}
}
