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
	ErrorArtistDoesNotExist     = "artist does not exits"
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

func GetSingleArtistFullData(ArtistsApi structures.ArtistsAPI, id string) (structures.ArtistFullData, error) {
	var ArtistFullData structures.ArtistFullData

	linkArtist := fmt.Sprintf("%s/%s", ArtistsApi.Artists, id)
	linkLocation := fmt.Sprintf("%s/%s", ArtistsApi.Locations, id)
	linkConcertDates := fmt.Sprintf("%s/%s", ArtistsApi.Dates, id)
	linkRelation := fmt.Sprintf("%s/%s", ArtistsApi.Relation, id)

	data, err := cache.RedisClient.Get(context.Background(), linkArtist).Result()
	if err == redis.Nil { //key does not exist
		//do nothing
	} else if err != nil {
		logger.ErrorLog.Println(errorGetFromCache, err)
	} else {
		err = json.Unmarshal([]byte(data), &ArtistFullData)
		if err != nil {
			logger.ErrorLog.Println(errorUnmarshalingBody)
		} else {
			return ArtistFullData, nil
		}
	}
	artist, err := getArtistData(linkArtist, id)
	if err != nil {
		logger.ErrorLog.Println("Error getting artist data", err)
		return ArtistFullData, err
	}
	ArtistFullData.Id = artist.Id
	ArtistFullData.Image = artist.Image
	ArtistFullData.Name = artist.Name
	ArtistFullData.Members = artist.Members
	ArtistFullData.CreationDate = artist.CreationDate
	ArtistFullData.FirstAlbum = artist.FirstAlbum
	
	location, err := getLocation(linkLocation, id)
	if err != nil {
		logger.ErrorLog.Println("Error getting artist location data", err)
		return ArtistFullData, err
	}
	ArtistFullData.Locations = location

	concerDates, err := getConcertDates(linkConcertDates, id)
	if err != nil {
		logger.ErrorLog.Println("Error getting artist concert date data", err)
		return ArtistFullData, err
	}
	ArtistFullData.ConcertDates = concerDates

	relation, err := getRelation(linkRelation, id)
	if err != nil {
		logger.ErrorLog.Println("Error getting artist relation data", err)
		return ArtistFullData, err
	}
	ArtistFullData.Relations = relation

	// resp, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	logger.ErrorLog.Println(errorRequest)
	// 	return Artist, err
	// } else if resp.StatusCode != http.StatusOK {
	// 	logger.InfoLog.Printf("From %s %s request was sent, status %s", req.RemoteAddr, linkArtist, resp.Status)
	// 	return Artist, errors.New(http.StatusText(http.StatusNotFound))
	// }

	// body, err := io.ReadAll(resp.Body)

	// if len(body) == 126 {
	// 	return Artist, errors.New(ErrorArtistDoesNotExist)
	// }

	// if err != nil {
	// 	logger.ErrorLog.Println(errorReadingBody, err)
	// 	return Artist, err
	// }

	// if err = json.Unmarshal(body, &Artist); err != nil {
	// 	logger.ErrorLog.Println(errorUnmarshalingBody, err)
	// 	return Artist, nil
	// }

	body, err := json.Marshal(ArtistFullData)
	if err != nil {
		logger.ErrorLog.Println(errorMarshallingData, err)
		return ArtistFullData, err
	}

	cache.RedisClient.Set(context.Background(), linkArtist, body, time.Minute*10)

	return ArtistFullData, nil
}

func getArtistData(linkArtist string, id string) (structures.Artist, error) {
	var Artist structures.Artist
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, linkArtist, nil)
	if err != nil {
		logger.ErrorLog.Println(errorRequest, err)
		return Artist, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.ErrorLog.Println(errorRequest)
		return Artist, err
	} else if resp.StatusCode != http.StatusOK {
		logger.InfoLog.Printf("From %s %s request was sent, status %s", req.RemoteAddr, linkArtist, resp.Status)
		return Artist, errors.New("internal server error")
	}

	body, err := io.ReadAll(resp.Body)

	if len(body) == 126 {
		return Artist, errors.New(ErrorArtistDoesNotExist)
	}

	if err != nil {
		logger.ErrorLog.Println(errorReadingBody, err)
		return Artist, err
	}

	if err = json.Unmarshal(body, &Artist); err != nil {
		logger.ErrorLog.Println(errorUnmarshalingBody, err)
		return Artist, nil
	}
	return Artist, nil
}

func getLocation(linkLocation string, id string) (structures.ArtistsLocation, error) {
	var Artist structures.ArtistsLocation
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, linkLocation, nil)
	if err != nil {
		logger.ErrorLog.Println(errorRequest, err)
		return Artist, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.ErrorLog.Println(errorRequest)
		return Artist, err
	} else if resp.StatusCode != http.StatusOK {
		logger.InfoLog.Printf("From %s %s request was sent, status %s", req.RemoteAddr, linkLocation, resp.Status)
		return Artist, errors.New("internal server error")
	}

	body, err := io.ReadAll(resp.Body)

	if len(body) == 126 {
		return Artist, errors.New(ErrorArtistDoesNotExist)
	}

	if err != nil {
		logger.ErrorLog.Println(errorReadingBody, err)
		return Artist, err
	}

	if err = json.Unmarshal(body, &Artist); err != nil {
		logger.ErrorLog.Println(errorUnmarshalingBody, err)
		return Artist, nil
	}
	return Artist, nil
}

func getConcertDates(linkConcertDates string, id string) (structures.ArtistsConcertDates, error) {
	var Artist structures.ArtistsConcertDates
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, linkConcertDates, nil)
	if err != nil {
		logger.ErrorLog.Println(errorRequest, err)
		return Artist, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.ErrorLog.Println(errorRequest)
		return Artist, err
	} else if resp.StatusCode != http.StatusOK {
		logger.InfoLog.Printf("From %s %s request was sent, status %s", req.RemoteAddr, linkConcertDates, resp.Status)
		return Artist, errors.New("iternal server error")
	}

	body, err := io.ReadAll(resp.Body)

	if len(body) == 126 {
		return Artist, errors.New(ErrorArtistDoesNotExist)
	}

	if err != nil {
		logger.ErrorLog.Println(errorReadingBody, err)
		return Artist, err
	}

	if err = json.Unmarshal(body, &Artist); err != nil {
		logger.ErrorLog.Println(errorUnmarshalingBody, err)
		return Artist, nil
	}
	return Artist, nil
}

func getRelation(linkRelation string, id string) (structures.ArtistsRealtion, error) {
	var Artist structures.ArtistsRealtion
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, linkRelation, nil)
	if err != nil {
		logger.ErrorLog.Println(errorRequest, err)
		return Artist, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.ErrorLog.Println(errorRequest)
		return Artist, err
	} else if resp.StatusCode != http.StatusOK {
		logger.InfoLog.Printf("From %s %s request was sent, status %s", req.RemoteAddr, linkRelation, resp.Status)
		return Artist, errors.New("internal server error")
	}

	body, err := io.ReadAll(resp.Body)

	if len(body) == 126 {
		return Artist, errors.New(ErrorArtistDoesNotExist)
	}

	if err != nil {
		logger.ErrorLog.Println(errorReadingBody, err)
		return Artist, err
	}

	if err = json.Unmarshal(body, &Artist); err != nil {
		logger.ErrorLog.Println(errorUnmarshalingBody, err)
		return Artist, nil
	}
	return Artist, nil
}