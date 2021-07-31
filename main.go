package main

import (

	"go-module/routes"

)

func main(){
    r := routes.SetupRouter()
    r.Run(":8090")
}