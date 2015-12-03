# haigo
YeSQL for Mongo on Go

### Overview
Haigo lets you write and use standard MongoDB queries, easily.


### TODO

- [x] Parse MongoDB Query YAML
- [x] Detect Queries
- [x] Replace Query Params
- [ ] Tests
- [ ] Docs

### Design

#### MongoDB query file parsing

1. Read in full file.
1. Find all instances of `--`.
1. Parse the `name`, use as key for Query map.
1. Parse the query and track params (look for `:` followed by param name).

#### Usage

1. Call named query and pass `params` map.
