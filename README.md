[![Circle CI](https://circleci.com/gh/ahamidi/haigo.svg?style=svg)](https://circleci.com/gh/ahamidi/haigo)
[![GoDoc](https://godoc.org/github.com/ahamidi/haigo?status.svg)](https://godoc.org/github.com/ahamidi/haigo)

# Haigo 
YeSQL-like abstraction layer for Mongo on Go

### What is Haigo?
Haigo lets you write and use standard MongoDB queries directly from within your
Go application.

### Why?
While working on a project that leveraged Mongo's aggregation framework
extensively, we found that writing complicated queries in Go was painful and
brittle with deeply nested `bson.M` or `map[string]interface{}` clauses.

With Haigo, you can use the JSON formatted Mongo queries directly.

### Usage

1. Connect to Mongo (as you would normally)

1. Load Query File:
    ```hf, err := haigo.LoadQueryFile("queries.yml")```

1. Configure Haigo Params:
    ```params := haigo.Params{"name": "Ali", "age": 32}```

1. Execute Query/Pipe:
    ```q, err := hf.Queries["FindUser"].Pipe(col, params)```

1. Handle result (as you would normally):
    ```q.Count()```

### TODO
- [x] Parse MongoDB Query YAML
- [x] Detect Queries
- [x] Replace Query Params
- [x] Integrate CI
- [ ] More/Better Tests
- [ ] Better Error Handling


