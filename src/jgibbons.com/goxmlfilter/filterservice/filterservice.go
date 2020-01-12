package filterservice

import (
	"fmt"
	"jgibbons.com/goxmlfilter/config"
)

const (
	APP_NAME    = "xmlfilter"
	APP_VERSION = "0.0.1"
	APP_META    = "by Jonathan Gibbons (c) 2020 All Rights Reserved"
)

func Start(config *config.FilterConfig) error {
	fmt.Printf("Appliction: %s, Version %s, %s\n", APP_NAME, APP_VERSION, APP_META)

	fmt.Printf("[%s] Calling [%s]\n", debugTStamp(), config.RestQuery)
	resp, err := callRestApi(config.RestQuery)
	if err != nil {
		return err
	}

	fmt.Printf("[%s] Processing Response of size [%d] bytes\n", debugTStamp(), resp.ContentLength)
	rowFields, fieldColNums := convertExtractsToMap(config.ExtractColumns)

	defer resp.Body.Close()
	err = decodeIOStream(resp.Body, config.DelimTag,
		config.RowsInFile,
		config.NumFiles,
		convertFiltersToMap(config.FiltersEquals),
		convertFiltersToMap(config.FiltersNotEquals),
		rowFields, fieldColNums)
	return err
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return err
	//}
	//fmt.Printf("%v", string(body))

}

// A filter could be listed for a field mutliple times.
func convertFiltersToMap(filters []config.FilterField) map[string][]string {
	m := make(map[string][]string)

	for _, fld := range filters {
		// Trick is zero value for slice can be appended to.
		// as per https://blog.golang.org/go-maps-in-action 1/2 way down.
		m[fld.Name] = append(m[fld.Name], fld.Value)
	}
	return m
}

func convertExtractsToMap(extractFields []config.ExtractField) (map[string]string, map[string]int) {
	m := make(map[string]string)
	n := make(map[string]int)

	for idx, fld := range extractFields {
		m[fld.InputName] = fld.OuputName
		n[fld.InputName] = idx
	}
	return m, n
}
