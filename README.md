# orderedmap
ordered map in golang generics way, support json.Unmarsal, json.Marshal

# example
```
package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/doublepu/orderedmap"
)

type Data struct {
	A int
	B int
}

func main() {
	o := orderedmap.NewOrderedMap[string, string]()
	//o.Set("a", "b")
	o.Set("c", "d")
	o.Set("a", "b")
	vv := map[string]string{
		"c": "d",
		"a": "b",
	}
	bb, _ := json.Marshal(vv)
	log.Println(string(bb))
	//o.Delete("c")
	for _, v := range o.List() {
		log.Println(v.K, v.V)
	}
	start := time.Now()
	b, err := json.Marshal(o)
	log.Println(string(b), err, time.Since(start))
	newO := orderedmap.NewOrderedMap[string, Data]()
	b = []byte(`{"c": {"A":2333}}`)
	log.Println(json.Unmarshal(b, newO))
	b, err = json.Marshal(newO)
	log.Println(string(b), err, time.Since(start))
}
```
