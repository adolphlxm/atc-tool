package orm

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/adolphlxm/atc-tool/commands"
	"github.com/adolphlxm/atc-tool/utils"

	"encoding/json"
	"github.com/adolphlxm/atc/orm"
	_ "github.com/adolphlxm/atc/orm/xorm"
)

var (
	j bool
)

var CmdOrm = &commands.Command{
	Usage: "orm [-json] driverName datasourceName tableName [generatedPath]",
	Use:   "reverse a database table to codes",
	Options: `

    -s                Generated one go file for every table
    driverName        Database driver name, now supported four: mysql mymysql sqlite3 postgres
    datasourceName    Database connection uri, for detail infomation please visit driver's project page
    tmplPath          Template dir for generated. the default templates dir has provide 1 template
    generatedPath     This parameter is optional, if blank, the default value is model, then will
                      generated all codes in model dir
`,
	Run: Run,
}

var dbStruct = `package {{.pkgName}}

type {{.dbName}} struct {
	{{.dbFiles}}
}
`

func init() {
	CmdOrm.Flag.BoolVar(&j, "json", false, "Support Json tag true or false.")
	commands.Register(CmdOrm)
}

var db orm.Orm

type OrmConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"dbname"`
	LogLevel string `json:"loglevel"`
}

func Run(cmd *commands.Command, args []string) int {
	if len(args) < 3 {
		fmt.Println("Usage: orm [-json] driverName datasourceName tableName [generatedPath]")
		fmt.Println(commands.ErrUseError)
		return 1
	}

	cmd.Flag.Parse(args)
	args = cmd.Flag.Args()
	if j {
		args = args[1:]
	}

	// Datasource parsing.
	sourceName := args[1]
	source1 := strings.Split(sourceName, "@") // root:123	tcp(%s)/test?charset=utf8
	source2 := strings.Split(source1[1], "/") // tcp(%s)	test?charset=utf8
	lenTcp := len(source2[0])
	user := strings.Split(source1[0], ":")
	i1 := strings.Index(source2[1], "?")
	cf := &OrmConfig{
		Driver:   args[0],
		Host:     source2[0][4 : lenTcp-1],
		User:     user[0],
		Password: user[1],
		DbName:   source2[1][:i1],
		LogLevel: "LOG_DEBUG",
	}
	cfJson, err := json.Marshal(cf)
	if err != nil {
		commands.Logger.Error("datasourceName json err :%v", err.Error())
		return 1
	}

	db, _ = orm.NewOrm("xorm")
	if err := db.Open(cf.DbName, string(cfJson)); err != nil {
		//commands.Logger.Error("init database orm err :%v", err.Error())
		return 1
	}
	var gen string
	if len(args) > 3 {
		gen = args[3]
	}
	if err := reverse(cf.DbName, args[2], gen, j); err != nil {
		return 1
	}

	return 2
}

//生成表结构
func reverse(aliasName, tableName, g string, isJson bool) error {
	var str string

	engine := db.Use(aliasName)
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

	// Generate the file path
	// e.g. a/b/c
	wr := os.Stdout
	var pkgName string
	if has := strings.HasSuffix(g, "/"); !has {
		g = g + "/"
	}
	fileName := g + tableName

	if g != "/" {
		generatePath := strings.Split(g, "/")
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
		commands.Logger.Error("Reverse %s->%s failed.", aliasName, tableName)
	}

	// TODO go fmt file
	if g != "/" {
		commands.Logger.Trace("Reverse create %s.go successfully!", fileName)
	} else {
		commands.Logger.Trace("Reverse struct successfully!")
	}
	return nil
}
