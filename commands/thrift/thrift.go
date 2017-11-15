package thrift

import (
	"fmt"
	"strings"

	"github.com/adolphlxm/atc-tool/commands"
	"github.com/adolphlxm/atc-tool/utils"
	"github.com/adolphlxm/atc/logs"
)

var CmdThrift = &commands.Command{
	Usage: "thrift [options] file",
	Use:   "Thrift files are generated",
	Options: `
	Please use 'thrift --help' to obtain.
`,
	Run: Run,
}

func init() {
	commands.Register(CmdThrift)
}

func Run(cmd *commands.Command, args []string) int {
	b, err := utils.ExeCmd("thrift", args...)
	if err != nil {
		logs.Error("%s", err.Error())
		return 1
	}

	// Replace thrift package
	gen := args[3]
	i1 := strings.LastIndex(gen, "/")
	i2 := strings.LastIndex(gen, ".")
	thriftGenName := gen[i1+1:i2]

	helper := utils.ReplaceHelper{
		Root:    "./gen-go/" + thriftGenName,
		OldText: "git.apache.org/thrift.git/lib/go/thrif",
		NewText: "github.com/adolphlxm/atc/rpc/thrift/lib/thrif",
	}
	err = helper.DoWrok()
	if err != nil {
		logs.Error("Replace thrift package err:%v", err)
		return 1
	}

	fmt.Printf("%s", b)

	return 2
}

