package app

import (
	"BookVault-API/database"
	"BookVault-API/handler"
	"BookVault-API/middleware"
	"BookVault-API/service"
	
	"net/http"
	"gorm.io/gorm"
)

var Routes =  map[string]string {
	"NotFound": 		"/",
	"Home":				"/user",
	"AdminHome":		"/admin",

	"Register":			"/user/register",
	"Login":			"/user/login",
	"CreateDetails":	"/user/createDetails/{userID}",
	"GetUserByID":		"/user/getById/{userID}",

	"CreateBook":		"/book/create",
	"GetByTitle":		"/book?title={bookTitle}",
	"GetBooks":			"/book/all",
	"GetBooksByAuthor":	"/book/author?author={bookAuthor}",
	"UpdateStock":		"/book/updateStock/{bookID}",

	"AddToCart":		"/cart/add/{userID}/{bookID}",
	"ClearCart":		"/cart/clear/{userID}",
	"RemoveFromCart":	"/cart/remove/{userID}/{bookID}",
	"UpdateQuantity":	"/cart/update/{userID}/{bookID}?quantity={quantity}",
	"GetCart":			"/cart/{cartID}",

	"CreateOrder":		"/order/create/{userID}?address={address}",
	"CancelOrder":		"/order/cancel/{orderID}",
	"GetOrder":			"/order/{orderID}",
	"GetUserOrders":	"/order/user/{userID}",
	"GetOrdersByStatus":"/order/status/?status={status}",
	"UpdateStatus":		"/order/update/{orderID}?status={status}",

	"AddReview":		"/review/add/{userID}/{bookID}",
	"GetReviewsByBook": "/review/get/{bookID}",
	"GetReviewsByUser": "/review/getByUser/{userID}",
	"UpdateReview":		"/review/update/{userID}/{bookID}",
	"DeleteReviewByID":	"/review/delete/{reviewID}",
}

type App struct {
	DB *gorm.DB

	UserService 	service.UserService
	BookService 	service.BookService
	CartService		service.CartService
	OrderService	service.OrderService
	ReviewService	service.ReviewService

	HomeHandler 	*handler.HomeHandler
	UserHandler		*handler.UserHandler
	BookHandler 	*handler.BookHandler
	CartHandler		*handler.CartHandler
	OrderHandler	*handler.OrderHandler
	ReviewHandler	*handler.ReviewHandler
}


func NewApp() *App {
	db := database.InitDB()

	userService 	:= service.NewUserService(db)
	bookService 	:= service.NewBookService(db)
	cartService 	:= service.NewCartService(db)
	orderService 	:= service.NewOrderService(db)
	reviewService 	:= service.NewReviewService(db)

	homeHandler 	:= handler.NewHomeHandler()
	userHandler 	:= handler.NewUserHandler(userService)
	bookHandler 	:= handler.NewBookHandler(bookService)
	cartHandler 	:= handler.NewCartHandler(cartService)
	orderHandler 	:= handler.NewOrderHandler(orderService)
	reviewHandler	:= handler.NewReviewHandler(reviewService)

	return &App{
		DB: db,
		UserService: userService,
		BookService: bookService,
		CartService: cartService,
		OrderService: orderService,
		ReviewService: reviewService,

		HomeHandler: homeHandler,
		UserHandler: userHandler,
		BookHandler: bookHandler,
		CartHandler: cartHandler,
		OrderHandler: orderHandler,
		ReviewHandler: reviewHandler,
	}
}


func (a *App) Init() {
	mux := http.NewServeMux()

	//homeHandlers
	mux.HandleFunc("/", 		a.HomeHandler.NotFound)
	mux.HandleFunc("/user", 	middleware.AuthMiddleware("admin", "user")(a.HomeHandler.Home))
	mux.HandleFunc("/admin", 	middleware.AuthMiddleware("admin")(a.HomeHandler.AdminHome))

	//userHandlers
	mux.HandleFunc("/user/register", 		a.UserHandler.Register)
	mux.HandleFunc("/user/login", 			a.UserHandler.Login)
	mux.HandleFunc("/user/createDetails/",	middleware.AuthMiddleware("admin", "user")(a.UserHandler.CreateDetails))
	mux.HandleFunc("/user/getById/", 		middleware.AuthMiddleware("admin", "user")(a.UserHandler.GetUserByID))

	//bookHandlers
	mux.HandleFunc("/book/create", 			middleware.AuthMiddleware("admin")(a.BookHandler.CreateBook))
	mux.HandleFunc("/book", 				middleware.AuthMiddleware("admin", "user")(a.BookHandler.GetByTitle))
	mux.HandleFunc("/book/all", 			middleware.AuthMiddleware("admin", "user")(a.BookHandler.GetBooks))
	mux.HandleFunc("/book/author", 			middleware.AuthMiddleware("admin", "user")(a.BookHandler.GetBooksByAuthor))
	mux.HandleFunc("/book/updateStock/", 	middleware.AuthMiddleware("admin")(a.BookHandler.UpdateStock))

	//cartHandlers
	mux.HandleFunc("/cart/add/", 		middleware.AuthMiddleware("admin", "user")(a.CartHandler.AddToCart))
	mux.HandleFunc("/cart/clear/", 		middleware.AuthMiddleware("admin", "user")(a.CartHandler.ClearCart))
	mux.HandleFunc("/cart/remove/", 	middleware.AuthMiddleware("admin", "user")(a.CartHandler.RemoveFromCart))
	mux.HandleFunc("/cart/update/", 	middleware.AuthMiddleware("admin", "user")(a.CartHandler.UpdateQuantity))
	mux.HandleFunc("/cart/", 			middleware.AuthMiddleware("admin", "user")(a.CartHandler.GetCart))

	//orderHandlers
	mux.HandleFunc("/order/create/", 	middleware.AuthMiddleware("admin", "user")(a.OrderHandler.CreateOrder))
	mux.HandleFunc("/order/cancel/", 	middleware.AuthMiddleware("admin", "user")(a.OrderHandler.CancelOrder))
	mux.HandleFunc("/order/", 			middleware.AuthMiddleware("admin", "user")(a.OrderHandler.GetOrder))
	mux.HandleFunc("/order/user/", 		middleware.AuthMiddleware("admin", "user")(a.OrderHandler.GetUserOrders))
	mux.HandleFunc("/order/status/", 	middleware.AuthMiddleware("admin", "user")(a.OrderHandler.GetOrdersByStatus))
	mux.HandleFunc("/order/update/", 	middleware.AuthMiddleware("admin")(a.OrderHandler.UpdateStatus))

	//reviewHandlers
	mux.HandleFunc("/review/add/",			middleware.AuthMiddleware("admin", "user")(a.ReviewHandler.AddReview))
	mux.HandleFunc("/review/get/", 			middleware.AuthMiddleware("admin", "user")(a.ReviewHandler.GetReviewsByBook))
	mux.HandleFunc("/review/getByUser/", 	middleware.AuthMiddleware("admin", "user")(a.ReviewHandler.GetReviewsByUser))
	mux.HandleFunc("/review/update/", 		middleware.AuthMiddleware("admin", "user")(a.ReviewHandler.UpdateReview))
	mux.HandleFunc("/review/delete/", 		middleware.AuthMiddleware("admin", "user")(a.ReviewHandler.DeleteReviewByID))	

	http.ListenAndServe(":8080", mux)
}