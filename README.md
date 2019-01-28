go-isbn
===
[![Build Status](https://travis-ci.org/nasa9084/go-isbn.svg?branch=master)](https://travis-ci.org/nasa9084/go-isbn)
[![GoDoc](https://godoc.org/github.com/nasa9084/go-isbn?status.svg)](https://godoc.org/github.com/nasa9084/go-isbn)

a package for parsing/validating ISBN(International Standard Book Number).

## SYNOPSIS

``` go
import (
    "log"

    "github.com/nasa9084/go-isbn"
)
var n = "ISBN978-4-00-310101-8"

func main() {
    // parse
	code, err := isbn.Parse(n)
	if err != nil {
		log.Fatal(err)
	}
    // validate
	if !code.IsValid() {
		log.Print("given ISBN code is invalid")
	}
    // update legacy 10-length code to 13-length code
	if code.IsLegacy {
		code, err = code.Update()
		if err != nil {
			log.Fatal(err)
		}
	}
    // format
	fmt.Printf("%s", code.String())
}
```
