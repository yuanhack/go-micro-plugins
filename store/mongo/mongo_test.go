package mongo

import (
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro/v2/store"
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
	res, err := db.Read(testKey, store.ReadPrefix())
	assert.NoError(t, err)
	fmt.Printf("res[0]: %+v\r\n", res[0])
	for i, re := range res {
		fmt.Println(i, re)
	}

	list, err := db.List(store.ListPrefix("test"))
	fmt.Println("listPrefix test :", err, list)

	list, err = db.List(store.ListSuffix("test"))
	fmt.Println("listSuffix test:", err, list)
}
