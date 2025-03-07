package garagefwk

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type DataObject struct {
	Id       string
	Reg      *map[string]interface{}
	Metadata *[]map[string]interface{}
	Owners   *[]string
	Perms    *map[string]interface{}
}

type DataObjectFilter struct {
	Field string
	Value interface{}
}

func connectDB(config Config) *sql.DB {
	strConn := fmt.Sprintf("%v:%v@tcp(%v)/%v", config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Db)
	db, err := sql.Open("mysql", strConn)

	if err != nil {
		panic(err)
	}

	return db
}

func createUuid() string {
	uuid, _ := uuid.NewRandom()

	return uuid.String()
}

func insertDataObject(db *sql.DB, user string, reg map[string]interface{}) *string {
	fields := ""
	values := ""
	valueSets := []interface{}{}

	fields += "metadata"
	values += "?"
	valueSets = append(valueSets, "[ { \"user\": \""+user+"\", \"timestamp\": \""+time.Now().Format("2006-01-02 15:04:05")+"\" } ]")

	fields += ",owners"
	values += ",?"
	valueSets = append(valueSets, "[ \""+user+"\" ]")

	for key := range reg {
		fields += "," + key
		values += ",?"
		dt := reflect.TypeOf(reg[key]).String()

		if dt == "*interface {}" || dt == "map[string]interface {}" {
			jsonData, _ := json.Marshal(reg[key])
			valueSets = append(valueSets, string(jsonData))
		} else {
			valueSets = append(valueSets, reg[key])
		}
	}

	_, err := db.Exec("INSERT INTO DataObjects("+fields+") VALUES ("+values+")", valueSets...)
	if err != nil {
		errStr := err.Error()
		return &errStr
	}

	return nil
}

func UpdateDataObject(db *sql.DB, user string, idReg string, reg map[string]interface{}) *string {
	values := ""
	valueSets := []interface{}{}

	values = "metadata = JSON_ARRAY_APPEND( metadata, '$', ?)"
	valueSets = append(valueSets, "[ \"user\": \""+user+"\", \"timestamp\": \""+time.Now().Format("2006-01-02 15:04:05")+"\" ]")

	for key := range reg {
		values += ", " + key + " = ?"
		dt := reflect.TypeOf(reg[key]).String()

		if dt == "*interface {}" || dt == "map[string]interface {}" {
			jsonData, _ := json.Marshal(reg[key])
			valueSets = append(valueSets, string(jsonData))
		} else {
			valueSets = append(valueSets, reg[key])
		}
	}

	valueSets = append(valueSets, idReg)

	_, err := db.Exec("UPDATE DataObjects SET "+values+" WHERE id=?", valueSets...)
	if err != nil {
		errStr := err.Error()
		return &errStr
	}

	return nil
}

func readDataObject(db *sql.DB, idReg string) *DataObject {

	rows, err := db.Query("select id,reg, metadata, owners, perms from DataObjects where id = ?", idReg)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	if rows.Next() {
		return readDataObjectsRow(rows)
	}

	return nil
}

func readDataObjects(db *sql.DB, user string, objecttype string) *[]DataObject {
	rows, err := db.Query("select id,reg, metadata, owners, perms from DataObjects where objecttype = ? AND JSON_CONTAINS(owners, ?, '$');", objecttype, "[\""+user+"\"]")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	dos := make([]DataObject, 0)
	for rows.Next() {
		do := readDataObjectsRow(rows)

		dos = append(dos, *do)
	}

	return &dos
}

func ReadDataObjectsByFilter(db *sql.DB, user string, objecttype string, filter *[]DataObjectFilter) *[]DataObject {
	sql := "select id,reg, metadata, owners, perms from DataObjects where objecttype = ?"
	params := []interface{}{objecttype}

	for _, f := range *filter {
		sql += " AND JSON_EXTRACT(reg, '$." + f.Field + "') = ?"
		params = append(params, f.Value)
	}

	if user != "" {
		sql += "  AND JSON_CONTAINS(owners, ?, '$');"
		params = append(params, "[\""+user+"\"]")
	}

	rows, err := db.Query(sql, params...)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	dos := make([]DataObject, 0)
	for rows.Next() {
		do := readDataObjectsRow(rows)

		dos = append(dos, *do)
	}

	return &dos
}

func readDataObjectsRow(rows *sql.Rows) *DataObject {
	r := DataObject{}

	var reg []byte
	var metadata []byte
	var owners []byte
	var perms []byte

	err := rows.Scan(&r.Id, &reg, &metadata, &owners, &perms)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(reg, &r.Reg)
	json.Unmarshal(metadata, &r.Metadata)
	json.Unmarshal(owners, &r.Owners)
	json.Unmarshal(perms, &r.Perms)

	return &r
}
