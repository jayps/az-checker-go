package excel

import (
	"fmt"
	"github.com/jayps/azure-checker-go/azure"
	"github.com/xuri/excelize/v2"
)

func addStyledCell(f *excelize.File, sheet string, cellLocation string, text string, style int) error {
	err := f.SetCellValue(sheet, cellLocation, text)
	if err != nil {
		return err
	}

	err = f.SetCellStyle(sheet, cellLocation, cellLocation, style)
	if err != nil {
		return err
	}

	return nil
}

func addHeading(f *excelize.File, sheet string, cellLocation string, text string) error {
	boldStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 20,
		},
	})

	if err != nil {
		return err
	}

	return addStyledCell(f, sheet, cellLocation, text, boldStyle)
}

func addBoldCell(f *excelize.File, sheet string, cellLocation string, text string) error {
	boldStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})

	if err != nil {
		return err
	}

	return addStyledCell(f, sheet, cellLocation, text, boldStyle)
}

func writeCell(f *excelize.File, sheet string, cellLocation string, text string) error {
	err := f.SetCellValue(sheet, cellLocation, text)
	if err != nil {
		return err
	}

	return nil
}

func writeResourceAlerts(f *excelize.File, sheet string, resources map[string]azure.Resource) error {
	lineIndex := 1
	err := addHeading(f, sheet, fmt.Sprintf("A%d", lineIndex), fmt.Sprintf("%d resources found", len(resources)))
	if err != nil {
		return err
	}
	lineIndex++

	for _, resource := range resources {
		err = addBoldCell(f, sheet, fmt.Sprintf("A%d", lineIndex), resource.Name)
		if err != nil {
			return err
		}
		err = addBoldCell(f, sheet, fmt.Sprintf("B%d", lineIndex), fmt.Sprintf("%d alerts configured", len(resource.AlertRules)))
		if err != nil {
			return err
		}
		lineIndex++

		if len(resource.AlertRules) > 0 {
			err = addBoldCell(f, sheet, fmt.Sprintf("B%d", lineIndex), "Name")
			if err != nil {
				return err
			}

			err = addBoldCell(f, sheet, fmt.Sprintf("C%d", lineIndex), "Criteria")
			if err != nil {
				return err
			}

			lineIndex++

			for _, alertRule := range resource.AlertRules {
				err = writeCell(f, sheet, fmt.Sprintf("B%d", lineIndex), alertRule.Name)
				if err != nil {
					return err
				}

				for _, criterion := range alertRule.Criteria.AllOf {
					criterionOutput := fmt.Sprintf("%s %s %s %s", criterion.TimeAggregation,
						criterion.MetricName,
						criterion.Operator,
						fmt.Sprintf("%.2f", criterion.Threshold),
					)
					err = writeCell(f, sheet, fmt.Sprintf("C%d", lineIndex), criterionOutput)
					if err != nil {
						return err
					}

					lineIndex++
				}
			}
		}
	}

	return nil
}

func writeResourceBackups(f *excelize.File, vms map[string]azure.Resource) error {
	lineIndex := 1
	for _, vm := range vms {
		err := writeCell(f, "Backups", fmt.Sprintf("A%d", lineIndex), vm.Name)
		if err != nil {
			return err
		}

		if vm.BackupVault != nil {
			err = writeCell(f, "Backups", fmt.Sprintf("B%d", lineIndex), vm.BackupVault.Name)
			if err != nil {
				return err
			}

		} else {
			err = writeCell(f, "Backups", fmt.Sprintf("B%d", lineIndex), "Not backed up")
			if err != nil {
				return err
			}

		}
		lineIndex++
	}

	return nil
}

