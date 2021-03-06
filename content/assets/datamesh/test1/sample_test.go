package main

import (
	"bytes"
	"fmt"
	"log"

	"cuelang.org/go/cue/load"
)

func ExampleGenerateOpenAPI() {
	src := bytes.NewBufferString(`#Identity: {
		// first name of the person
		first: =~ "[A-Z].*"
		// Last name of the person
		Last: =~ "[A-Z].*"
		// Age of the person
		Age?: number & < 130
	}
	`)
	b, err := generateOpenAPI("-", &load.Config{
		Stdin: src,
	})
	if err != nil {
		log.Fatal(err)
	}
	// This contains the OpenAPI.json definition
	fmt.Println(string(b))
	// output:
	// {
	//    "openapi": "3.0.0",
	//    "info": {
	//       "title": "Generated by cue.",
	//       "version": "no version"
	//    },
	//    "paths": {},
	//    "components": {
	//       "schemas": {
	//          "Identity": {
	//             "type": "object",
	//             "required": [
	//                "first",
	//                "Last"
	//             ],
	//             "properties": {
	//                "first": {
	//                   "description": "first name of the person",
	//                   "type": "string",
	//                   "pattern": "[A-Z].*"
	//                },
	//                "Last": {
	//                   "description": "Last name of the person",
	//                   "type": "string",
	//                   "pattern": "[A-Z].*"
	//                },
	//                "Age": {
	//                   "description": "Age of the person",
	//                   "type": "number",
	//                   "maximum": 130,
	//                   "exclusiveMaximum": true
	//                }
	//             }
	//          }
	//       }
	//    }
	// }
}
