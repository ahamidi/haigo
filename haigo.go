package haigo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"gopkg.in/mgo.v2"
	"gopkg.in/yaml.v2"
)

type param struct {
	Name string
	Type string
}

type HaigoQuery struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description,omitempty"`
	QueryString string  `yaml:"query"`
	params      []param // TODO
}

type HaigoParams map[string]interface{}

// Execute - Returns configured mgo Query.
func (h *HaigoQuery) Execute(col *mgo.Collection, params HaigoParams) (*mgo.Query, error) {

	q, err := h.Query(params)
	if err != nil {
		return nil, err
	}

	return col.Find(q), nil
}

// HaigoFile - YAML formatted file with MongoDB Queries.
//
//	---
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
type HaigoFile struct {
	Queries map[string]*HaigoQuery
}

func (m *HaigoFile) unmarshalYAML(data []byte) error {

	var hqs []HaigoQuery

	err := yaml.Unmarshal(data, &hqs)
	if err != nil {
		return err
	}

	qm := make(map[string]*HaigoQuery)

	for i := range hqs {
		qm[hqs[i].Name] = &hqs[i]
	}

	m.Queries = qm

	return nil
}

// sanitizeParams - Adds single quotes if param is a string (as needed by Mongo).
func sanitizeParams(params HaigoParams) HaigoParams {
	for k, v := range params {
		switch v.(type) {
		case string:
			params[k] = fmt.Sprintf("\"%s\"", v)
		}

	}
	return params
}

// Query - Accepts a params map and returns a map for use with the mgo `find()`
// function.
func (h *HaigoQuery) Query(params HaigoParams) (map[string]interface{}, error) {

	// Create the template
	t, err := template.New("haigo").Parse(h.QueryString)
	if err != nil {
		return nil, err
	}

	// Buffer to capture string
	buf := new(bytes.Buffer)

	// Execute template
	err = t.Execute(buf, sanitizeParams(params))
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON into Map
	var m map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// LoadQueryFile - Reads in Mongo Query File for use with Haigo.
func LoadQueryFile(file string) (*HaigoFile, error) {
	return parseMongoFile(file)
}
