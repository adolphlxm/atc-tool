package thrift

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/adolphlxm/atc-tool/commands"
	"github.com/adolphlxm/atc-tool/utils"
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
	b, err := exeCmd("thrift", args...)
	if err != nil {
		return 1
	}

	// Replace thrift package
	gen := args[3]
	i := strings.LastIndex(gen, "/")
	thriftGenName := gen[i+1 : i+6]
	helper := utils.ReplaceHelper{
		Root:    "./gen-go/" + thriftGenName,
		OldText: "git.apache.org/thrift.git/lib/go/thrif",
		NewText: "github.com/adolphlxm/atc/rpc/thrift/lib/thrif",
	}
	err = helper.DoWrok()
	if err != nil {
		commands.Logger.Error("Replace thrift package err:%v", err)
		return 1
	}

	fmt.Printf("%s", b)

	return 2
}

func exeCmd(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		commands.Logger.Error("StdoutPipe: %s", err.Error())
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		commands.Logger.Error("StderrPipe: %s ", err.Error())
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	bytesErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		commands.Logger.Error("ReadAll stderr: %s", err.Error())
		return bytesErr, err
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		commands.Logger.Error("ReadAll stdout: %s", err.Error())
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		fmt.Printf("%s", bytesErr)
		return nil, err
	}
	return bytes, nil
}
