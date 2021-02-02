package dbdriver

import (
	"database/sql"
	"fmt"
	"github.com/fibanez6/go-dbexporter/domain"
	"log"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "test"
)

var DB *sql.DB

const sqlStatement = `
		with device_select as (
		    SELECT id 
			FROM device
			WHERE name = $1 AND lastipaddress =$2
		), device_insert as (
		 	INSERT INTO device(name,lastipaddress)
		 	SELECT $1,$2
		 	WHERE NOT EXISTS (
				SELECT id
				FROM device 
				WHERE name = $1 AND lastipaddress =$2 limit 1)
		 	RETURNING id
		), device_id as (
			SELECT id FROM device_select union all SELECT * FROM device_insert
		), monitor_select as (
		    SELECT id 
			FROM monitor
			WHERE serialnumber = $3 AND resolution = $4
		), monitor_insert as (
		 	INSERT INTO monitor(serialnumber, resolution)
		 	SELECT $3, $4
		 	WHERE NOT EXISTS (
				SELECT id 
				FROM monitor 
				WHERE serialnumber = $3 AND resolution = $4)
			RETURNING id 
		), monitor_id as (
			SELECT id FROM monitor_select union all SELECT * FROM monitor_insert
		)
		INSERT INTO device_monitor(device_id,monitor_id)
		SELECT * from device_id as d, monitor_id as m
		ON CONFLICT DO NOTHING`

/*
Init is called prior to main.
*/
func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	DB = db
	log.Printf("Successfully connected to Postgres!")
}

func Close() {
	DB.Close()
}

/*
Given a device and a monitor, it executes the sql statement and saves them in the database
*/
func Write(device domain.Device, monitor domain.Monitor) error {
	_, err := DB.Query(sqlStatement, device.Name, device.LastIpAddress, monitor.SerialNumber, monitor.Resolution)
	return err
}
