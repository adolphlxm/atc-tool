package main

import (
	"flag"
	"os"
	"time"

	_ "github.com/adolphlxm/atc-tool/commands/orm"
	_ "github.com/adolphlxm/atc-tool/commands/thrift"
	"github.com/adolphlxm/atc-tool/commands"
	"fmt"
	"runtime"
)

const VERSION = "0.2.0"

func init() {
	version := flag.Bool("v", false, "Use -v <current version>")
	flag.Parse()
	// Show version
	if *version {
		fmt.Println("atc-tool version",VERSION, runtime.GOOS +"/" + runtime.GOARCH)
		os.Exit(0)
	}
}

func main() {
	//currentpath, _ := os.Getwd()
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		os.Exit(2)
		return
	}

	// Help
	if args[0] == "help" {

	}

	for _, c := range commands.AdapterCommands {
		if c.Name() == args[0] && c.Run != nil {
			//fmt.Println(args[1:])
			//c.Flag.Parse(args[1:])
			//args = c.Flag.Args()
			code := c.Run(c, args[1:])

			time.Sleep(1 * time.Millisecond)
			os.Exit(code)
		}
	}
}
