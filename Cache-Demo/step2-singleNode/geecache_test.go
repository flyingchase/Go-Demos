package step2_singleNode

import (
	"fmt"
	"log"
	"testing"
)

var db = map[string]string{
	"A": "1",
	"B": "2",
	"C": "3",
}

func TestGet(t *testing.T) {
	loadCoundts := make(map[string]int, len(db))
	gee := NewGroup("sorces", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] Search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCoundts[key]; !ok {
					loadCoundts[key] = 0
				}
				loadCoundts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exit", key)
		}))
	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of A")
		}
		if _, err := gee.Get(k); err != nil || loadCoundts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}
	if view, err := gee.Get("unknown"); err != nil {
		t.Fatalf("the value of unkown should be empty, but %s got", view)
	}
}
