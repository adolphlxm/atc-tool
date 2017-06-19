package utils

import (
	"html/template"
	"io"
	"os/exec"
	"fmt"
	"io/ioutil"
)

// 模板替换
func Tmpl(text string, data interface{}, wr io.Writer) error {

	t := template.New("Usage")
	template.Must(t.Parse(text))

	return t.Execute(wr, data)
}

// 命令行
func ExeCmd(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("StdoutPipe: %s", err.Error())
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("StderrPipe: %s ", err.Error())
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	bytesErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		return bytesErr, fmt.Errorf("ReadAll stderr: %s", err.Error())
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("ReadAll stdout: %s", err.Error())
	}
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("%s:err:%s", bytesErr,err.Error())
	}
	return bytes, nil
}
