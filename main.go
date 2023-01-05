package main

import (
	"fmt"
	"github.com/jayps/azure-checker-go/azure"
	"github.com/jayps/azure-checker-go/excel"
	"github.com/jayps/azure-checker-go/pdf"
	"log"
	"strings"
	"time"
)

func getSubscriptionIds() []string {
	fmt.Println("Enter a comma separated list of subscription IDs you are checking:")

	var subscriptionIdsInput string
	_, err := fmt.Scanln(&subscriptionIdsInput)
	if err != nil {
		log.Fatalln(err.Error())
	}

	subscriptionIds := strings.Split(subscriptionIdsInput, ",")

	return subscriptionIds
}

func getFilename() string {
	fmt.Println("Enter the name of the client:")

	var filename string
	_, err := fmt.Scanln(&filename)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return filename
}

func getYesNoChoice(question string, defaultAnswer bool) bool {
	defaultAnswerString := "Y/n"
	if !defaultAnswer {
		defaultAnswerString = "y/N"
	}

	fmt.Println(fmt.Sprintf("%s %s", question, defaultAnswerString))

	var answer string
	_, err := fmt.Scanln(&answer)
	if err != nil {
		log.Fatalln(err.Error())
	}

	answer = strings.ToLower(answer)

	if answer == "" {
		return defaultAnswer
	}

	return answer == "y"
}

func main() {

	subscriptionIds := getSubscriptionIds()
	clientName := getFilename()
	generatePdfFile := getYesNoChoice("Do you want to generate a PDF report?", true)
	generateExcelFile := getYesNoChoice("Do you want to generate an Excel report?", true)

	if !generatePdfFile && !generateExcelFile {
		fmt.Println("You must select at least one output file type.")
		return
	}

	for i := 0; i < len(subscriptionIds); i++ {
		subscriptionId := subscriptionIds[i]
		azure.SetSubscription(subscriptionId)
		fmt.Println(fmt.Sprintf("Set subscription ID to %s...", subscriptionId))

		// Fetch resources
		vms := azure.FetchVMs()
		aksClusters := azure.FetchAKSClusters()
		mySQLServers := azure.FetchMySQLServers()
		flexibleMySQLServers := azure.FetchFlexibleMySQLServers()
		sqlServers := azure.FetchSQLServers()
		storageAccounts := azure.FetchStorageAccounts()
		webApps := azure.FetchWebApps()
		alertRules := azure.FetchAlertRules()

		// Assign alert rules
		azure.AssignAlertRulesToResources(alertRules, vms)
		azure.AssignAlertRulesToResources(alertRules, aksClusters)
		azure.AssignAlertRulesToResources(alertRules, mySQLServers)
		azure.AssignAlertRulesToResources(alertRules, flexibleMySQLServers)
		azure.AssignAlertRulesToResources(alertRules, sqlServers)
		azure.AssignAlertRulesToResources(alertRules, storageAccounts)
		azure.AssignAlertRulesToResources(alertRules, webApps)

		azure.FetchVMBackups(vms)
		recommendations := azure.FetchAdvisorRecommendations()

		for _, vm := range vms {
			azure.AssessPatches(&vm)
			fmt.Println(fmt.Sprintf("%d cricial patches, %d other patches for %s", vm.PatchAssessmentResult.CriticalAndSecurityPatchCount, vm.PatchAssessmentResult.OtherPatchCount, vm.Name))
			break
		}

		now := time.Now()
		outputFilename := fmt.Sprintf("%s-%s-%d-%d-%d", clientName, subscriptionId, now.Year(), now.Month(), now.Day())

		if generatePdfFile {
			g := pdf.NewGenerator()
			g.ClientName = clientName
			g.OutputFilename = outputFilename
			g.VirtualMachines = vms
			g.AzureKubernetesServices = aksClusters
			g.MySQLServers = mySQLServers
			g.FlexibleMySQLServers = flexibleMySQLServers
			g.SqlServers = sqlServers
			g.StorageAccounts = storageAccounts
			g.WebApps = webApps
			g.Recommendations = recommendations
			err := g.GeneratePDF()
			if err != nil {
				log.Fatal(err)
			}
		}

		if generateExcelFile {
			excel.OutputExcelDocument(
				outputFilename,
				vms,
				aksClusters,
				mySQLServers,
				flexibleMySQLServers,
				sqlServers,
				storageAccounts,
				webApps,
				recommendations,
			)
		}
	}

	fmt.Println("All done.")
}
