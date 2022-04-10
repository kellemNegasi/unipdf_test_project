package main

import (
	"fmt"
	"regexp"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/contentstream"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/redactor"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(`d3ffb0da511a9d03403559278c5b167b2e308dc5b5b9d563e84974fee3ef4f8f`)
	if err != nil {
		fmt.Printf("ERROR: Failed to set metered key: %v\n", err)
		fmt.Printf("Make sure to get a valid key from https://cloud.unidoc.io\n")
		panic(err)
	}
}

func main() {
	inputFile := "./test_files/sample.pdf"

	destFile := "./output/redacted_file_new.pdf"
	// list of regex patterns and replacement strings
	redact(inputFile, destFile)
	printContentStream(inputFile)
	// groupTextBlocks(inputFile)
}
func redact(inputFile, destFile string) {
	patterns := []string{
		//`scrambled`,
		`Virtual Mechanics`,
	}

	replacements := []string{"r"}
	terms := []redactor.RedactionTerm{}
	for i, pattern := range patterns {
		regexp, err := regexp.Compile(pattern)
		if err != nil {
			panic(err)
		}
		replacement := replacements[i]
		redTerm := redactor.RedactionTerm{Pattern: regexp, Replacement: replacement}
		terms = append(terms, redTerm)
	}
	pdfReader, f, err := model.NewPdfReaderFromFile(inputFile, nil)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	//define RedactionOptions
	options := redactor.RedactionOptions{Terms: terms}
	red, err := redactor.New(pdfReader, &options)
	if err != nil {
		panic(err)
	}
	err = red.Redact()
	if err != nil {
		panic(err)
	}
	// write the redacted document to file
	err = red.WriteToFile(destFile)
	if err != nil {
		panic(err)
	}
}
func printContentStream(file string) {
	pdfReader, f, err := model.NewPdfReaderFromFile(file, nil)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	page, err := pdfReader.GetPage(1)
	if err != nil {
		panic(err)
	}
	ex, err := extractor.New(page)
	if err != nil {
		panic(err)
	}
	pageText, _, _, err := ex.ExtractPageText()
	if err != nil {
		panic(err)
	}
	ContentStream := pageText.GetContentStreamOps().String()
	fmt.Print(ContentStream)

}

func groupTextBlocks(file string) {
	pdfReader, f, err := model.NewPdfReaderFromFile(file, nil)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	page, err := pdfReader.GetPage(1)
	if err != nil {
		panic(err)
	}
	ex, err := extractor.New(page)
	if err != nil {
		panic(err)
	}
	pageText, _, _, err := ex.ExtractPageText()
	if err != nil {
		panic(err)
	}
	ContentStream := pageText.GetContentStreamOps()
	groups := redactor.TextBlocksFromContentStreamOps(ContentStream)
	fmt.Print(groups)
}
func findPdfObject(targetObj core.PdfObject, csops *contentstream.ContentStreamOperations) bool {
	equal := false
	for _, op := range *csops {
		operand := op.Operand
		if operand == "Tj" || operand == "TJ" {
			params := op.Params
			for _, directObject := range params {
				equal = directObject == targetObj
			}
		}

	}
	return equal
}
