package server

import (
	"groupie-tracker/internal/api"
	"groupie-tracker/internal/structures"
	"net/http"

	"github.com/gin-gonic/gin"
)

// home page
func homePage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var artistsData structures.Artists
		var err error
		artistsData.ArtistsAPI, err = api.GetArtistsAPI()
		if err != nil {
			// c.Status(http.StatusInternalServerError)
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		artistsData.ArtistsArray, err = api.GetArtistsData(artistsData.ArtistsAPI)
		if err != nil {
			// c.Status(http.StatusInternalServerError)
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		// c.JSON(200, artistsData)
		c.HTML(http.StatusOK, "HomePage.html", artistsData)
	}
}

func artistPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, exists := c.GetQuery("id")
		if !exists || id == "" {
			c.Redirect(303, "/")
			// c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "error": "404 page not found"})
			return
		}

		var artistsData structures.Artists
		var err error

		artistsData.ArtistsAPI, err = api.GetArtistsAPI()
		if err != nil {
			// c.Status(http.StatusInternalServerError)
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		artistData, err := api.GetSingleArtistData(artistsData.ArtistsAPI, id)
		if err != nil {
			if err.Error() == http.StatusText(http.StatusNotFound) {
				c.JSON(404, gin.H{"error": "Page not found"})
				return
			}
			// c.Status(http.StatusInternalServerError)
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		c.IndentedJSON(200, artistData)
	}
}
