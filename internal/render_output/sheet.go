package render_output

import (
	"sort"

	"golang.org/x/exp/maps"
	"github.com/xuri/excelize/v2"
	"github.com/cloudfoundry/go-cfclient/v3/resource"

)

func WriteSheet(singletonApps map[string]map[string]map[string][]*resource.Process, outputFile string) (err error) {
	tableArray := buildTableArray(singletonApps)
	err = renderSheet(tableArray, outputFile)
	return
}

func renderSheet(tableArray [][]string, outputFile string) (err error) {

	headers := []string{
		"Org Name",
		"Space Name",
		"App Name",
		"App ID",
		"Process Type",
	}

	columnWidths := []float64{
		20, 20, 20, 32, 15,
	}
	
	f := excelize.NewFile()
	defer func() {
		if err = f.Close(); err != nil {
			return
		}
	}()

	_, err = f.NewSheet("singleton-apps")
	if err != nil {
		return
	}
	setColumnWidths(f, "singleton-apps", columnWidths)
	
	// Write headers
	writeLine(f, "singleton-apps", headers, 0)
	
	// Write lines from array
	for row, line := range tableArray{
		writeLine(f, "singleton-apps", line, row + 1)
	}
	
	_ = f.DeleteSheet("Sheet1")

	err = f.SaveAs(outputFile)
	return
}

func buildTableArray(singletonApps map[string]map[string]map[string][]*resource.Process) (tableArray [][]string) {
	// Get keys to enable sorting so that each sheet has a predictable order
	orgs := maps.Keys(singletonApps)
	sort.Strings(orgs)
	for _, orgName := range orgs {
		spaces := maps.Keys(singletonApps[orgName])
		sort.Strings(spaces)
		for _, spaceName := range spaces {
			apps := maps.Keys(singletonApps[orgName][spaceName])
			sort.Strings(apps)
			for _, appName := range apps {
				for _, process := range singletonApps[orgName][spaceName][appName] {
					tableArray = append(tableArray, []string{
						orgName,
						spaceName,
						appName,
						process.Relationships.App.Data.GUID,
						process.Type,
					})
				}
			}
		}
	}
	return
}

// Set widths of colums
func setColumnWidths(f *excelize.File, sheetName string, columnWidths []float64) {
	for columnIdx, columnWidth := range columnWidths{
		columnName, _ := excelize.ColumnNumberToName(columnIdx + 1)
		_ = f.SetColWidth(sheetName, columnName, columnName, columnWidth)
	}
}

// Write line to a worksheet based on an array of strings
func writeLine(f *excelize.File, sheetName string, content []string, rowIdx int) {
	for columnIdx, cellContent := range content {
		cellName, _ := excelize.CoordinatesToCellName(columnIdx + 1, rowIdx + 1)
		_ = f.SetCellValue(sheetName, cellName, cellContent)
	}
}
