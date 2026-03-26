# Lem-in

##  Description

Lem-in is a Go program that simulates an ant colony finding the fastest path from a start room to an end room through a network of tunnels.

The program:

* Parses a colony description from a file
* Builds a graph of rooms and links
* Finds optimal paths
* Simulates ant movement turn by turn

---

## Usage

```bash
go run ./cmd <filename>
```

Example:

```bash
go run ./cmd testdata/valid/test0.txt
```

---

## Features

* Robust input parsing and validation
* Graph representation using adjacency lists
* Shortest path search (BFS)
* Optimized ant distribution across paths
* Step-by-step movement simulation

---

##  Error Handling

Invalid input will return:

```
ERROR: invalid data format
```

---

##  Tech Stack

* Go (standard library only)

---

## Project Structure

* `cmd/` → entry point
* `internal/parser/` → input parsing
* `internal/graph/` → graph + algorithms
* `internal/simulation/` → ant movement
* `testdata/` → input test files

---

## Learning Goals

* Graph algorithms
* Data parsing
* Algorithm optimization
* Clean architecture in Go

### Authors
* Geremy Mwaro
* Amon Ochuka
