package route

import (
	"os"
	"path/filepath"

	"github.com/0ne-zero/f4h/constansts"
	"github.com/0ne-zero/f4h/utilities/functions/general"
	template_func "github.com/0ne-zero/f4h/utilities/functions/template"
	wrapper_logger "github.com/0ne-zero/f4h/utilities/wrapper_logger"
	"github.com/0ne-zero/f4h/web/controller"
	"github.com/0ne-zero/f4h/web/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func MakeRoute() *gin.Engine {
	r := gin.Default()
	// Max memory for multipart forms
	r.MaxMultipartMemory = 12 << 20
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
	session_key := "secret"
	if gin.Mode() == "PRODUCTION" {
		if sk := os.Getenv("F4H_SESSION_KEY"); sk == "" {
			wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "F4H_SESSION_KEY isn't exists in environment variables", Fields: nil, ErrorLocation: general.GetCallerInfo(0)})
		}
	}
	store := memstore.NewStore([]byte(session_key))
	store.Options(sessions.Options{MaxAge: 0})
	r.Use(sessions.Sessions(constansts.AppName+"_Session_KEY", store))

	// Public routes
	r.GET("/Login", controller.Login_GET)
	r.POST("/Login", controller.Login_POST)
	r.POST("/Register", controller.Register_POST)
	r.GET("/admin", controller.Admin_Index)
	r.Use(middleware.SetSession())
	// Public authorized
	authorized := r.Group("/")
	authorized.Use(middleware.Authentication())
	{
		// Public authorized routes
		// Surface
		authorized.GET("/", controller.Index)
		authorized.GET("/Products/*category", controller.ProductList)
		authorized.GET("/ProductDetails/:id", controller.ProductDetails)
		authorized.GET("/Profile", controller.Profile_GET)
		authorized.GET("/Wishlist", controller.Wishlist)
		authorized.GET("/Cart", controller.Cart)
		authorized.GET("/AddProduct", controller.AddProduct_GET)
		authorized.POST("/AddProduct", controller.AddProduct_POST)
		authorized.GET("/EditProduct/:id", controller.EditProduct_GET)
		authorized.POST("/EditProduct/", controller.EditProduct_POST)

		// Forum
		authorized.GET("/Discussions", controller.Discussions)
		authorized.GET("/DiscussionForums/:discussion", controller.DiscussionForums)
		authorized.GET("/ForumTopics/:forum", controller.ForumTopics)
		authorized.GET("/Topic/:topic_id", controller.ShowTopic)
		authorized.GET("/AddTopic/:forum", controller.AddTopic_GET)
		authorized.POST("/AddTopic/:forum", controller.AddTopic_POST)
		authorized.GET("/EditTopic/:topic_id", controller.EditTopic_Get)
		authorized.POST("/EditTopic/:topic_id", controller.EditTopic_POST)

		// Sub-Routes
		authorized.GET("/DeleteCartItem/:id", controller.DeleteCartItem)
		authorized.GET("/IncreaseCartItemQuantity/:id", controller.IncreaseCartItemQuantity)
		authorized.GET("/DecreaseCartItemQuantity/:id", controller.DecreaseCartItemQuantity)
		authorized.POST("/AddToCart", controller.AddToCart)
		authorized.POST("/AddProductComment/", controller.AddProductComment)
		authorized.GET("/AddToWishlist/:p_id", controller.AddToWishlist)
		authorized.GET("/DeleteTopic/:id", controller.DeleteTopic)
		authorized.GET("/DeleteProduct/:id", controller.DeleteProduct)
		authorized.POST("/EditAccount", controller.EditAccount_POST)
		authorized.POST("/EditAvatar", controller.EditAvatar_POST)
		authorized.POST("/ManageWallet", controller.ManageWallet_POST)
		authorized.POST("/ManageAddress", controller.ManageAddress_POST)
		authorized.GET("/IncreaseProductCommentPositiveVote/:id", controller.IncreaseProductCommentPositiveVote)
		authorized.GET("/DecreaseProductCommentPositiveVote/:id", controller.DecreaseProductCommentPositiveVote)
		authorized.GET("/IncreaseProductCommentNegativeVote/:id", controller.IncreaseProductCommentPositiveVote)
		authorized.GET("/DecreaseProductCommentNegativeVote/:id", controller.DecreaseProductCommentNegativeVote)
		authorized.GET("/IncreaseTopicCommentPositiveVote/:id", controller.IncreaseTopicCommentPositiveVote)
		authorized.GET("/DecreaseTopicCommentPositiveVote/:id", controller.DecreaseTopicCommentPositiveVote)
		authorized.GET("/IncreaseTopicCommentNegativeVote/:id", controller.IncreaseTopicCommentNegativeVote)
		authorized.GET("/DecreaseTopicCommentNegativeVote/:id", controller.DecreaseTopicCommentNegativeVote)
	}
	constansts.Routes = r.Routes()

	return r
}
