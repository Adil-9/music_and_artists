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
			artistsData.Error = "Internal server error"
			c.HTML(500, "HomePage.html", artistsData)
			return
		}

		artistsData.ArtistsArray, err = api.GetArtistsData(artistsData.ArtistsAPI)
		if err != nil {
			// c.Status(http.StatusInternalServerError)
			artistsData.Error = "Internal server error"
			c.HTML(500, "HomePage.html", artistsData)
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
			artistsData.Error = "Internal server error"
			c.HTML(500, "SingleArtist.html", artistsData)
			return
		}

		artistData, err := api.GetSingleArtistData(artistsData.ArtistsAPI, id)
		if err != nil && err.Error() == api.ErrorArtistDoesNotExist {
			if err.Error() == api.ErrorArtistDoesNotExist {
				artistData.Error = "Page not found"
				c.HTML(404, "SingleArtist.html", artistData)
				return
			}
			// c.Status(http.StatusInternalServerError)
			artistData.Error = "Internal server error"
			c.HTML(500, "SingleArtist.html", artistData)
			return
		}

		c.HTML(200, "SingleArtist.html", artistData)
	}
}
