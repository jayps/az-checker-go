package excel

import (
	"azurechecker/azure"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
)

func addHeading(f *excelize.File, sheet string, cellLocation string, text string) {
	boldStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	err = f.SetCellValue(sheet, cellLocation, text)
	if err != nil {
		log.Fatalln("Failed to set cell value at", cellLocation, "to", text)
	}

	err = f.SetCellStyle(sheet, cellLocation, cellLocation, boldStyle)
	if err != nil {
		log.Fatalln("Failed to set cell style at", cellLocation)
	}
}

func writeCell(f *excelize.File, sheet string, cellLocation string, text string) {
	err := f.SetCellValue(sheet, cellLocation, text)
	if err != nil {
		log.Fatalln("Failed to set cell value at", cellLocation, "to", text)
	}
}

func writeResourceAlerts(f *excelize.File, title string, lineIndex *int, resources map[string]azure.Resource) {
	addHeading(f, "Alerts", fmt.Sprintf("A%d", *lineIndex), title)
	addHeading(f, "Alerts", fmt.Sprintf("B%d", *lineIndex), fmt.Sprintf("%d resources checked", len(resources)))
	*lineIndex++

	for _, resource := range resources {
		writeCell(f, "Alerts", fmt.Sprintf("A%d", *lineIndex), resource.Name)
		writeCell(f, "Alerts", fmt.Sprintf("B%d", *lineIndex), fmt.Sprintf("%d alerts configured", len(resource.AlertRules)))
		*lineIndex++
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

	// Add VMs to Alerts sheet
	f.SetSheetName("Sheet1", "Alerts")
	lineIndex := 1
	writeResourceAlerts(f, "Virtual machines", &lineIndex, vms)
	writeResourceAlerts(f, "AKS clusters", &lineIndex, aksClusters)
	writeResourceAlerts(f, "MySQL servers", &lineIndex, mySQLServers)
	writeResourceAlerts(f, "Flexible MySQL servers", &lineIndex, flexibleMySQLServers)
	writeResourceAlerts(f, "SQL servers", &lineIndex, sqlServers)
	writeResourceAlerts(f, "Storage accounts", &lineIndex, storageAccounts)
	writeResourceAlerts(f, "Web apps", &lineIndex, webApps)

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
