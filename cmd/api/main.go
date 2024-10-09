package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raphaelbertoldo/scraping-rentals-br/cmd/internal/ivan/ivanService"
	"github.com/raphaelbertoldo/scraping-rentals-br/cmd/internal/viva/vivaService"
)

type Server struct {
	router        *gin.Engine
	searchService *ivanService.Service
	vivaService   *vivaService.Service
}

func NewServer() *Server {
	s := &Server{
		router:        gin.Default(),
		searchService: ivanService.NewService(),
		vivaService:   vivaService.NewService(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.GET("/", s.handleSearch)
	s.router.GET("/health", s.checkHealth)
}
func (s *Server) checkHealth(c *gin.Context) {
	response := gin.H{}
	response["health"] = "okkkkk cu de egua"
	c.JSON(http.StatusOK, response)

}
func (s *Server) handleSearch(c *gin.Context) {
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

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

func main() {
	server := NewServer()
	if err := server.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	server := NewServer()
	server.router.ServeHTTP(w, r)
}
