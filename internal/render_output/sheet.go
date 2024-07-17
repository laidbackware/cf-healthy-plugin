package render_output

import (
	"sort"
	"strconv"

	"golang.org/x/exp/maps"
	"github.com/xuri/excelize/v2"
	"github.com/laidbackware/cf-healthy-plugin/internal/collect_data"

)

type sheetContent struct {
	sheetName			string
	sheetHeaders		[]string
	columnWidths	[]float64
	tableData			[][]string
}

func WriteSheet(healthState collect_data.HealthState, outputFile string) (err error) {
	allSheetsContents := []sheetContent{
		{
			sheetName: "singleton_apps",
			sheetHeaders: []string{
				"Org Name",
				"Space Name",
				"App Name",
				"App ID",
				"Process Type",
			},
			columnWidths: []float64{20, 20, 20, 32, 15},
			tableData: buildTableArray(healthState.SingletonApps),
		},
	}
	
	// tableArray := buildTableArray(healthState.SingletonApps)
	err = renderSheet(allSheetsContents, outputFile)
	return
}

func renderSheet(allSheetsContents []sheetContent, outputFile string) (err error) {
	// Initialize file
	f := excelize.NewFile()
	defer func() {
		if err = f.Close(); err != nil {
			return
		}
	}()

	// Walk sheets and render
	for _, sheetContents := range allSheetsContents {
		_, err = f.NewSheet("singleton-apps")
		if err != nil {
			return
		}
		setColumnWidths(f, "singleton-apps", sheetContents.columnWidths)
		
		// Write headers
		writeLine(f, "singleton-apps", sheetContents.sheetHeaders, 0)
		
		// Write lines from array
		for row, line := range sheetContents.tableData{
			writeLine(f, "singleton-apps", line, row + 1)
		}
		
		
	}
	_ = f.DeleteSheet("Sheet1")
	err = f.SaveAs(outputFile)
	return

}

func buildTableArray(sourceData map[string]map[string]map[string][]collect_data.Process) (tableArray [][]string) {
	// Get keys to enable sorting so that each sheet has a predictable order
	orgs := maps.Keys(sourceData)
	sort.Strings(orgs)
	idx := 0
	for _, orgName := range orgs {
		spaces := maps.Keys(sourceData[orgName])
		sort.Strings(spaces)
		for _, spaceName := range spaces {
			apps := maps.Keys(sourceData[orgName][spaceName])
			sort.Strings(apps)
			for _, appName := range apps {
				for _, process := range sourceData[orgName][spaceName][appName] {
					tableArray = append(tableArray, []string{
						orgName,
						spaceName,
						appName,
						process.AppGuid,
						process.Type,
						strconv.Itoa(process.Instances),

					})
					if process.HealthCheck != nil {
						tableArray[idx] = append(tableArray[idx], []string{
							process.HealthCheck.Type,
							strconv.Itoa(*process.HealthCheck.Data.InvocationTimeout),
							strconv.Itoa(*process.HealthCheck.Data.Timeout),
						}...)
					}
					idx++
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
