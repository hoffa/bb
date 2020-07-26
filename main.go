package main

import (
	"encoding/hex"
	"flag"
	"io"
	"io/ioutil"
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

func put(k string, r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	c.Put(k, b)
	return ioutil.WriteFile(k, b, 0644)
}

func handler(w http.ResponseWriter, r *http.Request) {
	k := dataDir + "/" + safeFilename(r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		if b, err := get(k); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(b)
		}
	// Makes using curl simpler (no need to specify -X)
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		if err := put(k, r.Body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func main() {
	addr := flag.String("addr", ":8080", "address to listen to")
	flag.Parse()

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler)
	panic(http.ListenAndServe(*addr, nil))
}
