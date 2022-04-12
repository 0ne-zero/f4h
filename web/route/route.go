package route

import (
	"os"
	"path/filepath"

	"github.com/0ne-zero/f4h/config/constansts"
	"github.com/0ne-zero/f4h/utilities/log"
	templatefunction "github.com/0ne-zero/f4h/utilities/template_function"
	"github.com/0ne-zero/f4h/web/controllers"
	"github.com/0ne-zero/f4h/web/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func MakeRoute() *gin.Engine {
	r := gin.Default()
	// Html template function
	templatefunction.AddFunctionsToRoute(r)
	// Statics
	r.Static("statics", filepath.Join(constansts.ExecutableDirectory, "/statics/"))
	// Htmls
	r.LoadHTMLGlob("statics/templates/**/*.html")

	// Too many request
	r.Use(middleware.TooManyRequest())
	// Not found
	r.NoRoute(middleware.NotFound())

	// // Use session
	session_key := "secret"
	if gin.Mode() == "PRODUCTION" {
		if sk := os.Getenv("F4H_SESSION_KEY"); sk == "" {
			log.Log(logrus.Fatal, "F4H_SESSION_KEY isn't exists in environment variables")
		}
	}
	store := cookie.NewStore([]byte(session_key))
	store.Options(sessions.Options{MaxAge: 0})
	r.Use(sessions.Sessions("authentication", store))

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
	}
	constansts.Routes = r.Routes()

	return r
}
