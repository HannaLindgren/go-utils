package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	// https://github.com/qax-os/excelize
	"github.com/xuri/excelize/v2"

	hio "github.com/HannaLindgren/go-utils/io"
)

var cols = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

var globalFont, boldFont *excelize.Font

func initAfterFlags(fieldSepFlag, hideColsFlag, centerColsFlag *string) {
	globalFont = &excelize.Font{
		Family: *fontFamily,
		Size:   *fontSize,
		Color:  "#000000",
	}

	boldFont = &excelize.Font{
		Family: *fontFamily,
		Bold:   true,
		Size:   *fontSize,
		Color:  "#000000",
	}
	fieldSep = *fieldSepFlag
	if fieldSep == "<tab>" {
		fieldSep = "\t"
	}
	hideCols = string2ints(hideColsFlag)
	for _, i := range string2ints(centerColsFlag) {
		centerCols[i] = true
	}
}

func getColLetter(ci int) string {
	if ci >= len(cols) {
		panic(fmt.Sprintf("Column index %v is out of range", ci))
	}
	return fmt.Sprintf("%s", cols[ci])
}

func getCellID(li, ci int) string {
	return fmt.Sprintf("%s%v", getColLetter(ci), li+1)
}

func setCellStyle(sheet *excelize.File, centering, bold bool, cellID string, nLines int) error {

	var alignment *excelize.Alignment
	var font = globalFont

	if centering {
		alignment = &excelize.Alignment{Horizontal: "center"}
	}
	if bold {
		font = boldFont
	}

	style, err := sheet.NewStyle(&excelize.Style{
		Font:      font,
		Alignment: alignment,
	})
	if err != nil {
		return fmt.Errorf("new style failed : %v", err)
	}

	sheet.SetCellStyle(*sheetName, cellID, cellID, style)
	return nil
}

func setColWidth(sheet *excelize.File, col string, width float64) error {
	if err := sheet.SetColWidth(*sheetName, col, col, width); err != nil {
		return fmt.Errorf("col width failed : %v", err)
	}
	return nil
}

func setRowHeight(sheet *excelize.File, row int, height float64) error {
	if err := sheet.SetRowHeight(*sheetName, row, height); err != nil {
		return fmt.Errorf("row height failed : %v", err)
	}
	return nil
}

func hideColumns(sheet *excelize.File, cols []int) error {
	for _, i := range cols {
		err := sheet.SetColVisible(*sheetName, getColLetter(i), false)
		if err != nil {
			return err
		}
	}
	return nil
}

func convertFile(txtFile string) (string, int, error) {

	lines, err := hio.ReadFileToLines(txtFile)
	if err != nil {
		return "", 0, fmt.Errorf("read failed : %v", err)
	}
	nLines := len(lines)

	ext := strings.TrimPrefix(filepath.Ext(txtFile), ".")
	if ext != "txt" && ext != "tsv" && ext != "csv" {
		return "", nLines, fmt.Errorf("input file has invalid extension %s", txtFile)
	}
	xlsxFile := fmt.Sprintf("%s.xlsx", hio.RemoveFileExtension(txtFile))

	if path.Base(txtFile) == path.Base(xlsxFile) {
		return "", nLines, fmt.Errorf("input and output file are have the same extension: %s", txtFile)
	}

	sheet := excelize.NewFile()

	for li, l := range lines {
		fs := strings.Split(l, fieldSep)

		// set cell values
		for ci, f := range fs {
			cellID := getCellID(li, ci)
			sheet.SetCellValue(*sheetName, cellID, f)

			centering := centerCols[ci]
			bold := (li == 0) && *lockHeader
			err := setCellStyle(sheet, centering, bold, cellID, len(lines))
			if err != nil {
				return "", nLines, fmt.Errorf("cell style failed : %v", err)
			}
			if *colWidth > 0 {
				colID := getColLetter(ci)
				setColWidth(sheet, colID, *colWidth)
			}
		}
		if *rowHeight > 0 {
			setRowHeight(sheet, li, *rowHeight)
		}
	}
	// freeze first row
	if *lockHeader {
		sheet.SetPanes(*sheetName, `{"freeze":true,"split":false,"x_split":0,"y_split":1,"top_left_cell":"A2","active_pane":"bottomLeft"}`)
	}

	// hide
	err = hideColumns(sheet, hideCols)
	if err != nil {
		return "", nLines, fmt.Errorf("hide columns failed : %v", err)
	}

	if err := sheet.SaveAs(xlsxFile); err != nil {
		return "", nLines, fmt.Errorf("save failed : %v", err)
	}
	return xlsxFile, nLines, nil
}

const cmdname = "csv2xlsx"

// flags
var fieldSep string
var lockHeader *bool
var fontFamily, sheetName *string
var fontSize, colWidth, rowHeight *float64
var hideCols []int
var centerCols = make(map[int]bool)

func string2ints(cols *string) []int {
	res := []int{}
	if *cols == "" {
		return res
	}
	for _, s := range strings.Split(*cols, ",") {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Fatalf("Couldn't parse colum index %v: %v", s, err)
		}
		res = append(res, int(i))
	}
	return res
}

func main() {

	fieldSepFlag := flag.String("sep", "<tab>", "field `separator`")
	lockHeader = flag.Bool("header", false, "lock header")
	sheetName = flag.String("sheet", "Sheet1", "sheet `name`")
	hideColsFlag := flag.String("hide", "", "hide columns")
	centerColsFlag := flag.String("center", "", "center columns")
	fontFamily = flag.String("ff", "Arial", "font family")
	fontSize = flag.Float64("fs", 9, "font size")
	colWidth = flag.Float64("cw", 0, "column width")
	rowHeight = flag.Float64("rh", 0, "row height")

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, "Convert tsv/csv files to xlsx")
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

	initAfterFlags(fieldSepFlag, hideColsFlag, centerColsFlag)

	if flag.NArg() < 1 {
		printUsage()
		os.Exit(0)
	}

	for _, f := range flag.Args() {
		xFile, n, err := convertFile(f)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s => %s (%v lines)\n", f, xFile, n)
	}

}
