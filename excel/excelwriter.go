package excel

import (
	"fmt"
	"github.com/jayps/azure-checker-go/azure"
	"github.com/xuri/excelize/v2"
	"log"
)

func addStyledCell(f *excelize.File, sheet string, cellLocation string, text string, style int) {
	err := f.SetCellValue(sheet, cellLocation, text)
	if err != nil {
		log.Fatalln("Failed to set cell value at", cellLocation, "to", text, "on sheet", sheet, ". Error: ", err.Error())
	}

	err = f.SetCellStyle(sheet, cellLocation, cellLocation, style)
	if err != nil {
		log.Fatalln("Failed to set cell style at", cellLocation)
	}
}

func addHeading(f *excelize.File, sheet string, cellLocation string, text string) {
	boldStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 20,
		},
	})

	if err != nil {
		fmt.Println(err)
	}

	addStyledCell(f, sheet, cellLocation, text, boldStyle)
}

func addBoldCell(f *excelize.File, sheet string, cellLocation string, text string) {
	boldStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})

	if err != nil {
		fmt.Println(err)
	}

	addStyledCell(f, sheet, cellLocation, text, boldStyle)
}

func writeCell(f *excelize.File, sheet string, cellLocation string, text string) {
	err := f.SetCellValue(sheet, cellLocation, text)
	if err != nil {
		log.Fatalln("Failed to set cell value at", cellLocation, "to", text, "on sheet", sheet, ". Error: ", err.Error())
	}
}

func writeResourceAlerts(f *excelize.File, sheet string, resources map[string]azure.Resource) {
	lineIndex := 1
	addHeading(f, sheet, fmt.Sprintf("A%d", lineIndex), fmt.Sprintf("%d resources found", len(resources)))
	lineIndex++

	for _, resource := range resources {
		addBoldCell(f, sheet, fmt.Sprintf("A%d", lineIndex), resource.Name)
		addBoldCell(f, sheet, fmt.Sprintf("B%d", lineIndex), fmt.Sprintf("%d alerts configured", len(resource.AlertRules)))
		lineIndex++

		if len(resource.AlertRules) > 0 {
			addBoldCell(f, sheet, fmt.Sprintf("B%d", lineIndex), "Name")
			addBoldCell(f, sheet, fmt.Sprintf("C%d", lineIndex), "Criteria")
			lineIndex++

			for _, alertRule := range resource.AlertRules {
				writeCell(f, sheet, fmt.Sprintf("B%d", lineIndex), alertRule.Name)

				for _, criterion := range alertRule.Criteria.AllOf {
					criterionOutput := fmt.Sprintf("%s %s %s %s", criterion.TimeAggregation,
						criterion.MetricName,
						criterion.Operator,
						fmt.Sprintf("%.2f", criterion.Threshold),
					)
					writeCell(f, sheet, fmt.Sprintf("C%d", lineIndex), criterionOutput)
					lineIndex++
				}
			}
		}
	}
}

func writeResourceBackups(f *excelize.File, vms map[string]azure.Resource) {
	lineIndex := 1
	for _, vm := range vms {
		writeCell(f, "Backups", fmt.Sprintf("A%d", lineIndex), vm.Name)
		if vm.BackupVault != nil {
			writeCell(f, "Backups", fmt.Sprintf("B%d", lineIndex), vm.BackupVault.Name)
		} else {
			writeCell(f, "Backups", fmt.Sprintf("B%d", lineIndex), "Not backed up")
		}
		lineIndex++
	}
}

func writeRecommendations(f *excelize.File, recommendations map[string][]azure.AdvisorRecommendation) {
	for category, categoryRecommendations := range recommendations {
		sheetTitle := fmt.Sprintf("%s Recs", category)
		f.NewSheet(sheetTitle)
		index := 1

		writeCell(f, sheetTitle, fmt.Sprintf("A%d", index), "Recommendation")
		writeCell(f, sheetTitle, fmt.Sprintf("B%d", index), "Impact")
		writeCell(f, sheetTitle, fmt.Sprintf("C%d", index), "Resource type")
		writeCell(f, sheetTitle, fmt.Sprintf("D%d", index), "Affected resource")
		writeCell(f, sheetTitle, fmt.Sprintf("E%d", index), "Resource group")

		index++

		for i := 0; i < len(categoryRecommendations); i++ {
			writeCell(f, sheetTitle, fmt.Sprintf("A%d", index+1), categoryRecommendations[i].Description.Problem)
			writeCell(f, sheetTitle, fmt.Sprintf("B%d", index+1), categoryRecommendations[i].Impact)
			writeCell(f, sheetTitle, fmt.Sprintf("C%d", index+1), categoryRecommendations[i].ResourceType)
			writeCell(f, sheetTitle, fmt.Sprintf("D%d", index+1), categoryRecommendations[i].AffectedResource)
			writeCell(f, sheetTitle, fmt.Sprintf("E%d", index+1), categoryRecommendations[i].ResourceGroup)
			index++
		}
	}
}

func OutputExcelDocument(
	filename string,
	vms map[string]azure.Resource,
	aksClusters map[string]azure.Resource,
	mySQLServers map[string]azure.Resource,
	flexibleMySQLServers map[string]azure.Resource,
	sqlServers map[string]azure.Resource,
	storageAccounts map[string]azure.Resource,
	webApps map[string]azure.Resource,
	recommendations map[string][]azure.AdvisorRecommendation,
) {
	f := excelize.NewFile()

	if len(vms) > 0 {
		f.SetSheetName("Sheet1", "VM Alerts")
		writeResourceAlerts(f, "VM Alerts", vms)
	} else {
		f.DeleteSheet("Sheet1")
	}

	if len(aksClusters) > 0 {
		f.NewSheet("AKS Cluster Alerts")
		writeResourceAlerts(f, "AKS Cluster Alerts", aksClusters)
	}

	if len(mySQLServers) > 0 {
		f.NewSheet("MySQL Server Alerts")
		writeResourceAlerts(f, "MySQL Server Alerts", mySQLServers)
	}

	if len(flexibleMySQLServers) > 0 {
		f.NewSheet("Flexible MySQL Server Alerts")
		writeResourceAlerts(f, "Flexible MySQL Server Alerts", flexibleMySQLServers)
	}

	if len(sqlServers) > 0 {
		f.NewSheet("SQL Server Alerts")
		writeResourceAlerts(f, "SQL Server Alerts", sqlServers)
	}

	if len(storageAccounts) > 0 {
		f.NewSheet("Storage Account Alerts")
		writeResourceAlerts(f, "Storage Account Alerts", storageAccounts)
	}

	if len(webApps) > 0 {
		f.NewSheet("Web App Alerts")
		writeResourceAlerts(f, "Web App Alerts", webApps)
	}

	// Backups
	f.NewSheet("Backups")
	writeResourceBackups(f, vms)

	writeRecommendations(f, recommendations)

	if err := f.SaveAs(fmt.Sprintf("%s.xlsx", filename)); err != nil {
		fmt.Println("Could not save output document", err)
	} else {
		fmt.Println("Saved output document")
	}
}
