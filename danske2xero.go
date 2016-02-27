// Copyright 2016 Tomislav Biscan. All rights reserved.
// Use of this source code is governed by a MIT license
// The license can be found in the LICENSE file.

// Danske2Xero is a simple utility tool that converts Danske Bank statements into Xero accounting software CSV.

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/gocarina/gocsv"
	"io"
	"os"
	"strings"
	"time"
)

type DateTime struct {
	time.Time
}

// Convert the internal date as CSV string
func (date *DateTime) MarshalCSV() (string, error) {
	return date.Time.Format("02/01/2006"), nil
}

// Convert the CSV string as internal date
func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("01/02/2006", csv)
	if err != nil {
		return err
	}
	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// "Booked date";"Interest date";"Text";"Number";"Amount in GBP";"Booked balance in GBP";"Status";"Bank's archive reference"
// *Date,*Amount,Payee,Description,Reference,Cheque Number
type Statement struct {
	Date        DateTime `csv:"Booked date"`
	Amount      string   `csv:"Amount in GBP"`
	Payee       string   `csv:""`
	Description string   `csv:"Text"`
	Reference   string   `csv:"Bank's archive reference"`
	CheckNumber string   `csv:""`
}

func main() {

	filename := ""
	if len(os.Args) < 2 {
		fmt.Println("Please provide filename as an argument")
		os.Exit(1)
	} else {
		filename = os.Args[1]
	}

	statementFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	check(err)
	defer statementFile.Close()

	statements := []*Statement{}

	gocsv.SetCSVReader(func(in io.Reader) *csv.Reader {
		reader := gocsv.LazyCSVReader(in) // Allows use of quotes in CSV
		reader.Comma = ';'
		return reader
	})

	if err := gocsv.UnmarshalFile(statementFile, &statements); err != nil { // Load statements from file
		panic(err)
	}

	csvContent, err := gocsv.MarshalString(&statements) // Get all statements as CSV string
	check(err)
	lines := strings.Split(csvContent, "\n")
	lines[0] = "*Date,*Amount,Payee,Description,Reference,Cheque Number"
	csvContent = strings.Join(lines, "\n")
	numStatements := len(lines) - 1

	outputFilename := strings.Replace(filename, ".csv", "_output.csv", -1)
	f, err := os.Create(outputFilename)
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	bytes, err := w.WriteString(csvContent)
	check(err)
	fmt.Printf("=== Wrote %d bytes to %s ======\n", bytes, outputFilename)
	w.Flush()

	fmt.Println(csvContent)
	fmt.Printf("=== %d staments converted ======\n", numStatements)
}
