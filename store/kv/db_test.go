package kv

import (
	"context"
	"github.com/SirMetathyst/go-document"
	"github.com/SirMetathyst/go-kv/store/bolt"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/bbolt"
	"log"
	"os"
	"testing"
)

//func TestDB_StoreDocument_DoesNotErrorWithoutBucket(t *testing.T) {
//	db, closeFn := BoltDB()
//	defer closeFn()
//
//	err := db.StoreDocument(context.Background(), nil)
//	assert.Nil(t, err)
//
//	err = db.StoreDocument(context.Background(), document.Bucket{})
//	assert.Nil(t, err)
//
//	err = db.StoreDocument(context.Background(), document.Bucket{}, []document.Marshaler{}...)
//	assert.Nil(t, err)
//
//	err = db.StoreDocument(context.Background(), nil, []document.Marshaler{}...)
//	assert.Nil(t, err)
//}

//func TestDB_StoreDocument_DoesNotErrorWithoutKV(t *testing.T) {
//	db, closeFn := BoltDB()
//	defer closeFn()
//
//	err := db.StoreDocument(context.Background(), document.Bucket("default"))
//	assert.Nil(t, err)
//
//  err = db.StoreDocument(context.Background(), document.Bucket("default"), []document.Marshaler{}...)
//	assert.Nil(t, err)
//}

func TestDB_StoreDocument_CreatesKeysWhenTheyDontExist(t *testing.T) {
	db, closeFn := BoltDB()
	defer closeFn()

	bucket := document.Bucket("default")
	data := []document.Marshaler{
		&Document{Key: document.Key("key1"), Value: []byte("value1")},
		&Document{Key: document.Key("key2"), Value: []byte("value2")},
	}

	err := db.StoreDocument(context.Background(), bucket, data)
	//err := db.StoreDocument(context.Background(), bucket, &Document{Key: document.Key("key1"), Value: []byte("value1")})
	//err := db.StoreDocument(context.Background(), bucket, &Document{Key: document.Key("key1"), Value: []byte("value1")}, &Document{Key: document.Key("key2"), Value: []byte("value2")})
	assert.Nil(t, err)

	list, err := db.FetchDocument(context.Background(), bucket, extractKeys(data)...)
	//list, err := db.FetchDocument(context.Background(), bucket, document.Key("key1"))
	//list, err := db.FetchDocument(context.Background(), bucket, document.Key("key1"), document.Key("key2"))
	assert.Nil(t, err)

	assert.Contains(t, list, &Document{Key: document.Key("key1"), Value: []byte("value1")})
	assert.Contains(t, list, &Document{Key: document.Key("key2"), Value: []byte("value2")})
}

func BoltDB() (*DB[*Document], func()) {

	db, err := bolt.Open("temp.db", 0600, bbolt.DefaultOptions)
	if err != nil {
		log.Fatal(err)
	}

	closer := func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
		err = os.Remove("temp.db")
		if err != nil {
			panic(err)
		}
	}

	store, err := New[*Document](db)
	if err != nil {
		log.Fatal(err)
	}

	return store, closer
}

func extractKeys(v []document.Marshaler) (keys [][]byte) {
	for _, n := range v {
		k, _, err := n.MarshalDocument()
		if err != nil {
			log.Fatal(err)
		}
		keys = append(keys, k)
	}
	return
}
