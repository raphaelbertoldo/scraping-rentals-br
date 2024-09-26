package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raphaelbertoldo/scraping-rentals-br/internal/ivan/search"
)

func main() {
	searchService := search.NewService()
	g := gin.Default()

	g.GET("/", func(c *gin.Context) {
		neighborhood := c.Query("neighborhood")
		min := c.Query("min")
		max := c.Query("max")

		if neighborhood == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "neighborhood query is required"})
			return
		}

		results, err := searchService.Search(neighborhood, min, max)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao realizar a busca"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": results,
		})
	})

	g.Run(":3000")
}
