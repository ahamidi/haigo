package haigo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"text/template"

	"gopkg.in/mgo.v2"
	"gopkg.in/yaml.v2"
)

type param struct {
	Name string
	Type string
}

type Query struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description,omitempty"`
	QueryString string  `yaml:"query"`
	params      []param // TODO
}

type Params map[string]interface{}

// Returns configured mgo Pipe.
func (h *Query) Pipe(col *mgo.Collection, params Params) (*mgo.Pipe, error) {

	q, err := h.Map(params)
	if err != nil {
		return nil, err
	}

	return col.Pipe(q), nil
}

// Returns configured mgo Query.
func (h *Query) Query(col *mgo.Collection, params Params) (*mgo.Query, error) {

	q, err := h.Map(params)
	if err != nil {
		return nil, err
	}

	return col.Find(q), nil
}

// Accepts a params map and returns a map for use with the mgo `find()`
// function.
func (h *Query) Map(params Params) (interface{}, error) {

	// Create the template
	t, err := template.New("haigo").Parse(h.QueryString)
	if err != nil {
		return nil, err
	}

	// Buffer to capture string
	buf := new(bytes.Buffer)

	// Execute template
	err = t.Execute(buf, params)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON into Map
	var m interface{}
	err = json.Unmarshal(buf.Bytes(), &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Pretty prints the configured query.
func (h *Query) Print(params Params) error {
	m, err := h.Map(params)
	if err != nil {
		return err
	}

	jsonM, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("Name: %s\nDescription: %s\nQuery:\n%s\n", h.Name, h.Description, string(jsonM))
	return nil
}

// Return configured query as a string.
func (h *Query) String(params Params) (string, error) {
	if h == nil {
		return "", errors.New("Query Not Found.")
	}
	return h.QueryString, nil
}

// YAML formatted file with MongoDB Queries.
//
//  ---
//    - name: basic-select
//      description: Basic MongoDB Select
//      query: '{"type": {{.type}} }'
//
//    - name: conditional
//      description: Conditional Query
//      query: '{
//         "type": "food",
//         "$or": [ { "qty": { "$gt": {{.qty}} } }, { "name": {{.name}} } ]
//      }'
type File struct {
	Queries map[string]*Query
}

func (m *File) unmarshalYAML(data []byte) error {

	var hqs []Query

	err := yaml.Unmarshal(data, &hqs)
	if err != nil {
		return err
	}

	qm := make(map[string]*Query)

	for i := range hqs {
		qm[hqs[i].Name] = &hqs[i]
	}

	m.Queries = qm

	return nil
}

// Reads in Mongo Query File for use with Haigo.
func LoadQueryFile(file string) (*File, error) {
	return parseMongoFile(file)
}
