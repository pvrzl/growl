# Growl

#### Example
```
package main

import (
  "fmt"

  "github.com/homina/growl"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)

type TestTable struct {
  Name string `valid:"required" gorm:"unique_index"`
  Id   int    `gorm:"AUTO_INCREMENT"`
}

func (test *TestTable) Db() (db growl.Db) {
  return db.SetData(test)
}

func migrateTestTable() {
  conn, _ := growl.Conn()
  conn.AutoMigrate(TestTable{})
}

func deleteTestTable() {
  conn, _ := growl.Conn()
  conn.DropTable(TestTable{})
}

func main() {
  growl.Config.Path = "conf.yaml"
  err := growl.Config.Load()
  if err != nil {
    fmt.Println(err)
    return
  }

  migrateTestTable()

  test := new(TestTable)
  test.Name = "test01"
  err = test.Db().Save().Error()
  if err != nil {
    fmt.Println(err)
    return
  }

  err = test.Db().Delete().Error()
  if err != nil {
    fmt.Println(err)
    return
  }

  deleteTestTable()
}

```
