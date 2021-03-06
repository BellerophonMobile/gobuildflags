GoFlagBuilder
============

GoFlagBuilder provides simple tools to construct command line flags
and a file parser to manipulate a given structure.  It uses reflection
to traverse a potentially hierarchical object of structs and maps and
install handlers in Go's standard flag package.  Constructed parsers
scan given documents line by line and apply key/value pairs found.

Constructed flags have the form Foo=..., Bar=..., Obj.Field=... and so
on.  Maps and exposed fields of structs with primitive types are
consumed, so in this case Foo and Bar might be map keys or public
struct fields to a primitive.  Nested maps and structs are followed,
producing dot-notation hierarchical keys such as Obj.Field.

Primitive types understood by GoFlagBuilder include bool, float64,
int64, int, string, uint64, and uint. Slices of these primitive types
are supported as well.

Primitive fields in the given object and sub-objects must be settable.
In general this means structs should be passed in as pointers.  Maps
may also be set directly.

[![Build Status](https://travis-ci.org/BellerophonMobile/goflagbuilder.svg?branch=master)](https://travis-ci.org/BellerophonMobile/goflagbuilder)
[![GoDoc](https://godoc.org/github.com/BellerophonMobile/goflagbuilder?status.svg)](https://godoc.org/github.com/BellerophonMobile/goflagbuilder)
[![Go Report Card](https://goreportcard.com/badge/github.com/BellerophonMobile/goflagbuilder)](https://goreportcard.com/report/github.com/BellerophonMobile/goflagbuilder)

## Example

A very simple example:

```go
package main

import (
	"flag"
	"log"

	"github.com/BellerophonMobile/goflagbuilder/v2"
)

type server struct {
	Domain string
	Port   int
}

func Example_Simple() {

	myserver := &server{}

	// Construct the flags
	if err := From(myserver); err != nil {
		log.Fatal("Error: " + err.Error())
	}

	// Read from the command line to establish the param
	flag.Parse()

}
```

This would establish the command line flags "-Port" and "-Domain".

A more elaborate example including nested structures and using the parser is
available
[here](https://github.com/BellerophonMobile/goflagbuilder/blob/master/doc_extended_test.go).
There are also a series of tests in the package outlining exactly what input
structures are valid.


## Changelog

 * **2018/10/15: v2.1.0** Clean up path from FlagSet name for environment
                          variable parsing.
 * **2018/10/03: v2.0.3** Ignore unexported fields of structs.
 * **2018/10/03: v2.0.2** Fix crash when printing values.
 * **2018/10/03: v2.0.1** Remove mistakenly included logging, and formatting.
 * **2018/08/16: v2.0.0** Simplified package a bit, split config file
   parsing to its own sub-package. Added a new sub-package to read
   environment variables.
 * **2014/11/03: v1.0.0** Though not mature at all, we consider
   GoFlagBuilder to be usable.

## License

GoFlagBuilder is provided under the open source
[MIT license](http://opensource.org/licenses/MIT):

> The MIT License (MIT)
>
> Copyright (c) 2018 [Bellerophon Mobile](http://bellerophonmobile.com/)
>
>
> Permission is hereby granted, free of charge, to any person
> obtaining a copy of this software and associated documentation files
> (the "Software"), to deal in the Software without restriction,
> including without limitation the rights to use, copy, modify, merge,
> publish, distribute, sublicense, and/or sell copies of the Software,
> and to permit persons to whom the Software is furnished to do so,
> subject to the following conditions:
>
> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.
>
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
> NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
> BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
> ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
> CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
> SOFTWARE.
