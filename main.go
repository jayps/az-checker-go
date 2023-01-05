package main

import (
	"fmt"
	"github.com/jayps/azure-checker-go/azure"
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
	fmt.Println("Enter a filename for the output document:")

	var filename string
	_, err := fmt.Scanln(&filename)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return filename
}

func main() {

	subscriptionIds := getSubscriptionIds()
	filename := getFilename()

	for i := 0; i < len(subscriptionIds); i++ {
		subscriptionId := subscriptionIds[i]
		azure.SetSubscription(subscriptionId)
		fmt.Println(fmt.Sprintf("Set subscription ID to %s...", subscriptionId))

		// Fetch resources
		vms := azure.FetchVMs()
		aksClusters := azure.FetchAKSClusters()
		//mySQLServers := azure.FetchMySQLServers()
		//flexibleMySQLServers := azure.FetchFlexibleMySQLServers()
		//sqlServers := azure.FetchSQLServers()
		//storageAccounts := azure.FetchStorageAccounts()
		//webApps := azure.FetchWebApps()
		alertRules := azure.FetchAlertRules()

		// Assign alert rules
		azure.AssignAlertRulesToResources(alertRules, vms)
		azure.AssignAlertRulesToResources(alertRules, aksClusters)
		//azure.AssignAlertRulesToResources(alertRules, mySQLServers)
		//azure.AssignAlertRulesToResources(alertRules, flexibleMySQLServers)
		//azure.AssignAlertRulesToResources(alertRules, sqlServers)
		//azure.AssignAlertRulesToResources(alertRules, storageAccounts)
		//azure.AssignAlertRulesToResources(alertRules, webApps)

		//azure.FetchVMBackups(vms)
		//recommendations := azure.FetchAdvisorRecommendations()

		now := time.Now()
		outputFilename := fmt.Sprintf("%s-%s-%d-%d-%d", filename, subscriptionId, now.Year(), now.Month(), now.Day())

		pdf.GeneratePDF()

		//fmt.Println(fmt.Sprintf("Saving checks for subscription ID %s to %s...", subscriptionId, outputFilename))
		//excel.OutputExcelDocument(
		//	outputFilename,
		//	vms,
		//	aksClusters,
		//	mySQLServers,
		//	flexibleMySQLServers,
		//	sqlServers,
		//	storageAccounts,
		//	webApps,
		//	recommendations,
		//)
	}

	fmt.Println("All done, press Enter to exit.")
	fmt.Scanln()
}
