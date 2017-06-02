package orm

import (
	"bytes"
	"html/template"
	"os"
	"strings"

	"github.com/adolphlxm/atc"
	"github.com/adolphlxm/atc-tool/commands"
	"github.com/adolphlxm/atc-tool/utils"
)

var (
	c string
	s string
	j bool
)

var CmdOrm = &commands.Command{
	UsageLine: "orm dbName table [-c=config path] [-s=Generate the file path] [-j= true|false]",
	Short:     "reverse a database table to codes",
	Long: `
according database's tables and columns to generate codes for Go.
    dbName            Database aliasname, now supported four: mysql mymysql sqlite3 postgres
    table	      Database table name
    -c		      Database config path
    -s                Generated one go file for every table
    -j                Support Json tag
`,
	Run: Run,
}

var dbStruct = `package {{.pkgName}}

type {{.dbName}} struct {
	{{.dbFiles}}
}
`

func init() {
	CmdOrm.Flag.StringVar(&c, "c", "", "Get database configuration file")
	CmdOrm.Flag.StringVar(&s, "s", "", "Generated one go file for every table")
	CmdOrm.Flag.BoolVar(&j, "j", false, "Support Json tag true or false.")
	commands.Register(CmdOrm)
}

// Initialize orm.
func initOrm() error {
	// Initialize config
	err := atc.ParseConfig(c)
	if err != nil {
		commands.Logger.Error("%v", err.Error())
		return err
	}

	atc.Aconfig.Debug = false

	atc.RunOrm()

	return err
}

func Run(cmd *commands.Command, args []string) int {
	if len(args) < 2 {

	}
	cmd.Flag.Parse(args[2:])
	if err := initOrm(); err != nil {
		return 1
	}

	if err := reverse(args[0], args[1], j); err != nil {

	}

	return 2
}

//生成表结构
func reverse(aliasName string, tableName string, isJson bool) error {
	var str string

	engine := atc.Orm.Use(aliasName)
	results, _ := engine.Query("show columns from " + tableName)

	for _, res := range results {
		//fmt.Println("Field->" + string(res["Field"]))
		//fmt.Println("Key->" + string(res["Key"]))
		//fmt.Println("Extra->" + string(res["Extra"]))
		//fmt.Println("Type->" + string(res["Type"]))
		//fmt.Println("NUll->" + string(res["Null"]))
		//fmt.Println("Default->" + string(res["Default"]))
		//fmt.Println("========================")

		//field
		fieldUpper := ""
		field := string(res["Field"])
		isUpper := false
		for k, v := range res["Field"] {
			if k == 0 {
				// Capitalize the first letter.
				fieldUpper = strings.ToUpper(string(v))
			} else {
				// 解决Xorm大写字段会按照_分割问题
				if string(v) == strings.ToUpper(string(v)) && string(v) != "_" {
					isUpper = true
				}
				fieldUpper += string(v)
			}
		}

		dbType := ""
		xormType := ""
		if string(res["Key"]) == "PRI" { //Primary Key
			autoincr := ""
			if string(res["Extra"]) == "auto_increment" {
				autoincr = "autoincr"
			}
			dbType = "int64"
			xormType = "pk" + " " + autoincr

		} else {
			trimUnsigned := bytes.TrimSuffix(res["Type"], []byte("unsigned"))
			xormType = string(trimUnsigned)
			isUnsigned := bytes.HasSuffix(res["Type"], []byte("unsigned"))

			switch {
			case bytes.HasPrefix(res["Type"], []byte("int")):
				dbType = "int64"
				if isUnsigned {
					dbType = "uint64"
				}
			case bytes.HasPrefix(res["Type"], []byte("smallint")):
				dbType = "int16"
				if isUnsigned {
					dbType = "uint16"
				}
			case bytes.HasPrefix(res["Type"], []byte("tinyint")):
				dbType = "int8"
				if isUnsigned {
					dbType = "byte"
				}
			case bytes.HasPrefix(res["Type"], []byte("bigint")):
				dbType = "int64"
				if isUnsigned {
					dbType = "uint64"
				}
			case bytes.HasPrefix(res["Type"], []byte("mediumint")):
				dbType = "int32"
				if isUnsigned {
					dbType = "uint32"
				}
			case bytes.HasPrefix(res["Type"], []byte("float")):
				dbType = "float64"
			case bytes.HasPrefix(res["Type"], []byte("double")):
				dbType = "float64"
			case bytes.HasPrefix(res["Type"], []byte("decimal")):
				dbType = "string"
			case bytes.HasPrefix(res["Type"], []byte("date")):
				dbType = "string"
			case bytes.HasPrefix(res["Type"], []byte("datetime")):
				dbType = "time.Time"
				xormType = "DateTime"
			case bytes.HasPrefix(res["Type"], []byte("time")):
				dbType = "string"
			case bytes.HasPrefix(res["Type"], []byte("varchar")):
				dbType = "string"
			case bytes.HasPrefix(res["Type"], []byte("char")):
				dbType = "string"
			case bytes.HasPrefix(res["Type"], []byte("text")):
				dbType = "string"
			}

			if isUpper {
				xormType += " '" + field + "'"
			}
		}

		if isJson {
			str += fieldUpper + " " + dbType + " `xorm:\"" + xormType + "\" json:\"" + field + "\"`\n"
		} else {
			str += fieldUpper + " " + dbType + " `xorm:\"" + xormType + "\"`\n"
		}
	}

	var (
		dbName string
		upperK int
	)
	for k, v := range tableName {
		if 0 == k || upperK == k {
			dbName += strings.ToUpper(string(v))
		} else if "_" == string(v) {
			upperK = k + 1
		} else {
			dbName += string(v)
		}

	}
	commands.Logger.Trace("%v",dbName)
	// Generate the file path
	// e.g. a/b/c

	wr := os.Stdout
	var pkgName string
	if has := strings.HasSuffix(s, "/"); !has {
		s = s + "/"
	}
	fileName :=  s + tableName

	if s != "" {
		generatePath := strings.Split(s, "/")
		length := len(generatePath)
		if length > 0 {
			if f, _ := os.Stat(fileName + ".go"); f != nil {
				commands.Logger.Error("%s.go File already exists, please retry after deletion.", fileName)
				return nil
			}
			//generateFile := generatePath[length-1]
			pkgName = generatePath[length-2]
			wr, _ = os.Create(fileName + ".go")
		}
	}

	data := template.FuncMap{"pkgName": pkgName, "dbName": dbName, "dbFiles": template.HTML(str)}
	if err := utils.Tmpl(dbStruct, data, wr); err != nil {
		commands.Logger.Error("Reverse %s->%s failed.")
	}

	// TODO go fmt file
	if s != "" {
		commands.Logger.Trace("Reverse create %s.go successfully created!", fileName)
	} else {
		commands.Logger.Trace("Reverse successfully created!")
	}
	return nil
}
