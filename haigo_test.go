package haigo

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
)

func TestParseMongoFile(t *testing.T) {
	_, err := parseMongoFile("queries.yml")

	assert.NoError(t, err)
}

func TestGenerateQuery(t *testing.T) {
	hf, _ := parseMongoFile("queries.yml")

	q, err := hf.Queries["basic-select"].Query(map[string]interface{}{"type": "hi"})

	assert.NoError(t, err)
	assert.Equal(t, q["type"], "hi")

}

func TestExecQuery(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	uri := os.Getenv("MONGO_URL")
	if uri == "" {
		uri = "127.0.0.1"
	}

	sess, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}

	// sess.SetSafe(&mgo.Safe{})
	dbname := os.Getenv("MONGO_DB")
	if dbname == "" {
		dbname = "haigo"
	}
	db := sess.DB(dbname)

	hf, _ := parseMongoFile("queries.yml")
	q, err := hf.Queries["basic-select"].Query(map[string]interface{}{"type": "Good"})

	res := db.C("things").Find(q)

	cnt, err := res.Count()

	assert.NoError(t, err)
	assert.Equal(t, 1, cnt)

}
