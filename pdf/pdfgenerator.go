package pdf

import (
	"fmt"
	"strings"
	"time"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/jayps/azure-checker-go/azure"
)

type Generator struct {
	Head                       string `default:"test"`
	ClientName                 string `default:"Client"`
	SubscriptionId             string
	OutputFilename             string
	VirtualMachines            map[string]azure.Resource
	VirtualMachinesDeallocated map[string]azure.Resource
	AzureKubernetesServices    map[string]azure.Resource
	MySQLServers               map[string]azure.Resource
	FlexibleMySQLServers       map[string]azure.Resource
	SqlServers                 map[string]azure.Resource
	StorageAccounts            map[string]azure.Resource
	WebApps                    map[string]azure.Resource
	Recommendations            map[string][]azure.AdvisorRecommendation
}

func NewGenerator() Generator {
	result := Generator{
		Head: `
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Montserrat:wght@300;400;700&display=swap" rel="stylesheet">
<style>
body {
font-family: 'Montserrat', sans-serif;
color: #666;
}

br {
	margin-bottom: 0;
}

p {
font-size: 1em;
line-height: 1.2em;
font-weight: 400; 
}

h1 {
	font-size: 3em;
	font-weight: 700;
}

h2 {
	font-size: 2em;
	font-weight: 700;
}

h3 {
	font-size: 1.2em;
	font-weight: 700;
	margin-bottom: 0;
	line-height: 2em;
}

strong {
	font-weight: 700;
}

small {
	font-size: 0.75em;
}

.bg-grey {
	background-color: #eee;
}

.p-1 {
	padding: 24px;
}

.mb-1 {
	margin-bottom: 24px;
}

.page-break-before {
	page-break-before: always;
}

.page-break-avoid {
	page-break-inside: avoid; /* TODO: Figure out why this makes random-ish page breaks sometimes. */
}

.warn {
	color: orange;
}

.danger {
	color: red;
}

.bg-warn {
	background-color: orange;
}

.bg-danger {
	background-color: red;
}

</style>
`,
	}

	return result
}

func (g Generator) GenerateAlertRulesSection(title string, resources map[string]azure.Resource) string {
	if len(resources) == 0 {
		return ""
	}

	output := "<div class='page-break-before'>"
	output += fmt.Sprintf("<h2>Monitoring: %s</h2>", title)
	if len(resources) == 0 {
		output += "No resources of this type."
		return output
	}
	for _, resource := range resources {
		output += "<div class='page-break-avoid'>"
		output += fmt.Sprintf("<h3>%s</h3>", resource.Name)
		if len(resource.AlertRules) == 0 {
			output += "<span class='danger'>No alert rules are configured for this resource.</span><br />"
			output += "<strong>Action to be performed: </strong>If this alert is used in production, create resource alert rules. We do not monitor non-production resources."
		} else {
			output += fmt.Sprintf("This resource has %d alert rules configured:<br /><br />", len(resource.AlertRules))
			output += "<div class='bg-grey p-1'>"
			for _, rule := range resource.AlertRules {
				output += "<small>"
				output += fmt.Sprintf("<strong>Rule: %s: </strong><br />", rule.Name)
				for _, criterion := range rule.Criteria.AllOf {
					output += fmt.Sprintf("<strong>Criteria: </strong>%s %s %s %s<br /><br />", criterion.TimeAggregation,
						criterion.MetricName,
						criterion.Operator,
						fmt.Sprintf("%.2f", criterion.Threshold),
					)
				}
				output += "</small>"
			}
			output += "<strong>Action to be performed:</strong> Review alert rules and confirm that they are appropriate for this resource."
			output += "</div>" // background grey
		}
		output += "</div>" // page break avoid
	}
	output += "</div>" // page break before

	return output
}

func (g Generator) GenerateBackupsSection() string {
	if len(g.VirtualMachines) == 0 {
		return ""
	}
	output := "<div class='page-break-before'>"
	output += fmt.Sprintf("<h2>Virtual Machine Backups</h2>")
	for _, vm := range g.VirtualMachines {
		output += fmt.Sprintf("<h3>%s</h3>", vm.Name)
		if vm.BackupVault != nil {
			output += fmt.Sprintf("This virtual machine is backed up to %s.<br />", vm.BackupVault.Name)
			output += fmt.Sprintf("<strong>Action to be performed:</strong> None")
		} else {
			output += fmt.Sprintf("<span class='danger'>This virtual machine is not backed up.</span><br />")
			output += fmt.Sprintf("<strong>Action to be performed:</strong> If this is a production machine, consider setting up backups using Azure Backup Vault. If an alternative backup solution is being used, this recommendation can be ignored.")
		}
	}
	output += "</div>" // page break before

	return output
}

func (g Generator) GenerateDeallocatedVMsSection() string {
	if len(g.VirtualMachinesDeallocated) == 0 {
		return ""
	}
	output := "<div class='page-break-before'>"
	output += fmt.Sprintf("<h2>Deallocated Virtual Machines</h2>")
	for _, vm := range g.VirtualMachines {
		output += fmt.Sprintf("<h3>%s</h3>", vm.Name)
	}
	output += "</div>" // page break before

	return output
}

