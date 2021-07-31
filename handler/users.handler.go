package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-module/helpers"
	"go-module/models"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/godror/godror"
)

var jwt_secret_Key = []byte("my_secret_key")

var TokenValidTime time.Duration = 1000
var UserHandlers = UserHandler{}

type Claims struct {
	Username string `json:"user_name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}
type UserHandler struct{}

func (u *UserHandler) Index() gin.HandlerFunc {
	//Do everything here, call model etc...
	return func(c *gin.Context) {
		users, err := models.UserModels.Get()
		if err == nil {
			//fmt.Println(users)
			//c.JSON(http.StatusOK, users)
			c.JSON(http.StatusOK, helpers.MessageResponse{Msg: "Get data Success", Data: users})
		} else {
			fmt.Println(err)
		}
	}
}

//SignUp
func (u *UserHandler) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

		var user, full_name, email, password string
		var IA int
		var myMap map[string]string
		json.NewDecoder(c.Request.Body).Decode(&myMap)
		user = fmt.Sprintf("%v", myMap["user"])
		full_name = fmt.Sprintf("%v", myMap["fullname"])
		email = fmt.Sprintf("%v", myMap["email"])
		password = fmt.Sprintf("%v", myMap["password"])
		//check Info User
		if user == "" || full_name == "" || email == "" || password == "" {
			c.JSON(http.StatusBadRequest, helpers.MessageResponse{Msg: "The parameters are not enough"})
		} else {
		
			Exist_user, err := models.UserModels.Check_user(user, email)
			IA = 1
			if err != nil {
				c.JSON(http.StatusBadRequest, helpers.MessageResponse{Msg: "Error running query"})
			}
			//UserName đã tồn tại
			if len(Exist_user) > 0 {
				if Exist_user[0].UserName == user && Exist_user[0].UserEmail != email {
					c.JSON(http.StatusConflict, helpers.MessageResponse{Msg: "User already exists, please choose another user"})
				}
			}
			//Email đã tồn tại
			if len(Exist_user) > 0 {
				if Exist_user[0].UserName != user && Exist_user[0].UserEmail == email {
					c.JSON(http.StatusConflict, helpers.MessageResponse{Msg: "Email already exists, please choose another email"})
				}
			}
			//Email và user name đã tồn tại
			if len(Exist_user) > 0 {
				if Exist_user[0].UserName == user && Exist_user[0].UserEmail == email {
					c.JSON(http.StatusConflict, helpers.MessageResponse{Msg: "User and email already exists, please choose another user and email"})
				}
			}
			//được phép đăng ký
			if len(Exist_user) == 0 {
				sample := password                                                 //password
				encodedString := base64.StdEncoding.EncodeToString([]byte(sample)) //encodeString password
				scr := models.User{UserName: user, UserFullName: full_name, UserEmail: email, UserPassword: encodedString, IsAdmin: IA}
				sm := models.UserModel{}
				if _, err := sm.InsertUser(scr); err != nil {
					c.JSON(http.StatusBadRequest, helpers.MessageResponse{Msg: "Error running query"})

				} else {
					//xét time hết hạn token
					expirationTime := time.Now().Add(TokenValidTime * time.Hour)
					//tạo claims
					claims := &Claims{
						Username: user,
						IsAdmin:  1,
						StandardClaims: jwt.StandardClaims{
							// In JWT, the expiry time is expressed as unix milliseconds
							ExpiresAt: expirationTime.Unix(),
						},
					}
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

					signedToken, err := token.SignedString(jwt_secret_Key)
					if err != nil {
						c.JSON(http.StatusBadRequest, helpers.MessageResponse{Msg: "have error with encode token"})
					} else {

						c.JSON(http.StatusOK, helpers.MessageResponse{Msg: "Sign Up Success", Data: struct{ Token string }{signedToken}})

					}

				}

			}
		}
	}
}

//login
func (u *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		//request body
		var myMap map[string]string
		var user, password string
		json.NewDecoder(c.Request.Body).Decode(&myMap)
		user = fmt.Sprintf("%v", myMap["user"])
		password = fmt.Sprintf("%v", myMap["password"])
		//kiểm tra user email passowrd null
		if user == "" || password == "" {
			c.JSON(http.StatusBadRequest, helpers.MessageResponse{Msg: "The parameters are not enough"})
		} else {
			//Check exist user
			Exist_user, err := models.UserModels.Check_user(user, user)

			if err != nil {
				c.JSON(http.StatusBadRequest, helpers.MessageResponse{Msg: "Error running query"})
			}
			//Nếu không có user nào
			if len(Exist_user) == 0 {
				c.JSON(http.StatusNotAcceptable, helpers.MessageResponse{Msg: "User not exists, please choose another user"})

			}
			//Có user
			if len(Exist_user) == 1 {
				//Check admin
				if Exist_user[0].IsAdmin == 0 {
					// mã hóa password
					originalStringBytes, err := base64.StdEncoding.DecodeString(Exist_user[0].UserPassword)
					if err != nil {
						c.JSON(http.StatusNotAcceptable, helpers.MessageResponse{Msg: "User not exists, please choose another user"})
						log.Fatalf("Some error occured during base64 decode. Error %s", err.Error())

					}
					//kiểm tra có trùng password ko
					if password == string(originalStringBytes) {
						//xét time hết hạn token
						expirationTime := time.Now().Add(TokenValidTime * time.Hour)
						//tạo claims
						claims := &Claims{
							Username: user,
							IsAdmin:  0,
							StandardClaims: jwt.StandardClaims{
								// In JWT, the expiry time is expressed as unix milliseconds
								ExpiresAt: expirationTime.Unix(),
							},
						}
						// Declare the token with the algorithm used for signing, and the claims

						token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

						signedToken, err := token.SignedString(jwt_secret_Key)
						if err != nil {
							// If there is an error in creating the JWT return an internal server error
							c.JSON(http.StatusBadRequest, helpers.MessageResponse{Msg: "have error with encode token"})
						} else {
							c.JSON(http.StatusOK, helpers.MessageResponse{Msg: "Login Success", Data: struct{ Token string }{signedToken}})
						}
					} else {
						c.JSON(http.StatusConflict, helpers.MessageResponse{Msg: "Wrong password, please try again!"})
					}
				}
			}
		}

	}
}

func (u *UserHandler) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		fmt.Println(id)
		Exist_user_exist, err := models.UserModels.Check_user_exist(id)
		if id == "" {
			c.JSON(http.StatusBadRequest, helpers.MessageResponse{Msg: "The parameters are not enough"})
		} else {
			if err != nil {
				c.JSON(http.StatusBadRequest, helpers.MessageResponse{Msg: "Error running query"})
			} else {
				if len(Exist_user_exist) == 0 {
					c.JSON(http.StatusBadRequest, helpers.MessageResponse{Msg: "User Id does not exist"})
				}
				if len(Exist_user_exist) > 0 {
					if _, err := models.UserModels.DeleteUser(id); err != nil {
						c.JSON(http.StatusBadRequest, helpers.MessageResponse{Msg: "Error running query"})
					} else {
						c.JSON(http.StatusOK, helpers.MessageResponse{Msg: "Delete user success"})

					}
				}
			}
		}

	}
}
