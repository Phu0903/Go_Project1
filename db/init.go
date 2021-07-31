package db
import (
	"database/sql"
    "fmt"
	"time"
	_ "github.com/godror/godror"
)
//var DbOracle = newDB()
type DBOracle struct{
	ConnStr string
	Db      *sql.DB
}
var db *sql.DB

func InitializeMySQL() {
	connStr := "user=\"jira\" password=\"MyGIAgVwrzQRj8lhD54k\" connectString=\"115.165.166.55:1521/COREDB\""
	dBConnection, err := sql.Open("godror", connStr)
	if err != nil {
		fmt.Println("Connection Failed!!")
	}
	err = dBConnection.Ping()
	if err != nil {
		fmt.Println("Ping Failed!!")
	}
	//db  = dBConnection
	
	db = dBConnection
	dBConnection.SetMaxOpenConns(10)
	dBConnection.SetMaxIdleConns(5)
	dBConnection.SetConnMaxLifetime(time.Second * 10)
}
func GetConnection() *sql.DB {
	if db == nil {
		InitializeMySQL()
	}
	return db
}
/*func (d *DBOracle) GetConnection() (*sql.DB, error){
	connStr := "user=\"jira\" password=\"MyGIAgVwrzQRj8lhD54k\" connectString=\"115.165.166.55:1521/COREDB\""
	db, err := sql.Open("godror", connStr)
	if err != nil{
		return nil,err
	}
	d.Db = db
	return d.Db,nil
}*/
