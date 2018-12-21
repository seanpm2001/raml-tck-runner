package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

// FileResult represents a single file parsing result
type FileResult struct {
	File    string `json:"file"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// Report represents a parser parsing results
type Report struct {
	Parser  string        `json:"parser"`
	Results []*FileResult `json:"results"`
}

func main() {
	parserFl := flag.String(
		"parser", "jumpscale",
		"Parser to test. Supported: jumpscale, go-raml, tsaikd.")
	flag.Parse()

	parsers := map[string]Parser{
		"jumpscale": Jumpscale,
		"go-raml":   Goraml,
		"tsaikd":    Tsaikd,
	}
	parser, ok := parsers[*parserFl]
	if !ok {
		fmt.Println("Not supported parser. See help (-h).")
		return
	}

	examplesFl := CloneTckRepo()
	fileList, err := ListRamls(examplesFl)
	if err != nil {
		fmt.Printf("Failed to list RAML files: %s\n", err)
		return
	}

	report := &Report{
		Parser:  *parserFl,
		Results: []*FileResult{},
	}

	for _, fpath := range fileList {
		err, notPanic := parser(fpath)
		if !notPanic {
			err = errors.New("Parser crashed")
		}

		result := &FileResult{
			File:    strings.Replace(fpath, examplesFl, "", -1),
			Success: err == nil,
			Error:   "",
		}
		if err != nil {
			result.Error = err.Error()
		}
		report.Results = append(report.Results, result)
	}

	SaveReport(report)
}
