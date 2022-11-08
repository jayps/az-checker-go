# Azure Checker

This tool will check some Azure resources for basic setup options and output a file containing the result. 
## Usage
You can run the tool with one of two methods:
`go run main.go` or `./azure-checker-go.exe`. The latter there will be different depending on your OS.

The tool will prompt you for subscription IDs - you can enter multiple if you want. They should be comma separated, with no spaces. 
You will also be prompted for a filename - this gets used as part of the filename for the output documentation. This just makes the files easier to identify if you're running the tool for multiple clients.  
Once you have satisfied the prompts, the tool will run through Azure resources and output an Excel file. The tool will tell you what the output filename is.

## Requesting Additional Features
If you want the tool to do more stuff, contact me or create an issue on the repo.

## Compiling from source
You can compile the tool yourself if you like. You'll need to have [Golang](https://go.dev/doc/install) installed on your machine. From there, you can just run `go run main.go`.