//	@title			Recipes API
//	@description	This is a recipes API Project.
//	@host			localhost:3000
//	@BasePath		/
//	@version		1.0.0
//	@contact.name	Tanveesh Singh Chaudhery
//	@contact.email	tanveeshs@gmail.com
//	@accept			application/json
//	@produce		application/json
//	@schemes		http

package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	redisstore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go-proj/docs"
	"go-proj/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/net/context"
	"log"
)

var ctx context.Context
var err error
var client *mongo.Client
var MongoUri string = "mongodb://tanveeshs:pass123@localhost:27017/test?authSource=admin"
var recipesHandler *handlers.RecipesHandler
var authHandler *handlers.AuthHandler

func init() {
	ctx = context.Background()
	client, err = mongo.Connect(ctx,
		options.Client().ApplyURI(MongoUri))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database("test").Collection("recipes")
	userCollection := client.Database("test").Collection("users")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0})
	status := redisClient.Ping(ctx)
	fmt.Println(status)
	recipesHandler = handlers.NewRecipesHandler(ctx,
		collection, redisClient)
	authHandler = handlers.NewAuthHandler(userCollection, ctx)
}

func main() {
	router := gin.Default()

	store, _ := redisstore.NewStore(10, "tcp",
		"localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("recipes_api", store))
	router.POST("/signin", authHandler.SignInHandler)
	router.GET("/signout", authHandler.SignOutHandler)
	router.GET("/refresh", authHandler.RefreshHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
		//authorized.PATCH("/recipes/:id", recipesHandler.UpdateRecipeHandler)
		//authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
		//authorized.GET("/recipes/search", recipesHandler.SearchRecipeHandler)
	}

	docs.SwaggerInfo.Title = "Recipes API"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err := router.Run(":3000")
	if err != nil {
		return
	}
}
