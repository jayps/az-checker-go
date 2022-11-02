package main

import (
	"azurechecker/azure"
	"azurechecker/excel"
	"fmt"
	"log"
	"time"
)

func getSubscriptionId() string {
	fmt.Println("Enter the subscription ID you are checking:")

	var subscriptionId string
	_, err := fmt.Scanln(&subscriptionId)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return subscriptionId
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

	subscriptionId := getSubscriptionId()
	filename := getFilename()

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

	now := time.Now()
	excel.OutputExcelDocument(
		fmt.Sprintf("%s-%s-%d-%d-%d", filename, subscriptionId, now.Year(), now.Month(), now.Day()),
		vms,
		aksClusters,
		mySQLServers,
		flexibleMySQLServers,
		sqlServers,
		storageAccounts,
		webApps,
		recommendations,
	)

	fmt.Println("All done, press Enter to exit.")
	fmt.Scanln()
}
