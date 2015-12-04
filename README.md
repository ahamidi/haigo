[![Circle CI](https://circleci.com/gh/ahamidi/haigo.svg?style=svg)](https://circleci.com/gh/ahamidi/haigo)
[![GoDoc](https://godoc.org/github.com/ahamidi/haigo?status.svg)](https://godoc.org/github.com/ahamidi/haigo)

# Haigo 
YeSQL-like abstraction layer for Mongo on Go

### Overview
Haigo lets you write and use standard MongoDB queries directly from within your
Go application.


### TODO

- [x] Parse MongoDB Query YAML
- [x] Detect Queries
- [x] Replace Query Params
- [x] Integrate CI
- [ ] Support Bindata
- [ ] Tests
- [ ] Docs

### Design

#### Haigo File

* YAML formatted file.
* MongoDB query as string.

#### Usage

1. Call named query and pass `params` map.
