# Growl

[![Build Status](https://travis-ci.org/homina/growl.svg?branch=master)](https://travis-ci.org/homina/growl)
[![codecov](https://codecov.io/gh/homina/growl/branch/master/graph/badge.svg)](https://codecov.io/gh/homina/growl)

#### Overview
Growl is another layer for https://github.com/jinzhu/gorm, https://github.com/go-redis/redis and https://github.com/patrickmn/go-cache

- Simple config file for db, redis
- ORM like
- Auto set/get cache on query

##### Installation
```bash
go get github.com/homina/growl
```

##### Import package in your project
```go
import (
    "github.com/homina/growl"
)
```

#### Config file 
```yaml
growl:
  database:
    driver: mysql 
    url: root:@/
    name: "growl_test?charset=utf8&parseTime=True&loc=Local&sql_mode='ALLOW_INVALID_DATES'" 
    prefix:  
    singulartable: false # default : false
    limit: 10 
    maxlimit: 32 
    timelayout: "02/01/2006"
  redis:
    host: localhost 
    port: "6379" 
    password: 
    channel: "channel"
    enable: true 
  misc:
    localcache: true # enable go-cache
    log: true # enable gorm log
    flushatinit: true # flush cache at start

```


#### Example
```go
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

  fmt.Printf("%+v", test)
  // &{Name:test01 Id:1}

  test.Name = "test02"
  err = test.Db().Model(test).Update().Error()
  if err != nil {
    fmt.Println(err)
    return
  }

  err = test.Db().First().Error()
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Printf("%+v", test)
  // &{Name:test02 Id:1}

  err = test.Db().Delete().Error()
  if err != nil {
    fmt.Println(err)
    return
  }

  deleteTestTable()
}
```

#### Validation
Reference : https://github.com/asaskevich/govalidator
