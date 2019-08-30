package repository

import (
	"database/sql"
	"fmt"
	"github.com/guru-invest/guru.api.preregistry/configuration"
	_ "github.com/lib/pq"
	"log"
)

var database *sql.DB

func connect(){
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s", configuration.CONFIGURATION.DATABASE.Username,
		configuration.CONFIGURATION.DATABASE.Password, configuration.CONFIGURATION.DATABASE.Database, configuration.CONFIGURATION.DATABASE.Url, configuration.CONFIGURATION.DATABASE.Port)
	db, err := sql.Open("postgres",dsn)
	if err != nil{
		log.Println("Error connecting to database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Println("Error on ping database: %v", err)
	}
	database = db
}


func mapResult(rows *sql.Rows, queryName string) map[string][]map[string]interface{}{
	cols, err := rows.Columns()
	if err != nil {
		log.Println("Error on mapping results: %v", err)
	}
	m := make(map[string]interface{})
	mapped := make(map[string][]map[string]interface{})

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			log.Println("Error on mapping results: %v", err)
		}
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		mapped[queryName] = append(mapped[queryName], m)
		m = make(map[string]interface{})
	}
	return mapped
}
