package main

import (
	"Cache-Demo/step3/geecache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"A": "1",
	"B": "2",
	"C": "3",
}

func main() {
	geecache.NewGroup("sores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {

			log.Println("[SolwDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	addr := "localhost:9908"
	peers := geecache.NewHTTPPool(addr)
	log.Println("geecache is running at ", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
