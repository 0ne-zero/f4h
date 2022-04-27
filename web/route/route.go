package route

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/0ne-zero/f4h/constansts"
	"github.com/0ne-zero/f4h/utilities/functions/general"
	template_func "github.com/0ne-zero/f4h/utilities/functions/template"
	wrapper_logger "github.com/0ne-zero/f4h/utilities/wrapper_logger"
	"github.com/0ne-zero/f4h/web/controllers"
	"github.com/0ne-zero/f4h/web/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func MakeRoute() *gin.Engine {
	r := gin.Default()
	// Html template function
	template_func.AddFunctionsToRoute(r)
	// Statics
	r.Static("statics", filepath.Join(constansts.ExecutableDirectory, "/statics/"))
	// Htmls
	r.LoadHTMLGlob("statics/templates/**/*.html")

	// Too many request
	r.Use(middleware.TooManyRequest())
	// Not found
	r.NoRoute(middleware.NotFound())

	// Use session
	session_key := "s"
	if gin.Mode() == "PRODUCTION" {
		if sk := os.Getenv("F4H_SESSION_KEY"); sk == "" {
			err_file_info, err := general.GetCallerInfo(0)
			if err != nil {
				log_msg := fmt.Sprintf("%s file='%s:%d'", "Error occurred during get caller info ", "/web/route/route.go", 38)
				general.AppendTextToFile(constansts.LogFilePath, log_msg)
				os.Exit(1)
			}
			wrapper_logger.Log(&wrapper_logger.FatalLevel{}, "F4H_SESSION_KEY isn't exists in environment variables", &err_file_info)
		}
	}
	store := memstore.NewStore([]byte(session_key))
	store.Options(sessions.Options{MaxAge: 0})
	r.Use(sessions.Sessions("authsdafentication", store))

	// Public routes
	r.GET("/login", controllers.Login_GET)
	r.POST("/login", controllers.Login_POST)
	r.POST("/register", controllers.Register_POST)

	// Public authorized
	authorized := r.Group("/")
	authorized.Use(middleware.Authentication())
	{
		// Public authorized routes
		authorized.GET("/", controllers.Index)
		authorized.GET("/products/*category", controllers.ProductList)
		authorized.GET("/productDetails/:id", controllers.ProductDetails)
		authorized.GET("/Discussions", controllers.Discussions)
		authorized.GET("/DiscussionForums/:discussion", controllers.DiscussionForums)
		authorized.GET("/ForumTopics/:forum", controllers.ForumTopics)
		authorized.GET("/AddTopic/:forum", controllers.AddTopic_GET)
		authorized.POST("/AddTopic/:forum", controllers.AddTopic_POST)
	}
	constansts.Routes = r.Routes()

	return r
}
