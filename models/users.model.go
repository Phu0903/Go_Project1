package models

import (
	"database/sql"
	"fmt"
	"go-module/db"

	_ "github.com/godror/godror"
)

var UserModels = UserModel{}

type User struct {
	UserId       int    `json:"User_Id"`
	UserName     string `json:"User_name"`
	UserFullName string `json:"User_Full_Name"`
	UserEmail    string `json:"User_Email"`
	UserPassword string `json:"User_Password"`
	IsAdmin      int    `json:"Is_Admin"`
	
	
}
type UserModel struct {
	Users []User
}

//Lấy thông tin tất cả user có trong db
func (u *UserModel) Get() ([]User, error) {
	var temp_user []User
	db := db.GetConnection()
	rows, err := db.Query("select * from NEW_JIRA_USER")
	if err == nil {
		for rows.Next() {
			user := User{}
			rows.Scan(&user.UserId,&user.UserName, &user.UserFullName, &user.UserEmail, &user.UserPassword, &user.IsAdmin)
			temp_user = append(temp_user, user)
		}
		//fmt.Printf("%q", rows)
		return temp_user, nil
	} else {
		return nil, err
	}
}

//Thêm User vào db
func (u *UserModel) InsertUser(user User) (sql.Result, error) {
	smt := `INSERT INTO "NEW_JIRA_USER"("USER_NAME","USER_FULL_NAME", "USER_EMAIL", "USER_PASSWORD", "USER_GLOBAL_ROLE") VALUES (:1, :2, :3, :4, :5)`
	db := db.GetConnection()
	fmt.Println(user.IsAdmin)
	return db.Exec(smt,user.UserName, user.UserFullName,user.UserEmail, user.UserPassword, user.IsAdmin)
}
//Check user có trong db hay ko và trả về thông tin người đó
func (um *UserModel) Check_user(ad string, am string) ([]User, error) {
	var temp_admin []User
	query := fmt.Sprintf("SELECT * FROM \"NEW_JIRA_USER\" WHERE \"USER_NAME\" = '%v' OR \"USER_EMAIL\" = '%v'", ad, am)
	db := db.GetConnection()
	rows, err := db.Query(query)
	if err == nil {
		for rows.Next() {
			user := User{}
			rows.Scan(&user.UserId, &user.UserName, &user.UserFullName, &user.UserEmail, &user.UserPassword, &user.IsAdmin)
			temp_admin = append(temp_admin, user)
		}
		return temp_admin, nil
	} else {
		return nil, err
	}
}

//xóa thông tin user
func (sm *UserModel) DeleteUser(id string) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM JIRA_USER WHERE USER_ID = '%v'", id)
	db := db.GetConnection()
	return db.Exec(query)
}

//Update thông tin User
func (sm *UserModel) UpdateUser(id int, strUser string, strFull_name string, strEmail string, strPassword string, isAdmin string) (sql.Result, error) {
	var UserQuery, PasswordQuery, AdminQuery string
	db := db.GetConnection()
	if strUser != "" {
		UserQuery = fmt.Sprintf("USERNAME = '%v',", strUser)
	} else {
		UserQuery = "USERNAME=USERNAME,"
	}
	if strPassword != "" {
		PasswordQuery = fmt.Sprint("USER_PASSWORD = '%v',", strPassword)

	} else {
		PasswordQuery = "USER_PASSWORD=USER_PASSWORD,"
	}
	if isAdmin == "" {
		AdminQuery = "IS_ADMIN=IS_ADMIN,"
	} else {
		AdminQuery = fmt.Sprintf("IS_ADMIN = %v,", isAdmin)
	}
	smt := fmt.Sprintf(`UPDATE "JIRA_USER" SET %v %v %v "USER_FULL_NAME"=:1, "USER_EMAIL"=:2 WHERE "USER_ID"=:3`, UserQuery, PasswordQuery, AdminQuery)
	return db.Exec(smt, strFull_name, strEmail, id)
}

//Check xem user có tồn tại ko
func (ue *UserModel) Check_user_exist(id string) ([]User, error) {
	var temp_exist []User
	fmt.Println(id)
	query := fmt.Sprintf("SELECT * FROM JIRA_USER WHERE USER_ID = '%v'", id)
	db := db.GetConnection()
	rows, err := db.Query(query)

	if err == nil {
		for rows.Next() {
			user_ := User{}

			rows.Scan(&user_.UserId, &user_.UserName, &user_.UserFullName, &user_.UserEmail, &user_.UserPassword, &user_.IsAdmin)

			temp_exist = append(temp_exist, user_)
		}
		return temp_exist, nil
	} else {
		return nil, err
	}
}
