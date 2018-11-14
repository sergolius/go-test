package sql_builder

import (
	"fmt"
	"strings"
)

const (
	insertFmt      = "INSERT INTO `%s` (%s) VALUES(%s)"
	updateFmt      = "UPDATE `%s` SET %s %s"
	existsFmt      = "SELECT count(`id`) FROM `%s` %s"
	createTableSQL = `CREATE TABLE IF NOT EXISTS %s (
		id INT(10) NOT NULL AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255),
		email VARCHAR(255),
		mobile_number VARCHAR(255)
	)`
)

func GetCreateTableSQL(table string) string {
	return fmt.Sprintf(createTableSQL, table)
}
func GetInsertSQL(table string, columns []string, data map[string]string) string {
	var values []string
	for _, key := range columns {
		if key == "id" {
			values = append(values, data[key])
		} else {
			values = append(values, fmt.Sprintf("%q", data[key]))
		}
	}
	return fmt.Sprintf(insertFmt, table, strings.Join(columns, ","), strings.Join(values, ","))
}

func GetUpdateSQL(table string, columns []string, data map[string]string) string {
	var values []string
	for _, key := range columns {
		if key == "id" {
			values = append(values, fmt.Sprintf("`%s` = %s", key, data[key]))
		} else {
			values = append(values, fmt.Sprintf("`%s` = '%s'", key, data[key]))
		}
	}
	return fmt.Sprintf(updateFmt, table, strings.Join(values, ","), fmt.Sprintf("WHERE `id` = %s", data["id"]))
}

func GetExistsSQL(table string, id string) string {
	return fmt.Sprintf(existsFmt, table, fmt.Sprintf("WHERE `id` = %s", id))
}
