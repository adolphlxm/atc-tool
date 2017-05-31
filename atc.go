package main

import (
	"flag"
	"github.com/adolphlxm/atc-tool/commands"
	"os"

	_ "github.com/adolphlxm/atc-tool/commands/orm"
	"time"
)

const version = "0.1.0"

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
		if c.Run != nil {
			c.Flag.Parse(args[1:])
			args = c.Flag.Args()
			code := c.Run(c, args)

			time.Sleep(1 * time.Millisecond)
			os.Exit(code)
		}
	}

	//fmt.Println(currentpath)
}
