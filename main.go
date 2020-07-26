package main

import (
	"encoding/hex"
	"flag"
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

func put(k string, b []byte) error {
	c.Put(k, b)
	return ioutil.WriteFile(k, b, 0644)
}

func handler(w http.ResponseWriter, r *http.Request) {
	k := dataDir + "/" + safeFilename(r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		b, err := get(k)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(b)
	// Makes using curl simpler (no need to specify -X)
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := put(k, b); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
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
