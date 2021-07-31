package auth

import (
	"fmt"
	"go-module/helpers"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)
type Claims struct {
	Username string `json:"user_name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}
func CheckUserLoged(c *gin.Context) {
	var tknStr string
	var tknStr1 string
	//serect key
	var jwtKey = []byte("my_secret_key")
	var rule bool
	//request.header
	auth := c.Request.Header["Authorization"]
	//nếu auth > 0 
	if len(auth) > 0 {
		//Cắt Beaer
		tknStr = strings.Trim(auth[0], "Bearer")
		//cắt " "
		tknStr1 = strings.Trim(tknStr, " ")
       
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr1, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		fmt.Println(claims.IsAdmin)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, helpers.MessageResponse{Msg: "Invalid signature"})
				c.Abort()
			}

			rule = false
		} else {
			rule = true
		}

		if rule && tkn != nil {

			c.Next()
		}

		if !rule && tkn != nil {
			c.JSON(http.StatusUnauthorized, helpers.MessageResponse{Msg: "Token expired, please login again"})
			c.Abort()
		}

		if tkn == nil {

			c.JSON(http.StatusUnauthorized, helpers.MessageResponse{Msg: "Token invalid"})
			c.Abort()
		}

	} else {
		c.Abort()
		c.JSON(http.StatusUnauthorized, helpers.MessageResponse{Msg: "Token not found"})
	}
	c.Abort()
}
