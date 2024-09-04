package db

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/prashant1k99/URL-Shortner/types"
	"github.com/prashant1k99/URL-Shortner/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShortURLResponse struct {
	Id           primitive.ObjectID
	ShortenedUrl string
}

// NOTE: This function does not returns a Universal unique code, so in order to tackle the duplicate shortcode issue, Add unique index in mongo for this
func generateShortCode() (string, error) {
	// Get the current timestamp (4 bytes)
	timestamp := uint32(time.Now().Unix())

	// Generate 2 random bytes
	randomBytes := make([]byte, 2)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Combine timestamp and random bytes
	combined := make([]byte, 6)
	binary.BigEndian.PutUint32(combined, timestamp)
	copy(combined[4:], randomBytes)

	// Encode to base62
	return utils.Base62Encode(combined), nil
}

func CreateShortUrl(url *types.Url, baseURL string) (ShortURLResponse, error) {
	urlCollection, err := GetCollection("urls")
	if err != nil {
		return ShortURLResponse{}, err
	}

	// Generate a unique shortcode
	var shortenedUrl string
	retries := 0
	for {
		shortenedUrl, err = generateShortCode()
		if err != nil {
			return ShortURLResponse{}, fmt.Errorf("Error while generating shortcode: %v", err)
		}

		// Check if the shortcode already exists
		count, err := urlCollection.CountDocuments(context.TODO(), bson.M{"shortenedurl": shortenedUrl})
		if err != nil {
			return ShortURLResponse{}, fmt.Errorf("Error while checking shortcode uniqueness: %v", err)
		}
		if count == 0 {
			break
		}
		if retries < 10 {
			// So to break the loop if the unique shortcode is not generated in 10 attempts
			return ShortURLResponse{}, fmt.Errorf("Internal Server Error")
		} else {
			retries++
		}
	}
	// Insert the document to get the Id
	result, err := urlCollection.InsertOne(context.TODO(), types.UrlWithShortVersion{
		Url:          *url,
		ShortenedUrl: shortenedUrl,
	})
	if err != nil {
		return ShortURLResponse{}, err
	}

	insertedId := result.InsertedID.(primitive.ObjectID)
	return ShortURLResponse{Id: insertedId, ShortenedUrl: baseURL + shortenedUrl}, nil
}

func GetAllShortUrlsForUser(userId primitive.ObjectID, baseURL string) ([]types.UrlWithShortVersion, error) {
	urlCollection, err := GetCollection("urls")
	if err != nil {
		return []types.UrlWithShortVersion{}, err
	}
	urlCursor, err := urlCollection.Find(context.TODO(), bson.M{"user": userId})
	if err != nil {
		return []types.UrlWithShortVersion{}, err
	}
	defer urlCursor.Close(context.TODO())

	var urls []types.UrlWithShortVersion
	for urlCursor.Next(context.TODO()) {
		var url types.UrlWithShortVersion
		if err := urlCursor.Decode(&url); err != nil {
			log.Fatal("Error decoding document:", err)
		}
		url.ShortenedUrl = baseURL + url.ShortenedUrl
		urls = append(urls, url)
	}

	// Check for any errors encountered during iteration
	if err := urlCursor.Err(); err != nil {
		fmt.Println("Error while getting cursor")
		return []types.UrlWithShortVersion{}, err
	}
	return urls, nil
}

func GetUrlFromShortUrl(shortUrl string) (types.UrlWithShortVersion, error) {
	urlCollection, err := GetCollection("urls")
	if err != nil {
		return types.UrlWithShortVersion{}, err
	}
	var url types.UrlWithShortVersion
	err = urlCollection.FindOne(context.TODO(), bson.M{"shUrl": shortUrl}).Decode(&url)
	if err != nil {
		return types.UrlWithShortVersion{}, err
	}
	return url, nil
}
