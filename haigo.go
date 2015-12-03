package haigo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"
)

type Param struct {
	Name string
	Type string
}

type HaigoQuery struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	QueryString string `yaml:"query"`
	Params      []Param
}

type HaigoFile struct {
	Queries map[string]*HaigoQuery
}

type QueryString string

func (qi *QueryString) bytes() []byte {
	return []byte(string(*qi))
}

func (qi QueryString) GetBSON() (interface{}, error) {
	br := bson.Raw{3, qi.bytes()}
	return br, nil
}

func (m *HaigoFile) unmarshalYAML(data []byte) error {

	var hqs []HaigoQuery

	err := yaml.Unmarshal(data, &hqs)
	if err != nil {
		return err
	}

	qm := make(map[string]*HaigoQuery)

	for _, q := range hqs {
		qm[q.Name] = &q
	}

	m.Queries = qm

	return nil
}

// sanitizeParams - Adds single quotes if param is a string (as needed by Mongo).
func sanitizeParams(params map[string]interface{}) map[string]interface{} {
	for k, v := range params {
		switch v.(type) {
		case string:
			params[k] = fmt.Sprintf("\"%s\"", v)
		}

	}
	return params
}

// Query - Accepts a params map and returns a bson.M.
func (h *HaigoQuery) Query(params map[string]interface{}) (map[string]interface{}, error) {

	// Create the template
	t, err := template.New("what").Parse(h.QueryString)
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

	//qs := QueryString(buf.String())

	// Unmarshal JSON into Map
	var m map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
