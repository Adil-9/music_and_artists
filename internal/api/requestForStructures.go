package api

import (
	"context"
	"encoding/json"
	"groupie-tracker/internal/api/structures"
	"groupie-tracker/internal/cache"
	"groupie-tracker/logger"
	"io"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

const artistsApiAPI = "https://groupietrackers.herokuapp.com/api"

func GetArtistsAPI() structures.ArtistsAPI {
	var artistsAPI structures.ArtistsAPI

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	//retrieveing cached data if exists
	data, err := cache.RedisClient.Get(ctx, artistsApiAPI).Result()
	if err == redis.Nil {
		// Key does not exist
		// fmt.Println("Key not found in Redis")
		//do nothing
	} else if err != nil { // Other errors
		logger.ErrorLog.Println("Error retrieving data from redis cache:", err)
		// return artistsAPI
	} else { // Key exists, print the value
		if err = json.Unmarshal([]byte(data), &artistsAPI); err != nil {
			logger.ErrorLog.Println("Error unmarshalling data from cache:", err)
		} else {
			return artistsAPI
		}
	}

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

	marshaledData, err := json.Marshal(artistsAPI)
	if err != nil {
		logger.ErrorLog.Println("Error marshaling data:", err)
	} else {
		//setting time duration for cache expareation as 10 min and sending data to cache
		cache.RedisClient.Set(context.Background(), artistsApiAPI, marshaledData, time.Minute*10)
	}

	return artistsAPI
}
