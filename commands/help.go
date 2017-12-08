package commands

import (
	"github.com/adolphlxm/atc-tool/utils"
	"os"
	"fmt"
	"html/template"
)

var usage = `Atc-tool is a tool for managing ATC framework.

Usage:

	atc-tool command [arguments]

The commands are:
`
var usageTemplate = `	{{.Name}}	{{.Use | printf "%s"}}
`

var usageHelpTemplate = `
Use "atc-tool help [command]" for more information about a command.
`

var helpTemplate = `
usage: {{.Usage | printf "%s"}}

{{.Use | printf "%s"}}
{{.Options | printf "%s"}}
`


func Usage() {
	fmt.Println(usage)
	for _, cmd := range AdapterCommands {
		data := template.FuncMap{"Name":cmd.Name(), "Use":cmd.Use}
		utils.Tmpl(usageTemplate, data,os.Stderr)
	}
	fmt.Println(usageHelpTemplate)
}

func Help(args []string) {
	if len(args) == 0 {
		Usage()
		return
	}

	if len(args) != 1 {
		fmt.Println("atc-tool: Too many arguments")
		return
	}

	arg := args[0]
	for _, c := range AdapterCommands {
		if c.Name() == arg{
			utils.Tmpl(helpTemplate, c,os.Stderr)
			return
		}
	}

	fmt.Println("atc-tool: Unknown help topic")
}