func (g Generator) GeneratePatchesSection() string {
	if len(g.VirtualMachines) == 0 {
		return ""
	}
	output := "<div class='page-break-before'>"
	output += fmt.Sprintf("<h2>Virtual Machine Patches</h2>")
	for _, vm := range g.VirtualMachines {
		output += fmt.Sprintf("<h3>%s</h3>", vm.Name)
		output += fmt.Sprintf("<span class='mb-1'>%d patches available.", len(vm.PatchAssessmentResult.AvailablePatches))
		for _, patch := range vm.PatchAssessmentResult.AvailablePatches {
			output += "<div class='mb-1 page-break-avoid bg-grey p-1'>"
			output += fmt.Sprintf("<strong>Patch Name: %s</strong><br />", patch.Name)
			output += fmt.Sprintf("<strong>Patch ID: %s</strong><br />", patch.PatchId)
			output += fmt.Sprintf("<strong>KB ID: %s</strong><br />", patch.KbId)
			output += fmt.Sprintf("<strong>Version: %s</strong><br />", patch.Version)
			output += fmt.Sprintf("<strong>Reboot: %s</strong><br />", patch.RebootBehavior)
			output += "</div>" // page break avoid
		}
	}
	output += "</div>" // page break before

	return output
}

func (g Generator) GenerateRecommendationsSections() string {
	output := "<div class='page-break-before'>"
	for category, categoryRecommendations := range g.Recommendations {
		output += fmt.Sprintf("<h2>Advisory Recommendations: %s</h2>", category)

		if len(categoryRecommendations) == 0 {
			output += "No recommendations in this category. Looking good!"
			return output
		}

		for _, rec := range categoryRecommendations {
			color := "" // don't change the color if we're looking at a low recommendation
			if rec.Impact == "Medium" {
				color = "warn"
			}
			if rec.Impact == "High" {
				color = "danger"
			}
			output += "<div class='mb-1 page-break-avoid bg-grey p-1'>"
			output += fmt.Sprintf("<strong>%s</strong><br />", rec.Description.Problem)
			output += fmt.Sprintf("<span>Impact:</span> <span class='%s'>%s</span><br />", color, rec.Impact)
			output += fmt.Sprintf("<span>Resource Type:</span> %s<br />", rec.ResourceType)
			output += fmt.Sprintf("<span>Affected Resource:</span> %s<br />", rec.AffectedResource)
			output += fmt.Sprintf("<span>Resource Group:</span> %s", rec.ResourceGroup)
			output += "</div>" // page break avoid
		}
	}
	output += "</div>" // page break before

	return output
}

func (g Generator) GeneratePDF() error {
	pdfGenerator, err := wkhtml.NewPDFGenerator()
	if err != nil {
		return err
	}
	htmlStr := `<html>
<head>
{headContent}
</head>
<body>
<div style="margin-top: 320px; text-align: center;">
<img height="150px; margin-top: " src="data:image/png;base64,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">
<h1>
	{title}
</h1>
<h2>
{clientName}
</h2>
<h3>
Subscription ID: {subscriptionId}
</h3>
<p>
Document Date: {date}
</p>
</div>
{vmAlerts}
{aksAlerts}
{mySQLServers}
{flexibleMySQLServers}
{sqlServers}
{storageAccounts}
{webApps}
{backups}
{deallocatedVMs}
{recommendations}
{patches}
</body>
</html>`

	now := time.Now()
	documentReplacer := strings.NewReplacer(
		"{headContent}", g.Head,
		"{title}", "Tangent Solutions Managed Services Report",
		"{clientName}", g.ClientName,
		"{subscriptionId}", g.SubscriptionId,
		"{date}", fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day()),
		"{vmAlerts}", g.GenerateAlertRulesSection("Virtual Machines", g.VirtualMachines),
		"{aksAlerts}", g.GenerateAlertRulesSection("Azure Kubernetes Services", g.AzureKubernetesServices),
		"{mySQLServers}", g.GenerateAlertRulesSection("MySQL Servers", g.MySQLServers),
		"{flexibleMySQLServers}", g.GenerateAlertRulesSection("Flexible MySQL Servers", g.FlexibleMySQLServers),
		"{sqlServers}", g.GenerateAlertRulesSection("SQL Servers", g.SqlServers),
		"{storageAccounts}", g.GenerateAlertRulesSection("Storage Accounts", g.StorageAccounts),
		"{webApps}", g.GenerateAlertRulesSection("Web Apps", g.WebApps),
		"{backups}", g.GenerateBackupsSection(),
		"{patches}", g.GeneratePatchesSection(),
		"{recommendations}", g.GenerateRecommendationsSections(),
		"{deallocatedVMs}", g.GenerateDeallocatedVMsSection(),
	)
	populatedHtml := documentReplacer.Replace(htmlStr)

	var margin uint = 16
	pdfGenerator.MarginRight.Set(margin)
	pdfGenerator.MarginLeft.Set(margin)
	pdfGenerator.MarginTop.Set(margin)
	pdfGenerator.MarginBottom.Set(margin)
	pdfGenerator.AddPage(wkhtml.NewPageReader(strings.NewReader(populatedHtml)))

	// Create PDF document in internal buffer
	err = pdfGenerator.Create()
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("./%s.pdf", g.OutputFilename)
	err = pdfGenerator.WriteFile(filename)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Saved PDF report to %s", filename))

	return nil
}
