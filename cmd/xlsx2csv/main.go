package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	// https://github.com/qax-os/excelize
	"github.com/xuri/excelize/v2"

	hio "github.com/HannaLindgren/go-utils/io"
)

func readFile(f string) ([][]string, error) {
	res := [][]string{}
	x, err := excelize.OpenFile(f)
	if err != nil {
		return res, fmt.Errorf("failed to open file : %v", err)
	}

	sheets := x.GetSheetList()
	var selectedSheet string
	if *sheetName == "" {
		if len(sheets) != 1 {
			return res, fmt.Errorf("multiple sheets found in %s, use -sheet flag to select which one to export: %v", f, sheets)
		}
		selectedSheet = sheets[0]
	} else {
		for _, sheet := range sheets {
			if sheet == *sheetName {
				selectedSheet = sheet
				break
			}
		}
		if selectedSheet == "" {
			return res, fmt.Errorf("invalid sheet name %s, found: %v", *sheetName, sheets)
		}
	}
	//fmt.Fprintf(os.Stderr, "Using sheet %s\n", selectedSheet)

	rows, err := x.GetRows(selectedSheet)
	if err != nil {
		return res, fmt.Errorf("failed to read rows : %v", err)
	}
	var firstLineLen int
	for ri, row := range rows {
		line := []string{}
		for _, cell := range row {
			line = append(line, strings.TrimSuffix(cell, "\n")) // trim final newline if it exists
		}
		if ri == 0 {
			firstLineLen = len(line)
		}
		if ri > 0 {
			for len(line) < firstLineLen {
				line = append(line, "")
			}
		}
		res = append(res, line)
	}
	return res, nil
}

func convertFile(xlsxFile, newExt string) (string, int, error) {

	var nLines int

	ext := strings.TrimPrefix(filepath.Ext(xlsxFile), ".")
	if ext != "xlsx" {
		return "", nLines, fmt.Errorf("input file has invalid extension %s", xlsxFile)
	}
	outFile := fmt.Sprintf("%s.%s", hio.RemoveFileExtension(xlsxFile), newExt)
	if path.Base(outFile) == path.Base(xlsxFile) {
		return "", nLines, fmt.Errorf("input and output file are have the same extension: %s", xlsxFile)
	}
	outWriter, err := os.Create(outFile)
	if err != nil {
		return "", nLines, fmt.Errorf("Failed to create file: %v", err)
	}
	defer outWriter.Close()

	lines, err := readFile(xlsxFile)
	if err != nil {
		return "", 0, fmt.Errorf("read failed : %v", err)
	}
	nLines = len(lines)

	for _, fs := range lines {
		l := strings.Join(fs, fieldSep)
		outWriter.WriteString(l)
		outWriter.WriteString("\n")
	}
	return outFile, nLines, nil
}

const cmdname = "xlsx2csv"

// flags
var fieldSep string
var sheetName *string

func main() {

	fieldSepFlag := flag.String("sep", "<tab>", "field `separator`")
	sheetName = flag.String("sheet", "", "Sheet `name` to export")
	ext := flag.String("ext", "csv", "output `extension`")

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, "Convert xlsx file sto tsv/csv")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Usage: %s <files>", cmdname)
		fmt.Fprintln(os.Stderr, "\nOptional flags:")
		flag.PrintDefaults()
	}

	flag.Usage = func() {
		printUsage()
		os.Exit(0)
	}

	flag.Parse()

	if flag.NArg() < 1 {
		printUsage()
		os.Exit(0)
	}

	fieldSep = *fieldSepFlag
	if fieldSep == "<tab>" {
		fieldSep = "\t"
	}

	for _, f := range flag.Args() {
		outFile, n, err := convertFile(f, *ext)
		if err != nil {
			log.Fatal(err)
		}
		if *sheetName != "" {
			fmt.Printf("%s [Sheet:%s] => %s (%v lines)\n", f, *sheetName, outFile, n)
		} else {
			fmt.Printf("%s => %s (%v lines)\n", f, outFile, n)
		}
	}

}