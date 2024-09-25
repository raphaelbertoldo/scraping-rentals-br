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
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro de busca 'q' é obrigatório"})
			return
		}

		results, err := searchService.Search(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao realizar a busca"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"results": results,
			"message": "Hello World",
		})
	})

	g.Run(":3000")
}
