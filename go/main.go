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

// ParserMeta represents a single parser meta-data
type ParserMeta struct {
	Language string `json:"language"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	Version  string `json:"version"`
}

// Report represents a parser parsing results
type Report struct {
	Parser  ParserMeta    `json:"parser"`
	Results []*FileResult `json:"results"`
	Branch  string        `json:"branch"`
}

func main() {
	parserFl := flag.String(
		"parser", "jumpscale",
		"Parser to test. Supported: jumpscale, go-raml.")
	outdirFl := flag.String(
		"outdir", "./", "Output report directory path.")
	branchFl := flag.String(
		"branch", "", "raml-tck directory to load RAML files from.")
	flag.Parse()

	parsersRunners := map[string]Parser{
		"jumpscale": Jumpscale,
		"go-raml":   Goraml,
		"tsaikd":    Tsaikd,
	}

	/**
	 * Parsers meta-data which helps generating pretty reports.
	 * Required fields are: url, version.
	 */
	parsersMeta := map[string]ParserMeta{
		"jumpscale": {
			Language: "go",
			Name:     "jumpscale",
			Url:      "https://github.com/Jumpscale/go-raml/tree/master/raml",
			Version:  "0.1",
		},
		"go-raml": {
			Language: "go",
			Name:     "go-raml",
			Url:      "https://github.com/go-raml/raml",
			Version:  "not versioned",
		},
	}
	parser, ok := parsersRunners[*parserFl]
	if !ok {
		fmt.Println("Not supported parser. See help (-h).")
		return
	}

	examplesFl := CloneTckRepo(*branchFl)
	fileList, err := ListRamls(examplesFl)
	if err != nil {
		fmt.Printf("Failed to list RAML files: %s\n", err)
		return
	}

	report := &Report{
		Parser:  parsersMeta[*parserFl],
		Results: []*FileResult{},
		Branch:  *branchFl,
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

	SaveReport(report, *outdirFl)
}
