package commands

import (
	"flag"
	"github.com/adolphlxm/atc/logs"
)

var AdapterCommands = []*Command{}
var Logger *logs.AtcLogger

func Register(cmd *Command) {
	if cmd == nil {
		Logger.Fatal("ATC command: Register command is nil")
	}
	AdapterCommands = append(AdapterCommands, cmd)
}

type Command struct {
	Run       func(cmd *Command, args []string) int
	UsageLine string
	Short     string
	Long      string

	Flag flag.FlagSet
}

func init(){
	// Initialize log
	Logger = logs.NewLogger(10000)
	Logger.SetHandler("stdout", ``)
	Logger.SetLevel(logs.LevelDebug)
}