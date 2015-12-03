package haigo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMongoFile(t *testing.T) {
	_, err := parseMongoFile("queries.yml")

	assert.NoError(t, err)
}

func TestExecQuery(t *testing.T) {
	hf, _ := parseMongoFile("queries.yml")

	q, err := hf.Queries["basic-select"].Query(map[string]interface{}{"type": "hi"})

	t.Log("Q:", q)

	assert.NoError(t, err)
}
