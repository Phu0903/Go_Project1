package routes
import(
	. "go-module/handler"
	"go-module/middeware"
	"github.com/gin-gonic/gin"

)

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {

        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
/*var UserRouter = UserRoute{GroupName: "/api/users"}

type UserRoute struct {
	GroupName   string
	RouterGroup *gin.RouterGroup
}*/
func SetupRouter() *gin.Engine{
	r:= gin.Default()
	r.Use(CORSMiddleware())
	v1 := r.Group("/")
	{
		v1.GET("user",auth.CheckUserLoged,UserHandlers.Index())
		v1.POST("register",UserHandlers.SignUp())
		v1.DELETE("user",UserHandlers.DeleteUser())
		v1.POST("user",UserHandlers.Login())
	}
	return r
}
/*func (u *UserRoute) Init(router *gin.Engine) {
	u.RouterGroup = router.Group(u.GroupName)
	{
		//u.RouterGroup.Use(auth.CheckUserLoged, auth.CheckAdmin)
	
		u.RouterGroup.GET("/", u.Index())
		
	}
}

func (u *UserRoute) Index() gin.HandlerFunc {
	//u.RouterGroup.Use(auth.CheckUserLoged, auth.CheckAdmin)
	return UserHandlers.Index()
}*/