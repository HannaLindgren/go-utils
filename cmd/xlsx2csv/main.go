package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	// https://github.com/qax-os/excelize
	"github.com/xuri/excelize/v2"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	hio "github.com/HannaLindgren/go-utils/io"
)

func readFile(f string) ([][]string, string, error) {
	res := [][]string{}
	x, err := excelize.OpenFile(f)
	if err != nil {
		return res, "", fmt.Errorf("failed to open file : %v", err)
	}

	sheets := x.GetSheetList()
	var selectedSheet string
	if len(sheetNames) == 0 {
		if len(sheets) != 1 {
			return res, "", fmt.Errorf("multiple sheets found in %s, use -sheets flag to select which ones to export: %v", f, sheets)
		}
		selectedSheet = sheets[0]
	} else {
		var selectedSheets = []string{}
		for _, sheet := range sheets {
			if sheetNames[sheet] {
				selectedSheets = append(selectedSheets, sheet)
			}
		}
		if len(selectedSheets) == 0 {
			requestedSheetNames := maps.Keys(sheetNames)
			slices.Sort(requestedSheetNames)
			return res, "", fmt.Errorf("requested sheets %v, found: %v", requestedSheetNames, sheets)
		}
		if len(selectedSheets) > 1 {
			return res, "", fmt.Errorf("cannot select more than one sheet per file, found: %v", selectedSheets)
		}
		selectedSheet = selectedSheets[0]
	}
	//fmt.Fprintf(os.Stderr, "Using sheet %s\n", selectedSheet)

	rows, err := x.GetRows(selectedSheet)
	if err != nil {
		return res, "", fmt.Errorf("failed to read rows : %v", err)
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
	return res, selectedSheet, nil
}

func convertFile(xlsxFile, newExt string) (string, string, int, error) {

	var nLines int

	lines, selectedSheet, err := readFile(xlsxFile)
	if err != nil {
		return "", "", 0, fmt.Errorf("read failed : %v", err)
	}
	nLines = len(lines)

	ext := strings.TrimPrefix(filepath.Ext(xlsxFile), ".")
	if ext != "xlsx" {
		return "", "", nLines, fmt.Errorf("input file has invalid extension %s", xlsxFile)
	}
	outFile := fmt.Sprintf("%s.%s", hio.RemoveFileExtension(xlsxFile), newExt)
	if path.Base(outFile) == path.Base(xlsxFile) {
		return "", "", nLines, fmt.Errorf("input and output file are have the same extension: %s", xlsxFile)
	}
	outWriter, err := os.Create(outFile)
	if err != nil {
		return "", "", nLines, fmt.Errorf("Failed to create file: %v", err)
	}
	defer outWriter.Close()

	for i, fs := range lines {
		for _, f := range fs {
			if strings.Contains(f, fieldSep) {
				msg := fmt.Sprintf("Input field <%s> on line %v contains field sep", f, i+1)
				panic(msg)
			}
			if strings.Contains(f, "\n") {
				msg := fmt.Sprintf("Input field <%s> on line %v contains newline", f, i+1)
				panic(msg)
			}
		}
		l := strings.Join(fs, fieldSep)
		outWriter.WriteString(l)
		outWriter.WriteString("\n")
	}
	return outFile, selectedSheet, nLines, nil
}

const cmdname = "xlsx2csv"

// flags
var fieldSep string
var sheetNames = map[string]bool{}

func main() {

	fieldSepFlag := flag.String("sep", "<tab>", "field `separator`")
	sheetNamesFlag := flag.String("sheets", "", "Sheet `names` to export (comma-separated list)")
	ext := flag.String("ext", "csv", "output `extension`")

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, "Convert xlsx file sto tsv/csv")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Usage: %s <files>\n", cmdname)
		fmt.Fprintln(os.Stderr, "       OR")
		fmt.Fprintf(os.Stderr, "       cat <files> | %s\n", cmdname)
		fmt.Fprintln(os.Stderr, "\nOptional flags:")
		flag.PrintDefaults()
	}

	flag.Usage = func() {
		printUsage()
		os.Exit(0)
	}

	flag.Parse()

	files := flag.Args()
	if flag.NArg() == 0 {
		var err error
		files, err = hio.ReadStdinToLines()
		for i, f := range files {
			files[i] = strings.TrimSpace(f)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read stdin: %v\n", err)
			os.Exit(1)
		}
	}
	if len(files) == 0 {
		printUsage()
		os.Exit(0)
	}

	fieldSep = *fieldSepFlag
	if fieldSep == "<tab>" {
		fieldSep = "\t"
	}
	for _, s := range strings.Split(*sheetNamesFlag, ",") {
		s = strings.TrimSpace(s)
		sheetNames[s] = true
	}

	for _, f := range files {
		outFile, selectedSheet, n, err := convertFile(f, *ext)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Convert failed: %v\n", err)
			os.Exit(1)
		}
		if selectedSheet != "" {
			fmt.Printf("%s [%s] => %s (%v lines)\n", f, selectedSheet, outFile, n)
		} else {
			fmt.Printf("%s => %s (%v lines)\n", f, outFile, n)
		}
	}

}
