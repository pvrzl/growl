package growl

import (
	"strings"

	"github.com/jinzhu/inflection"
)

func (db Db) GetTableName() string {
	config := YamlConfig.Growl.Database
	rawSplit := strings.Split(GetStructName(db.data), ".")
	name := strings.ToLower(ToSnake(rawSplit[len(rawSplit)-1]))

	if !config.SingularTable {
		return YamlConfig.Growl.Database.Prefix + inflection.Plural(name)
	}

	return YamlConfig.Growl.Database.Prefix + name
}
