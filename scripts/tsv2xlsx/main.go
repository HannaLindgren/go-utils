package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/HannaLindgren/go-utils/io"
)

const fieldSep = "\t"
const sheetName = "Sheet1"

// ColumnNumberToName provides a function to convert the integer to Excel sheet column title.
// Code copied from https://github.com/360EntSecGroup-Skylar/excelize/blob/v2.4.0/lib.go#L164-L177
// since I failed to install the latest version of excelize
func ColumnNumberToName(num int) (string, error) {
	if num < 1 {
		return "", fmt.Errorf("incorrect column number %d", num)
	}
	if num > 100 /*TotalColumns*/ {
		return "", fmt.Errorf("column number exceeds maximum limit")
	}
	var col string
	for num > 0 {
		col = string(rune((num-1)%26+65)) + col
		num = (num - 1) / 26
	}
	return col, nil
}

const cmd = "tsv2xlsx"

func main() {
	header := flag.Bool("header", false, "Data has header (default false)")
	help := flag.Bool("help", false, "Print usage and exit")
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s <flags> <files>\n", cmd)
		fmt.Fprintf(os.Stderr, "FLAGS:\n")
		flag.PrintDefaults()
	}
	if *help || len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	for _, tsvFile := range flag.Args() {
		baseName := strings.TrimSuffix(tsvFile, filepath.Ext(tsvFile))
		xlsxFile := fmt.Sprintf("%s.xlsx", baseName)

		content, err := io.ReadFileToString(tsvFile)
		if err != nil {
			log.Fatalf("Failed to read file %s: %v", tsvFile, err)
		}
		xlsx := excelize.NewFile()
		var headerStyle int
		if *header {
			headerStyle, err = xlsx.NewStyle(`{"font":{"bold":true}}`)
			if err != nil {
				log.Fatalf("Failed to create header style: %v", err)
			}
		}

		lines := strings.Split(content, "\n")

		index := xlsx.NewSheet(sheetName)

		// Set column width
		// xlsx.SetColWidth(sheetName, "C", "C", 82)
		// xlsx.SetColWidth(sheetName, "E", "E", 82)

		for li, l := range lines {
			for fi, f := range strings.Split(l, fieldSep) {
				row := li + 1
				col, err := ColumnNumberToName(fi + 1)
				if err != nil {
					log.Fatalf("Failed to convert column number %v to name: %v", fi, err)
				}
				cell := fmt.Sprintf("%s%v", col, row)
				xlsx.SetCellValue(sheetName, cell, f)
				if li == 0 && *header {
					xlsx.SetCellStyle(sheetName, cell, cell, headerStyle)
				}
			}
		}

		xlsx.SetActiveSheet(index)
		err = xlsx.SaveAs(xlsxFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s => %s", tsvFile, xlsxFile)
	}
}
