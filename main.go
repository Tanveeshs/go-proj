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
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go-proj/docs"
	"net/http"
	"os"
	"strings"
	"time"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)
	file, _ := os.ReadFile("recipes.json")
	_ = json.Unmarshal(file, &recipes)
}

// ListRecipesHandler @Summary      List all recipes
//
//	@Description	To Get All Recipes
//	@Tags			recipes
//	@Produce		json
//	@Success		200	{array}	Recipe
//	@Router			/recipes [get]
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// UpdateRecipeHandler @Summary      Update Recipe
//
//	@Description	Update Recipe
//	@Tags			recipes
//	@Param			id	path	string	true	"Id of Recipe"	Format(string)
//	@Produce		json
//	@Success		200	{object}	Recipe
//	@Router	/recipes/{id} [put]
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	index := -1
	for i := 0; i < len(recipes); i++ {
		if id == recipes[i].ID {
			index = i
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe Not Found"})
		return
	}
	recipes[index] = recipe
	c.JSON(http.StatusOK, recipe)
}
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1
	for i := 0; i < len(recipes); i++ {
		if id == recipes[i].ID {
			index = i
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe Not Found"})
		return
	}
	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been Deleted",
	})
}
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}
func SearchRecipeHandler(c *gin.Context) {
	tag := c.Query("tag")
	if len(tag) < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tag not defined",
		})
	}
	listOfRecipes := make([]Recipe, 0)
	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
				break
			}
		}
		if found {
			listOfRecipes = append(listOfRecipes, recipes[i])
		}
	}
	c.JSON(http.StatusOK, listOfRecipes)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PATCH("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipeHandler)
	docs.SwaggerInfo.Title = "Recipes API"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err := router.Run(":3000")
	if err != nil {
		return
	}
}
