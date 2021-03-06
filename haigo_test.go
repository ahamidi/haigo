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

	// Setup
	err := seedDB()
	if err != nil {
		log.Println("Seed DB error:", err)
	}

	retCode := m.Run()

	// Teardown
	err = dropDB()
	if err != nil {
		log.Println("Unable to drop collection:", err)
	}

	os.Exit(retCode)

}

func TestParseMongoFile(t *testing.T) {
	_, err := parseMongoFile("examples/queries.yml")

	assert.NoError(t, err)
}

func TestGenerateQuery(t *testing.T) {
	hf, _ := parseMongoFile("examples/queries.yml")

	q, err := hf.Queries["basic-select"].Map(map[string]interface{}{
		"type": "hi",
	})

	assert.NoError(t, err)
	assert.Equal(t, "hi", q.(map[string]interface{})["type"])
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

	hf, _ := parseMongoFile("examples/queries.yml")
	q, err := hf.Queries["basic-select"].Map(map[string]interface{}{
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

	hf, _ := parseMongoFile("examples/queries.yml")
	q, err := hf.Queries["conditional"].Map(map[string]interface{}{
		"qty":  100,
		"name": "apple",
	})

	res := db.C("testcol").Find(q)

	cnt, err := res.Count()

	assert.NoError(t, err)
	assert.Equal(t, 2, cnt)

}

func TestExecute(t *testing.T) {
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
	col := sess.DB(dbname).C("testcol")

	hf, _ := LoadQueryFile("examples/queries.yml")
	params := Params{
		"type": "Good",
	}
	res, err := hf.Queries["basic-select"].Query(col, params)

	cnt, err := res.Count()

	assert.NoError(t, err)
	assert.Equal(t, 1, cnt)

}

func TestPrint(t *testing.T) {
	hf, _ := parseMongoFile("examples/queries.yml")
	err := hf.Queries["conditional"].Print(map[string]interface{}{
		"qty":  100,
		"name": "apple",
	})

	assert.NoError(t, err)
}

func TestString(t *testing.T) {
	hf, _ := parseMongoFile("examples/queries.yml")
	q, err := hf.Queries["conditional"].String(map[string]interface{}{
		"qty":  100,
		"name": "apple",
	})

	assert.NoError(t, err)
	assert.IsType(t, "hello", q)
}

func TestMissingMongoFile(t *testing.T) {
	_, err := parseMongoFile("examples/querie.yml")

	assert.Error(t, err)
}

func TestMissingQuery(t *testing.T) {
	hf, _ := parseMongoFile("examples/queries.yml")
	_, err := hf.Queries["missing"].String(map[string]interface{}{
		"qty":  100,
		"name": "apple",
	})

	assert.Error(t, err)
}

// Here is an example of Haigo's basic usage.
//
// In this case we're loading the "basic-select" query from the
// "examples/queries.yml" MongoDB query file and then executing it against the
// "mycol" collection on "mydb".
//
// "Query" simply returns a mgo.Query struct and as such supports all of the
// expected methods (in this case we're calling Count().
func Example() {
	// Dial MongoDB Server
	sess, _ := mgo.Dial("127.0.0.1")
	col := sess.DB("mydb").C("mycol")

	// Load MongoDB Query File
	hf, _ := LoadQueryFile("examples/queries.yml")
	params := Params{
		"type": "Good",
	}

	// Call "basic-select" query.
	res, _ := hf.Queries["basic-select"].Query(col, params)

	// Returns mgo.Query struct
	cnt, _ := res.Count()

	fmt.Println(cnt)
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
	return db.C("testcol").DropCollection()
}
