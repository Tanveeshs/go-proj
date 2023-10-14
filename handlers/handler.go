package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go-proj/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

type RecipesHandler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewRecipesHandler(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client) *RecipesHandler {
	return &RecipesHandler{
		collection:  collection,
		ctx:         ctx,
		redisClient: redisClient,
	}
}

// ListRecipesHandler @Summary      List all recipes
//
//	@Description	To Get All Recipes
//	@Tags			recipes
//	@Produce		json
//	@Success		200	{array}	Recipe
//	@Router			/recipes [get]
func (handler *RecipesHandler) ListRecipesHandler(c *gin.Context) {
	val, err := handler.redisClient.Get(handler.ctx, "recipes").Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		cur, err := handler.collection.Find(handler.ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": err.Error()})
			return
		}
		defer func(cur *mongo.Cursor, ctx context.Context) {
			err := cur.Close(ctx)
			if err != nil {
				fmt.Print(err)
				return
			}
		}(cur, handler.ctx)
		recipes := make([]models.Recipe, 0)
		for cur.Next(handler.ctx) {
			var recipe models.Recipe
			err := cur.Decode(&recipe)
			if err != nil {
				return
			}
			recipes = append(recipes, recipe)
		}
		valToSave, _ := json.Marshal(recipes)
		handler.redisClient.Set(handler.ctx, "recipes", string(valToSave), 0)
		c.JSON(http.StatusOK, recipes)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	} else {
		recipes := make([]models.Recipe, 0)
		err := json.Unmarshal([]byte(val), &recipes)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, recipes)
	}

}

// UpdateRecipeHandler @Summary      Update Recipe
//
//	@Description	Update Recipe
//	@Tags			recipes
//	@Param			id	path	string	true	"Id of Recipe"	Format(string)
//	@Produce		json
//	@Success		200	{object}	Recipe
//	@Router	/recipes/{id} [put]
func (handler *RecipesHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	objectID, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.UpdateOne(handler.ctx, bson.M{
		"_id": objectID,
	}, bson.D{{"$set", bson.D{
		{"name", recipe.Name},
		{"instructions", recipe.Instructions},
		{"ingredients", recipe.Ingredients},
		{"tags", recipe.Tags}}}})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
}

func (handler *RecipesHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.DeleteOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe Not Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been Deleted",
	})
}

func (handler *RecipesHandler) NewRecipeHandler(c *gin.Context) {
	if c.GetHeader("X-API-KEY") == "DEMO" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "API key not provided or invalid"})
		return
	}
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err := handler.collection.InsertOne(handler.ctx, recipe)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error while inserting a new recipe"})
		return
	}
	handler.redisClient.Del(handler.ctx, "recipes")
	c.JSON(http.StatusOK, recipe)
}

func (handler *RecipesHandler) SearchRecipeHandler(c *gin.Context) {
	tag := c.Query("tag")
	if len(tag) < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tag not defined",
		})
	}
	listOfRecipes := make([]models.Recipe, 0)
	cur, err := handler.collection.Find(handler.ctx, bson.M{
		"tags": tag,
	})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			fmt.Print(err)
			return
		}
	}(cur, handler.ctx)
	for cur.Next(handler.ctx) {
		var recipe models.Recipe
		err := cur.Decode(&recipe)
		if err != nil {
			return
		}
		listOfRecipes = append(listOfRecipes, recipe)
	}

	c.JSON(http.StatusOK, listOfRecipes)
}
