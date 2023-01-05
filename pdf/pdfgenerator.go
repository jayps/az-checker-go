package pdf

import (
	"fmt"
	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"log"
	"strings"
)

func GeneratePDF() {
	pdfGenerator, err := wkhtml.NewPDFGenerator()
	if err != nil {
		return
	}
	htmlStr := `<html><body><h1 style="color:red;">This is an html from pdf to test color<h1><img src="http://api.qrserver.com/v1/create-qr-code/?data=HelloWorld" alt="img" height="42" width="42"></img></body></html>`

	pdfGenerator.AddPage(wkhtml.NewPageReader(strings.NewReader(htmlStr)))

	// Create PDF document in internal buffer
	err = pdfGenerator.Create()
	if err != nil {
		log.Fatal(err)
	}

	//Your Pdf Name
	err = pdfGenerator.WriteFile("./Your_pdfname.pdf")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")

}
