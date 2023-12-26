package api

import (
	"context"
	"encoding/json"
	"groupie-tracker/internal/api/structures"
	"groupie-tracker/logger"
	"io"
	"net/http"
	"time"
)

const artistsApiAPI = "https://groupietrackers.herokuapp.com/api"

func GetArtistsAPI() structures.ArtistsAPI {
	var artistsAPI structures.ArtistsAPI

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, artistsApiAPI, nil)
	if err != nil {
		logger.ErrorLog.Println("Error creating request:", err)
		return artistsAPI
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.ErrorLog.Println("Error requsting:", err)
		return artistsAPI
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLog.Println("Error reading body:", err)
		return artistsAPI
	}

	if err = json.Unmarshal(body, &artistsAPI); err != nil {
		logger.ErrorLog.Println("Error unmarshalling body:", err)
		return artistsAPI
	}

	return artistsAPI
}
