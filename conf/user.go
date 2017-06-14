package conf

type User struct {
	Id int64 `xorm:"pk autoincr"`
Name string `xorm:"varchar(255)"`

}
