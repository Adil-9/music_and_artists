package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"groupie-tracker/internal/cache"
	"groupie-tracker/internal/structures"
	"groupie-tracker/logger"
	"io"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	errorGetFromCache           = "Error retrieving data from redis cache:"
	errorUnmarshallingFromCache = "Error unmarshalling data from cache:"
	errorCreatingRequet         = "Error creating request:"
	errorRequest                = "Error requsting:"
	errorReadingBody            = "Error reading body:"
	errorUnmarshalingBody       = "Error unmarshalling body:"
	errorMarshallingData        = "Error marshaling data:"
)

const artistsApiAPI = "https://groupietrackers.herokuapp.com/api"

func GetArtistsAPI() (structures.ArtistsAPI, error) {
	var artistsAPI structures.ArtistsAPI

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	//retrieveing cached data if exists
	data, err := cache.RedisClient.Get(ctx, artistsApiAPI).Result()
	if err == redis.Nil { // Key does not exist
		// fmt.Println("Key not found in Redis")
		// do nothing
	} else if err != nil { // Other errors
		logger.ErrorLog.Println(errorGetFromCache, err)
		// return artistsAPI
	} else { // Key exists, print the value
		if err = json.Unmarshal([]byte(data), &artistsAPI); err != nil {
			logger.ErrorLog.Println(errorUnmarshallingFromCache, err)
		} else {
			return artistsAPI, nil
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, artistsApiAPI, nil)
	if err != nil {
		logger.ErrorLog.Println(errorCreatingRequet, err)
		return artistsAPI, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.ErrorLog.Println(errorRequest, err)
		return artistsAPI, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLog.Println(errorReadingBody, err)
		return artistsAPI, err
	}

	if err = json.Unmarshal(body, &artistsAPI); err != nil {
		logger.ErrorLog.Println(errorUnmarshalingBody, err)
		return artistsAPI, err
	}

	// marshaledData, err := json.Marshal(artistsAPI)
	// if err != nil {
	// 	logger.ErrorLog.Println(errorMarshallingData, err)
	// } else {
	// 	//setting time duration for cache expareation as 10 min and sending data to cache
	// 	cache.RedisClient.Set(context.Background(), artistsApiAPI, marshaledData, time.Minute*10)
	// }

	return artistsAPI, nil
}

func GetArtistsData(ArtistsApi structures.ArtistsAPI) ([]structures.Artist, error) {
	var artists []structures.Artist
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	data, err := cache.RedisClient.Get(ctx, ArtistsApi.Artists).Result()
	if err == redis.Nil { //key does not exist
		// do nothing
	} else if err != nil {
		logger.ErrorLog.Println(errorGetFromCache, err)
		// return artistsAPI //do not do this becouse we still can send request using api and get data
	} else {
		if err = json.Unmarshal([]byte(data), &artists); err != nil {
			logger.ErrorLog.Println(errorUnmarshallingFromCache, err)
		} else {
			return artists, nil
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ArtistsApi.Artists, nil)
	if err != nil {
		logger.ErrorLog.Println(errorCreatingRequet, err)
		return artists, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.ErrorLog.Println(errorRequest, err)
		return artists, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLog.Println(errorReadingBody, err)
		return artists, err
	}

	if err = json.Unmarshal(body, &artists); err != nil {
		logger.ErrorLog.Println(errorUnmarshalingBody, err)
		return artists, err
	}

	// marshalled, err := json.Marshal(artists)
	// if err != nil {
	// 	logger.ErrorLog.Println(errorUnmarshalingBody, err)
	// } else {
	// 	cache.RedisClient.Set(context.Background(), ArtistsApi.Artists, marshalled, time.Minute*10)
	// }

	return artists, nil
}

func GetSingleArtistData(ArtistsApi structures.ArtistsAPI, id string) (structures.Artist, error) {
	var Artist structures.Artist
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	link := fmt.Sprintf("%s/%s", ArtistsApi.Artists, id)
	
	data, err := cache.RedisClient.Get(ctx, link).Result()
	if err == redis.Nil { //key does not exist
		//do nothing
	} else if err != nil {
		logger.ErrorLog.Println(errorGetFromCache, err)
	} else {
		err = json.Unmarshal([]byte(data), &Artist)
		if err != nil {
			logger.ErrorLog.Println(errorUnmarshalingBody)
		} else {
			return Artist, nil
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		logger.ErrorLog.Println(errorCreatingRequet, err)
		return Artist, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.ErrorLog.Println(errorRequest)
		return Artist, err
	} else if resp.StatusCode != http.StatusOK {
		logger.InfoLog.Printf("From %s %s request was sent, status %s", req.RemoteAddr, link, resp.Status)
		return Artist, errors.New(http.StatusText(http.StatusNotFound))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLog.Println(errorReadingBody, err)
		return Artist, err
	}

	if err = json.Unmarshal(body, &Artist); err != nil {
		logger.ErrorLog.Println(errorUnmarshalingBody, err)
		return Artist, nil
	}

	cache.RedisClient.Set(context.Background(), link, body, time.Minute*10)

	return Artist, nil
}
