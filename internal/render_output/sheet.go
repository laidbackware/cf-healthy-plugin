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
	headersWithHealth := []string{
		"Org Name",
		"Space Name",
		"App Name",
		"App ID",
		"Process Type",
		"Instances",
		"Health Check Type",
		"Interval",
		"Timeout",
		"Invocation Timeout",
		"HTTP Endpoint",
	}
	widthsWithHealth := []float64{20, 20, 20, 32, 15, 15, 25, 20, 15, 25, 20}
	allSheetsContents := []sheetContent{
		{
			sheetName: "all_apps",
			sheetHeaders: headersWithHealth,
			columnWidths: widthsWithHealth,
			tableData: buildTableArray(healthState.AllProcesses, true),
		},
		{
			sheetName: "singleton_apps",
			sheetHeaders: []string{
				"Org Name",
				"Space Name",
				"App Name",
				"App ID",
				"Process Type",
				"Instances",
			},
			columnWidths: []float64{20, 20, 20, 32, 15, 15},
			tableData: buildTableArray(healthState.SingletonApps, false),
		},
		{
			sheetName: "port_health_check",
			sheetHeaders: headersWithHealth,
			columnWidths: widthsWithHealth,
			tableData: buildTableArray(healthState.PortHealthCheck, true),
		},
		{
			sheetName: "default_http_check",
			sheetHeaders: headersWithHealth,
			columnWidths: widthsWithHealth,
			tableData: buildTableArray(healthState.DefaultHttpTime, true),
		},
	}
	
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
		_, err = f.NewSheet(sheetContents.sheetName)
		if err != nil {
			return
		}
		setColumnWidths(f, sheetContents.sheetName, sheetContents.columnWidths)
		
		// Write headers
		writeLine(f, sheetContents.sheetName, sheetContents.sheetHeaders, 0)
		
		// Write lines from array
		for row, line := range sheetContents.tableData{
			writeLine(f, sheetContents.sheetName, line, row + 1)
		}
		
		
	}
	_ = f.DeleteSheet("Sheet1")
	err = f.SaveAs(outputFile)
	return

}

func buildTableArray(sourceData map[string]map[string]map[string][]collect_data.Process, incHealth bool) (tableArray [][]string) {
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
					if incHealth {
						addHealth(tableArray, idx, process)
					}
					idx++
				}
			}
		}
	}
	return
}

func addHealth(tableArray [][]string, idx int, process collect_data.Process){
	appendArray(tableArray, idx, process.HealthCheck.Type)
	appendIntIfNotNil(tableArray, idx, process.HealthCheck.Data.Interval)
	appendIntIfNotNil(tableArray, idx, process.HealthCheck.Data.Timeout)
	appendIntIfNotNil(tableArray, idx, process.HealthCheck.Data.InvocationTimeout)
	if process.HealthCheck.Type == "http" {
		appendArray(tableArray, idx, *process.HealthCheck.Data.Endpoint)
	}
}

// Call appendArray if pointer is not nil
func appendIntIfNotNil(tableArray [][]string, idx int, content *int){
	if content != nil {
		appendArray(tableArray, idx, strconv.Itoa(*content))
		return
	}
	appendArray(tableArray, idx, "")
}

// Append string to 3d string array
func appendArray(tableArray [][]string, idx int, content string){
	tableArray[idx] = append(tableArray[idx], []string{
		content,
	}...)
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
