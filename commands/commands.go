package commands

import (
	"flag"
	"strings"

	"github.com/adolphlxm/atc/logs"
)

var ErrUseError = "Use atc-tool -help for a list"

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
	Usage string
	Use     string
	Options      string

	Flag flag.FlagSet
}

// Name returns the command's name: the first word in the Usage line.
func (c *Command) Name() string {
	name := c.Usage
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func init(){
	// Initialize log
	Logger = logs.NewLogger(10000)
	Logger.SetHandler("stdout", ``)
	Logger.SetLevel(logs.LevelDebug)
}