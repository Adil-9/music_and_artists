package server

import (
	"groupie-tracker/internal/api"
	"groupie-tracker/internal/structures"

	"github.com/gin-gonic/gin"
)

// home page
func homePage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var artistsData structures.Artists
		artistsData.ArtistsAPI = api.GetArtistsAPI()
		artistsData.ArtistsArray = api.GetArtistsData(artistsData.ArtistsAPI)
		// c.JSON(200, artistsData)
		c.IndentedJSON(200, artistsData)
	}
}
