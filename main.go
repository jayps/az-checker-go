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

func main() {
	subscriptionIds := getSubscriptionIds()
	clientName := getFilename()

	for i := 0; i < len(subscriptionIds); i++ {
		subscriptionId := subscriptionIds[i]
		err := azure.SetSubscription(subscriptionId)
		if err != nil {
			log.Fatalln("Could not set subscription: ", err.Error())
		}
		fmt.Println(fmt.Sprintf("Set subscription ID to %s...", subscriptionId))

		// Fetch resources
		vms, err := azure.FetchVMs()
		if err != nil {
			log.Fatalln("Could not fetch VMs: ", err.Error())
		}

		aksClusters, err := azure.FetchAKSClusters()
		if err != nil {
			log.Fatalln("Could not fetch AKS clusters: ", err.Error())
		}

		mySQLServers, err := azure.FetchMySQLServers()
		if err != nil {
			log.Fatalln("Could not fetch MySQL servers: ", err.Error())
		}

		flexibleMySQLServers, err := azure.FetchFlexibleMySQLServers()
		if err != nil {
			log.Fatalln("Could not fetch flexible MySQL servers: ", err.Error())
		}

		sqlServers, err := azure.FetchSQLServers()
		if err != nil {
			log.Fatalln("Could not fetch SQL servers: ", err.Error())
		}

		storageAccounts, err := azure.FetchStorageAccounts()
		if err != nil {
			log.Fatalln("Could not fetch storage accounts: ", err.Error())
		}

		webApps, err := azure.FetchWebApps()
		if err != nil {
			log.Fatalln("Could not fetch web apps: ", err.Error())
		}

		alertRules, err := azure.FetchAlertRules()
		if err != nil {
			log.Fatalln("Could not fetch alert rules: ", err.Error())
		}

		// Assign alert rules
		azure.AssignAlertRulesToResources(alertRules, vms)
		azure.AssignAlertRulesToResources(alertRules, aksClusters)
		azure.AssignAlertRulesToResources(alertRules, mySQLServers)
		azure.AssignAlertRulesToResources(alertRules, flexibleMySQLServers)
		azure.AssignAlertRulesToResources(alertRules, sqlServers)
		azure.AssignAlertRulesToResources(alertRules, storageAccounts)
		azure.AssignAlertRulesToResources(alertRules, webApps)

		err = azure.FetchVMBackups(vms)
		if err != nil {
			log.Fatalln("Could not fetch VM backups: ", err.Error())
		}

		recommendations, err := azure.FetchAdvisorRecommendations()
		if err != nil {
			log.Fatalln("Could not fetch advisor recommendations: ", err.Error())
		}

		for key, vm := range vms {
			err = azure.AssessPatches(&vm)
			if err != nil {
				log.Fatalln(fmt.Sprintf("Could not assess patches for VM %s: ", vm.Name), err.Error())
			}
			fmt.Println(fmt.Sprintf("%d cricial patches, %d other patches for %s", vm.PatchAssessmentResult.CriticalAndSecurityPatchCount, vm.PatchAssessmentResult.OtherPatchCount, vm.Name))
			vms[key] = vm // range gives us a copy of the vm we are working with, so we reassign it back to the map.
		}

		for _, vm := range vms {
			fmt.Println(fmt.Sprintf("POST ASSESSMENT: %d cricial patches, %d other patches for %s", vm.PatchAssessmentResult.CriticalAndSecurityPatchCount, vm.PatchAssessmentResult.OtherPatchCount, vm.Name))
		}

		now := time.Now()
		outputFilename := fmt.Sprintf("%s-%s-%d-%d-%d", clientName, subscriptionId, now.Year(), now.Month(), now.Day())

		g := pdf.NewGenerator()
		g.ClientName = clientName
		g.SubscriptionId = subscriptionId
		g.OutputFilename = outputFilename
		g.VirtualMachines = vms
		g.AzureKubernetesServices = aksClusters
		g.MySQLServers = mySQLServers
		g.FlexibleMySQLServers = flexibleMySQLServers
		g.SqlServers = sqlServers
		g.StorageAccounts = storageAccounts
		g.WebApps = webApps
		g.Recommendations = recommendations
		err = g.GeneratePDF()
		if err != nil {
			log.Fatalln(err)
		}
		err = excel.OutputExcelDocument(
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
		if err != nil {
			log.Fatalln("Could not generate excel file: ", err.Error())
		}
	}

	fmt.Println("All done.")
}
