package filterservice

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type XmlDecoder struct {
	delimeterTag        string
	rowsInFile          int
	numFiles            int
	debugOutput         bool
	filterEquals        map[string][]string
	filterNotEquals     map[string][]string
	extractFields       map[string]string
	extractFieldsColNum map[string]int

	// the ones used during processing
	pathStack        []string
	firstElem        bool
	failedValueCheck bool
	rowCount         int
	fileCount        int
	opFile           *os.File
}

func NewDecoder(delimeterTag string,
	rowsInFile int,
	numFiles int,
	debugOutput bool,
	filterEquals map[string][]string,
	filterNotEquals map[string][]string,
	extractFields map[string]string,
	extractFieldsColNum map[string]int) *XmlDecoder {
	dec := XmlDecoder{delimeterTag,
		rowsInFile,
		numFiles,
		debugOutput,
		filterEquals,
		filterNotEquals,
		extractFields,
		extractFieldsColNum,
		make([]string, 0),
		true,
		false,
		0,
		0,
		nil,
	}
	return &dec
}

func (dec *XmlDecoder) decodeIOStream(reader io.Reader) error {
	decoder := xml.NewDecoder(reader)

	headerLine := dec.constructHeader()
	newFile, err := dec.openNewCsv(headerLine)
	if err != nil {
		return err
	}
	dec.opFile = newFile

	row := make([]string, len(dec.extractFields))
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		// Inspect the type of the token just read.
		//		ty := reflect.TypeOf(t)
		//		fmt.Printf("%v\n", ty)
		switch se := t.(type) {
		case xml.StartElement:
			stopNow, err := dec.handleStartElement(se.Name.Local, row, headerLine)
			if err != nil {
				return err
			}
			if stopNow {
				break
			}
		case xml.CharData:
			dec.handleCharData(string(se), row)
		case xml.EndElement:
			if n := len(dec.pathStack) - 1; n > 0 {
				dec.pathStack = dec.pathStack[:n] // pop
			}
		}
	}
	if !dec.firstElem && dec.opFile != nil {
		dec.dumpRow(row)
		dec.opFile.Close()
	}
	println("Processing Complete")
	return nil
}

// Return true if should stop processing as output enough as per config
func (dec *XmlDecoder) handleStartElement(el string, row []string, headerLine string) (bool, error) {
	// If we just read a StartElement token
	// ...and its name is "page"
	if el == dec.delimeterTag {
		if !dec.firstElem && !dec.failedValueCheck {
			err := dec.dumpRow(row)
			if err != nil {
				return false, err
			}
			dec.rowCount += 1
			if dec.rowCount >= dec.rowsInFile {
				dec.opFile.Close()
				dec.opFile = nil
				dec.rowCount = 0
				if dec.fileCount >= dec.numFiles {
					return true, nil
				}
				var err error
				dec.opFile, err = dec.openNewCsv(headerLine)
				if err != nil {
					return false, err
				}
			}
		}
		row = make([]string, len(dec.extractFields))
		dec.failedValueCheck = false
		dec.firstElem = false
	}
	dec.pathStack = append(dec.pathStack, el) // push

	return false, nil
}

func (dec *XmlDecoder) handleCharData(cdata string, row []string) {
	path := dec.createPath()
	if dec.debugOutput {
		fmt.Printf("   %s: %s\n", path, cdata)
	}
	if colNum := dec.comparePathToFields(path); colNum >= 0 {
		row[colNum] = cdata
	}
	if !dec.failedValueCheck {
		dec.failedValueCheck = dec.discardRow(path, cdata)
	}
}

// Compares the current full path to the defined columns to extract.
// The columns to extract could just be the very tail of the path rather than the full one
// eg Path could be /obj/something/andanother/item but the extract could just have andanother/item
// return -1 if no match, else the index of the column to extract
func (dec *XmlDecoder) comparePathToFields(path string) int {

	for k, _ := range dec.extractFields {
		if strings.HasSuffix(path, k) {
			colNum := dec.extractFieldsColNum[k]
			return colNum
		}
	}
	return -1
}

func (dec *XmlDecoder) createPath() string {
	ret := ""
	for _, p := range dec.pathStack {
		ret += "/" + p
	}
	return ret
}

// Generates the scv first line holding the AS values for the output columns
func (dec *XmlDecoder) constructHeader() string {
	row := make([]string, len(dec.extractFieldsColNum))
	for k, v := range dec.extractFieldsColNum {
		row[v] = dec.extractFields[k]
	}
	ret := ""
	for _, colName := range row {
		if len(ret) > 0 {
			ret += ","
		}
		ret += colName
	}
	return ret
}

// path is of the form el1/el2/el3
// The filters may only include the end part of the path, eg el2/el3
func (dec *XmlDecoder) discardRow(path string, elValue string) bool {

	failedToEqual := false
	if eqVals, ok := findMatchingRule(path, dec.filterEquals); ok {
		for _, comp := range eqVals {
			if comp != elValue {
				failedToEqual = true
				break
			}
		}
	}

	failedToNotEqual := false
	if neqVals, ok := findMatchingRule(path, dec.filterNotEquals); ok {
		for _, comp := range neqVals {
			if comp == elValue {
				failedToNotEqual = true
				break
			}
		}
	}
	return failedToEqual || failedToNotEqual
}

func (dec *XmlDecoder) openNewCsv(headerLine string) (*os.File, error) {
	dec.fileCount += 1
	fileNameBase := time.Now().Format("20060102_150405")
	fileName := fileNameBase + "_" + strconv.Itoa(dec.fileCount) + ".csv"
	fmt.Printf("[%s] Writing to file [%s]\n", debugTStamp(), fileName)

	opFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	_, err = opFile.WriteString(headerLine + "\n")
	return opFile, err
}

func (dec *XmlDecoder) dumpRow(row []string) error {
	var err error = nil
	//	fmt.Printf("discard [%t] row: %v\n", discardRow, row)
	if !dec.failedValueCheck {
		for idx, cell := range row {
			if idx > 0 {
				_, err = dec.opFile.WriteString(",")
			}
			if err == nil {
				_, err = dec.opFile.WriteString(cell)
			}
		}
		if err == nil {
			_, err = dec.opFile.WriteString("\n")
		}
	}
	return err
}
