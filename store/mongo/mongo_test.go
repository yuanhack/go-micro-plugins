package mongo

import (
	"encoding/json"
	"fmt"
	"github.com/macheal/go-micro/v2/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestModel struct {
	Key string
	A   int
}

func TestMongo(t *testing.T) {
	url := "mongodb://root:example@172.16.0.210:30007/setting?authSource=admin"
	db := NewStore(URI(url))
	bytes, err := json.Marshal(TestModel{
		Key: "k2",
		A:   1,
	})
	//if err != nil {
	//	return err
	//}
	testKey := "test"
	data := store.Record{
		Key:   testKey,
		Value: bytes,
		//Expiry: 10 * time.Second,
	}
	err = db.Write(&data)
	assert.NoError(t, err)
	res, err := db.Read(testKey)
	fmt.Printf("len:%d,res[0]: %+v\r\n", len(res), res[0])
	res, err = db.Read(testKey, store.ReadPrefix())
	assert.NoError(t, err)
	for i, re := range res {
		fmt.Println(i, re)
	}

	list, err := db.List(store.ListPrefix("test"))
	fmt.Println("listPrefix test :", err, list)

	list, err = db.List(store.ListSuffix("test"))
	fmt.Println("listSuffix test:", err, list)
}

func TestMongoList(t *testing.T) {
	url := "mongodb://root:example@172.16.0.210:30007/setting?authSource=admin"
	db := NewStore(URI(url), store.Database("micro"), store.Table("store"))

	db.Write(&store.Record{
		Key:   "aabbccdd",
		Value: []byte("this is a test"),
	})

	ss, err := db.List()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("---all---")
	for _, v := range ss {
		fmt.Println(v)
	}

	ss, err = db.List(store.ListPrefix("aa"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("---prefix aa---")
	for _, v := range ss {
		fmt.Println(v)
	}

	ss, err = db.List(store.ListPrefix("aa"), store.ListSuffix("dd"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("---prefix aa, suffix dd---")
	for _, v := range ss {
		fmt.Println(v)
	}
}
