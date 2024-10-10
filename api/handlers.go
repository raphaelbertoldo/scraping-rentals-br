package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raphaelbertoldo/scraping-rentals-br/api/internal/ivan/ivanService"
	"github.com/raphaelbertoldo/scraping-rentals-br/api/internal/viva/vivaService"
)

type Server struct {
	Router        *gin.Engine
	searchService *ivanService.Service
	vivaService   *vivaService.Service
}

func NewServer() *Server {
	s := &Server{
		Router:        gin.Default(),
		searchService: ivanService.NewService(),
		vivaService:   vivaService.NewService(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.Router.GET("/rentals", s.handleSearch)
	s.Router.GET("/", s.checkHealth)
}

func (s *Server) handleSearch(c *gin.Context) {
	fmt.Println("🚀 ~ EXECUTANDO FUNC handleSearch : ")
	neighborhood := c.Query("neighborhood")
	min := c.Query("min")
	max := c.Query("max")

	if neighborhood == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "neighborhood query is required"})
		return
	}

	vivaResults, vivaErr := s.vivaService.Search(neighborhood, min, max)
	ivanResults, ivanErr := s.searchService.Search(neighborhood, min, max)

	if vivaErr != nil && ivanErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while searching both services"})
		return
	}

	response := gin.H{}
	if vivaErr == nil {
		response["viva"] = vivaResults
	}
	if ivanErr == nil {
		response["ivan"] = ivanResults
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) checkHealth(c *gin.Context) {
	response := gin.H{}
	response["health"] = "okkkkk cu de egua"
	c.JSON(http.StatusOK, response)
}

// Handler é a função que o Vercel vai chamar
func Handler(w http.ResponseWriter, r *http.Request) {
	gin.SetMode(gin.ReleaseMode)
	server := NewServer()
	server.Router.ServeHTTP(w, r)
}