func writeRecommendations(f *excelize.File, recommendations map[string][]azure.AdvisorRecommendation) error {
	for category, categoryRecommendations := range recommendations {
		sheetTitle := fmt.Sprintf("%s Recs", category)
		f.NewSheet(sheetTitle)
		index := 1

		err := writeCell(f, sheetTitle, fmt.Sprintf("A%d", index), "Recommendation")
		if err != nil {
			return err
		}
		err = writeCell(f, sheetTitle, fmt.Sprintf("B%d", index), "Impact")
		if err != nil {
			return err
		}

		err = writeCell(f, sheetTitle, fmt.Sprintf("C%d", index), "Resource type")
		if err != nil {
			return err
		}

		err = writeCell(f, sheetTitle, fmt.Sprintf("D%d", index), "Affected resource")
		if err != nil {
			return err
		}

		err = writeCell(f, sheetTitle, fmt.Sprintf("E%d", index), "Resource group")
		if err != nil {
			return err
		}

		index++

		for i := 0; i < len(categoryRecommendations); i++ {
			err = writeCell(f, sheetTitle, fmt.Sprintf("A%d", index+1), categoryRecommendations[i].Description.Problem)
			if err != nil {
				return err
			}

			err = writeCell(f, sheetTitle, fmt.Sprintf("B%d", index+1), categoryRecommendations[i].Impact)
			if err != nil {
				return err
			}

			err = writeCell(f, sheetTitle, fmt.Sprintf("C%d", index+1), categoryRecommendations[i].ResourceType)
			if err != nil {
				return err
			}

			err = writeCell(f, sheetTitle, fmt.Sprintf("D%d", index+1), categoryRecommendations[i].AffectedResource)
			if err != nil {
				return err
			}

			err = writeCell(f, sheetTitle, fmt.Sprintf("E%d", index+1), categoryRecommendations[i].ResourceGroup)
			if err != nil {
				return err
			}

			index++
		}
	}

	return nil
}

func OutputExcelDocument(
	outputFilename string,
	vms map[string]azure.Resource,
	aksClusters map[string]azure.Resource,
	mySQLServers map[string]azure.Resource,
	flexibleMySQLServers map[string]azure.Resource,
	sqlServers map[string]azure.Resource,
	storageAccounts map[string]azure.Resource,
	webApps map[string]azure.Resource,
	recommendations map[string][]azure.AdvisorRecommendation,
) error {
	f := excelize.NewFile()

	if len(vms) > 0 {
		f.SetSheetName("Sheet1", "VM Alerts")
		err := writeResourceAlerts(f, "VM Alerts", vms)
		if err != nil {
			return err
		}

	} else {
		f.DeleteSheet("Sheet1")
	}

	if len(aksClusters) > 0 {
		f.NewSheet("AKS Cluster Alerts")
		err := writeResourceAlerts(f, "AKS Cluster Alerts", aksClusters)
		if err != nil {
			return err
		}

	}

	if len(mySQLServers) > 0 {
		f.NewSheet("MySQL Server Alerts")
		err := writeResourceAlerts(f, "MySQL Server Alerts", mySQLServers)
		if err != nil {
			return err
		}

	}

	if len(flexibleMySQLServers) > 0 {
		f.NewSheet("Flexible MySQL Server Alerts")
		err := writeResourceAlerts(f, "Flexible MySQL Server Alerts", flexibleMySQLServers)
		if err != nil {
			return err
		}

	}

	if len(sqlServers) > 0 {
		f.NewSheet("SQL Server Alerts")
		err := writeResourceAlerts(f, "SQL Server Alerts", sqlServers)
		if err != nil {
			return err
		}

	}

	if len(storageAccounts) > 0 {
		f.NewSheet("Storage Account Alerts")
		err := writeResourceAlerts(f, "Storage Account Alerts", storageAccounts)
		if err != nil {
			return err
		}

	}

	if len(webApps) > 0 {
		f.NewSheet("Web App Alerts")
		err := writeResourceAlerts(f, "Web App Alerts", webApps)
		if err != nil {
			return err
		}

	}

	// Backups
	f.NewSheet("Backups")
	err := writeResourceBackups(f, vms)
	if err != nil {
		return err
	}

	err = writeRecommendations(f, recommendations)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s.xlsx", outputFilename)
	if err := f.SaveAs(filename); err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("Saved excel file to %s", filename))

	return nil
}
