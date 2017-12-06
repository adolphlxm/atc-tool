package utils

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os/exec"
	"os"
	"runtime"
	"path/filepath"
	"strings"
	"errors"
)

// WriteToFile creates a file and writes content to it
func WriteToFile(filename, content string) {
	f, err := os.Create(filename)
	if err != nil {

	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {

	}
}


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
		return nil, fmt.Errorf("%s:err:%s", bytesErr, err.Error())
	}
	return bytes, nil
}

func CheckEnv(appname string) (packpath string, err error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" && strings.Compare(runtime.Version(), "go1.8") >= 0 {
		gopath = DefaultGOPATH()
	}

	currpath, _ := os.Getwd()
	currpath = filepath.Join(currpath, appname)
	gopathStr := filepath.SplitList(gopath)

	for _, gpath := range gopathStr {
		gsrcpath := filepath.Join(gpath, "src")
		if strings.HasPrefix(strings.ToLower(currpath), strings.ToLower(gsrcpath)) {
			packpath = strings.Replace(currpath[len(gsrcpath)+1:], string(filepath.Separator), "/", -1)
			return
		}
	}

	return packpath, errors.New("You current workdir is not inside $GOPATH/src.")
}

func DefaultGOPATH() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		return filepath.Join(home, "go")
	}
	return ""
}