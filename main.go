package main

import (
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/hoffa/bb/cache"
)

const dataDir = "data"

var c = cache.New(128)

func safeFilename(s string) string {
	return hex.EncodeToString([]byte(s))
}

func get(k string) ([]byte, error) {
	if b := c.Get(k); b != nil {
		return b, nil
	}
	b, err := ioutil.ReadFile(k)
	if err != nil {
		return nil, err
	}
	c.Put(k, b)
	return b, nil
}

func put(k string, b []byte) error {
	return ioutil.WriteFile(k, b, 0644)
}

func handler(w http.ResponseWriter, r *http.Request) {
	k := dataDir + "/" + safeFilename(r.URL.Path)
	switch r.Method {
	case "GET":
		b, err := get(k)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(b)
	case "PUT":
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := put(k, b); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
