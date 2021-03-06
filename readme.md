# Simple DSN url parser

[![Build Status](https://travis-ci.com/sidmal/dsn-parser.svg?branch=main)](https://travis-ci.com/sidmal/dsn-parser)
[![codecov](https://codecov.io/gh/sidmal/dsn-parser/branch/master/graph/badge.svg)](https://codecov.io/gh/sidmal/dsn-parser)

Simple Go library to parse DSN urls.

## Installation

`go get github.com/sidmal/dsn-parser`

## Usage

```go
package main

import (
	dsnParser "github.com/sidmal/dsn-parser"
	"log"
)

func main() {
    parsedDsn, err := dsnParser.New("postgres://user:password@db1:5432/test?sslmode=disable")
	
    if err != nil {
        log.Fatalln(err)
    }
    
    log.Printf("DSN url parsed successfully. Parsed DSN url: %s", parsedDsn.Dsn)
}
```
