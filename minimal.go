package main

import (
	"fmt"

	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

func testfunc() {
	inputFile := "./test_files/sample.pdf"
	pdfReader, f, err := model.NewPdfReaderFromFile(inputFile, nil)
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
	contentStreamOps := pageText.GetContentStreamOps()
	textMarks := pageText.Marks()
	samplePdfObject := textMarks.Elements()[0].DirectObject
	exists := false
	for _, op := range *contentStreamOps {
		params := op.Params
		operand := op.Operand
		if operand == "TJ" || operand == "Tj" {
			for _, directObject := range params {
				if directObject == samplePdfObject {
					exists = true
				}
			}
		}
	}
	fmt.Print(exists)
}
