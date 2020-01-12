package filterservice

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func decodeIOStream(reader io.Reader, delimeterTag string,
	rowsInFile int,
	numFiles int,
	filterEquals map[string][]string,
	filterNotEquals map[string][]string,
	extractFields map[string]string,
	extractFieldsColNum map[string]int) error {

	decoder := xml.NewDecoder(reader)

	headerLine := constructHeader(extractFields, extractFieldsColNum)
	rowCount := 0
	fileCount := 1
	opFile, err := openNewCsv(headerLine, fileCount)
	if err != nil {
		return err
	}

	var currentElement string = ""

	firstElem := true
	failedValueCheck := false
	row := make([]string, len(extractFields))
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
			el := se.Name.Local
			// If we just read a StartElement token
			// ...and its name is "page"
			if el == delimeterTag {
				if !firstElem && !failedValueCheck {
					dumpRow(row, failedValueCheck, opFile)
					rowCount = rowCount + 1
					if rowCount >= rowsInFile {
						opFile.Close()
						opFile = nil
						rowCount = 0
						fileCount += 1
						if fileCount > numFiles {
							break
						}
						opFile, err = openNewCsv(headerLine, fileCount)
						if err != nil {
							return err
						}
					}
				}
				row = make([]string, len(extractFields))
				failedValueCheck = false
				firstElem = false
			}
			currentElement = el
		case xml.CharData:
			cdata := string(se)
			//			fmt.Printf("   %s: %s\n", currentElement, cdata)
			if _, ok := extractFields[currentElement]; ok {
				colNum := extractFieldsColNum[currentElement]
				row[colNum] = cdata
			}
			if !failedValueCheck {
				failedValueCheck = discardRow(currentElement, cdata, filterEquals, filterNotEquals)
			}

		case xml.EndElement:
			currentElement = ""
		}
	}
	if !firstElem && opFile != nil {
		dumpRow(row, failedValueCheck, opFile)
		opFile.Close()
	}
	println("Processing Complete")
	return nil
}

// Generates the scv first line holding the AS values for the output columns
func constructHeader(extractFields map[string]string,
	extractFieldsColNum map[string]int) string {
	row := make([]string, len(extractFieldsColNum))
	for k, v := range extractFieldsColNum {
		row[v] = extractFields[k]
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

func debugTStamp() string {
	return time.Now().Format("15:04:05.000")
}

func openNewCsv(headerLine string, fileCount int) (*os.File, error) {
	fileNameBase := time.Now().Format("20060102_150405")
	fileName := fileNameBase + "_" + strconv.Itoa(fileCount) + ".csv"
	fmt.Printf("[%s] Writing to file [%s]\n", debugTStamp(), fileName)

	opFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	opFile.WriteString(headerLine + "\n")
	return opFile, err
}

func discardRow(elName string, elValue string,
	filterEquals map[string][]string,
	filterNotEquals map[string][]string) bool {

	failedToEqual := false
	if eqVals, ok := filterEquals[elName]; ok {
		for _, comp := range eqVals {
			if comp != elValue {
				failedToEqual = true
				break
			}
		}
	}

	failedToNotEqual := false
	if eqVals, ok := filterNotEquals[elName]; ok {
		for _, comp := range eqVals {
			if comp == elValue {
				failedToNotEqual = true
				break
			}
		}
	}
	return failedToEqual || failedToNotEqual
}

func dumpRow(row []string, discardRow bool, opfile *os.File) {
	//	fmt.Printf("discard [%t] row: %v\n", discardRow, row)
	if !discardRow {
		for idx, cell := range row {
			if idx > 0 {
				opfile.WriteString(",")
			}
			opfile.WriteString(cell)
		}
		opfile.WriteString("\n")
	}
}
