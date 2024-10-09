package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raphaelbertoldo/scraping-rentals-br/internal/ivan/ivanService"
	"github.com/raphaelbertoldo/scraping-rentals-br/internal/viva/vivaService"
)

func main() {
	searchService := ivanService.NewService()
	vivaService := vivaService.NewService()
	g := gin.Default()

	g.GET("/", func(c *gin.Context) {
		neighborhood := c.Query("neighborhood")
		min := c.Query("min")
		max := c.Query("max")

		if neighborhood == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "neighborhood query is required"})
			return
		}
		results1, err := vivaService.Search(neighborhood, min, max)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao realizar a busca"})
			return
		}
		results, err := searchService.Search(neighborhood, min, max)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao realizar a buscar "})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ivan": results,
			"viva": results1,
		})
	})

	g.Run(":3000")
}
