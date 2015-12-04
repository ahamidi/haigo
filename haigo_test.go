package haigo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
)

func TestMain(m *testing.M) {
	err := seedDB()
	if err != nil {
		log.Println("Seed DB error:", err)
	}

	os.Exit(m.Run())

	dropDB()
}

func TestParseMongoFile(t *testing.T) {
	_, err := parseMongoFile("queries.yml")

	assert.NoError(t, err)
}

func TestGenerateQuery(t *testing.T) {
	hf, _ := parseMongoFile("queries.yml")

	q, err := hf.Queries["basic-select"].Query(map[string]interface{}{
		"type": "hi",
	})

	assert.NoError(t, err)
	assert.Equal(t, "hi", q["type"])
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
		dbname = "test"
	}
	db := sess.DB(dbname)

	hf, _ := parseMongoFile("queries.yml")
	q, err := hf.Queries["basic-select"].Query(map[string]interface{}{
		"type": "Good",
	})

	res := db.C("testcol").Find(q)

	cnt, err := res.Count()

	assert.NoError(t, err)
	assert.Equal(t, 1, cnt)
}

func TestConditionalQuery(t *testing.T) {
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
		dbname = "test"
	}
	db := sess.DB(dbname)

	hf, _ := parseMongoFile("queries.yml")
	q, err := hf.Queries["conditional"].Query(map[string]interface{}{
		"qty":  100,
		"name": "apple",
	})

	res := db.C("testcol").Find(q)

	cnt, err := res.Count()

	assert.NoError(t, err)
	assert.Equal(t, 2, cnt)

}

func readSeedFile(file string) ([]interface{}, error) {

	ms := []interface{}{}

	seedFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer seedFile.Close()

	scanner := bufio.NewScanner(seedFile)
	for scanner.Scan() {
		var m map[string]interface{}
		err := json.Unmarshal(scanner.Bytes(), &m)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	return ms, nil
}

func seedDB() error {
	uri := os.Getenv("MONGO_URL")
	if uri == "" {
		uri = "127.0.0.1"
	}

	sess, err := mgo.Dial(uri)
	if err != nil {
		log.Printf("Can't connect to mongo, go error %v\n", err)
		return err
	}

	dbname := os.Getenv("MONGO_DB")
	if dbname == "" {
		dbname = "test"
	}
	db := sess.DB(dbname)

	seedData, _ := readSeedFile("seed.json")

	b := db.C("testcol").Bulk()
	b.Insert(seedData...)

	_, err = b.Run()
	if err != nil {
		return err
	}

	return nil
}

func dropDB() error {
	uri := os.Getenv("MONGO_URL")
	if uri == "" {
		uri = "127.0.0.1"
	}

	sess, err := mgo.Dial(uri)
	if err != nil {
		log.Printf("Can't connect to mongo, go error %v\n", err)
		return err
	}

	dbname := os.Getenv("MONGO_DB")
	if dbname == "" {
		dbname = "test"
	}
	db := sess.DB(dbname)
	return db.DropDatabase()
}
