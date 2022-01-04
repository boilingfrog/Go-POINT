package main

import (
	"fmt"

	"github.com/golang/groupcache"
)

type SlowDB struct {
	data map[string]string
}

func (db *SlowDB) Get(key string) string {
	fmt.Printf("getting %s\n", key)
	return db.data[key]
}

func (db *SlowDB) Set(key string, value string) {
	fmt.Printf("setting %s to %s\n", key, value)
	db.data[key] = value
}

func NewSlowDB() *SlowDB {
	ndb := new(SlowDB)
	ndb.data = make(map[string]string)
	return ndb
}

func main() {

	db := NewSlowDB()

	db.Set("foo", "bar")
	db.Set("one", "two")

	var stringcache = groupcache.NewGroup("SlowDBCache", 64<<20, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			result := db.Get(key)
			dest.SetBytes([]byte(result))
			return nil
		}))

	var data []byte

	err := stringcache.Get(nil, "foo", groupcache.AllocatingByteSliceSink(&data))

	err2 := stringcache.Get(nil, "one", groupcache.AllocatingByteSliceSink(&data))

	db.Set("foo", "bar2")
	err3 := stringcache.Get(nil, "foo", groupcache.AllocatingByteSliceSink(&data))

	if err != nil {
		fmt.Println("error")
	}

	if err2 != nil {
		fmt.Println("error2")
	}

	if err3 != nil {
		fmt.Println("error3")
	}

	fmt.Printf("data was %s\n", data)

}
