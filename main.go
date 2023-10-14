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
	"github.com/gin-gonic/gin"
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
var collection *mongo.Collection
var MONGO_URI string = "mongodb://tanveeshs:pass123@localhost:27017/test?authSource=admin"
var recipesHandler *handlers.RecipesHandler

func init() {
	ctx = context.Background()
	client, err = mongo.Connect(ctx,
		options.Client().ApplyURI(MONGO_URI))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection = client.Database("test").Collection("recipes")
	recipesHandler = handlers.NewRecipesHandler(ctx,
		collection)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.PATCH("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	router.GET("/recipes/search", recipesHandler.SearchRecipeHandler)
	docs.SwaggerInfo.Title = "Recipes API"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err := router.Run(":3000")
	if err != nil {
		return
	}
}
