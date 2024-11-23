package services

import (
	"context"
	"encoding/json"
	"fmt"
	"interview-project/internal/cache"
	"interview-project/internal/repository"
	"interview-project/pkg/models"
	"interview-project/pkg/util"
	"strconv"
	"time"
)

func CreateVideo(video *models.VideoMetadata) error {
	newConnection, _ := repository.LoadConnection()
	err := newConnection.InsertData(util.GenerateVideoUUID(), video)
	if err != nil {
		return fmt.Errorf("failed to insert video data: %v", err)
	}

	go timestampUpdate()

	return nil
}

func timestampUpdate() {
	ctx := context.Background()
	timestamp := time.Now().Unix() // Current Unix timestamp
	err := cache.SetCache(ctx, "global_video_timestamp", fmt.Sprintf("%d", timestamp), 0)
	if err != nil {
		fmt.Printf("Error updating global timestamp: %v\n", err)
	}
}

//func GetVideos(ctx context.Context, limitoffset *models.LimitOffset) ([]models.VideoMetadataWithSuggestion, error) {
//	var videos []models.VideoMetadataWithSuggestion
//	newConnection := NewPostgresDB()
//
//	videos, err := newConnection.FetchPaginatedData(limitoffset)
//	if err != nil {
//		return nil, err
//	}
//
//	cacheKey := "video_suggestions"
//	var videoWithSuggestions []string
//	// Attempt to get cached suggestions from Redis
//	cachedSuggestions, err := cache.GetCache(ctx, cacheKey)
//	if err != nil {
//		// If the cache is empty or another error occurs, generate new suggestions
//		for _, video := range videos {
//			videoWithSuggestions = append(videoWithSuggestions, strings.Join(video, ", "))
//		}
//
//		// Store the generated suggestions in Redis (as a JSON string)
//		suggestionsJSON, _ := json.Marshal(videoWithSuggestions)
//		err = cache.SetCache(ctx, cacheKey, string(suggestionsJSON), 0)
//		if err != nil {
//			return nil, fmt.Errorf("could not set cache: %v", err)
//		}
//	} else {
//		// If cached data exists, unmarshal into the suggestions
//		err = json.Unmarshal([]byte(cachedSuggestions), &videoWithSuggestions)
//		if err != nil {
//			return nil, fmt.Errorf("failed to unmarshal Redis cached data: %v", err)
//		}
//	}
//
//	// Return videos and suggestions
//	return videos, nil
//
//}

func GetVideos(ctx context.Context, limitoffset *models.LimitOffset) ([]models.VideoMetadataWithSuggestion, error) {
	var videos []models.VideoMetadataWithSuggestion
	newConnection, _ := repository.LoadConnection()

	// Fetch videos with pagination
	//videos, err := newConnection.FetchPaginatedData(limitoffset)
	//if err != nil {
	//	return nil, err
	//}

	// Cache key for video suggestions
	cacheKey := "video_suggestions"
	cacheTimestampKey := "global_video_timestamp"

	// Get the global timestamp from Redis
	globalTimestampStr, _ := cache.GetCache(ctx, cacheTimestampKey)
	globalTimestamp, _ := strconv.ParseInt(globalTimestampStr, 10, 64)

	// Get the cached suggestions and timestamp
	cachedSuggestions, err := cache.GetCache(ctx, cacheKey)
	cachedTimestampStr, _ := cache.GetCache(ctx, cacheKey+"_timestamp")
	cachedTimestamp, _ := strconv.ParseInt(cachedTimestampStr, 10, 64)

	if err == nil && cachedSuggestions != "" && cachedTimestamp >= globalTimestamp {
		//var videoWithSuggestions []string
		err = json.Unmarshal([]byte(cachedSuggestions), &videos)
		if err == nil {
			return videos, nil
		}
	} else {
		videos, err := newConnection.FetchPaginatedData(limitoffset)
		if err != nil {
			return nil, fmt.Errorf("could not set cache: %v", err)
		}
		err = UpdateNewRedisCache(ctx, cacheKey, videos)
		if err != nil {
			return nil, fmt.Errorf("could not set cache: %v", err)
		}
		err = UpdateCacheTimestamp(ctx, cacheKey)
		if err != nil {
			return nil, fmt.Errorf("could not update cache timestamp: %v", err)
		}
	}

	//// Cache is invalid or missing; regenerate it
	//var videoWithSuggestions []string
	//for _, video := range videos {
	//	videoWithSuggestions = append(videoWithSuggestions, strings.Join(video, ", "))
	//}
	//
	//// Update the cache with new suggestions
	//suggestionsJSON, _ := json.Marshal(videoWithSuggestions)
	//err = cache.SetCache(ctx, cacheKey, string(suggestionsJSON), 3600) // TTL: 1 hour
	//if err != nil {
	//	return nil, fmt.Errorf("could not set cache: %v", err)
	//}
	//
	//// Update the cache timestamp
	//err = cache.SetCache(ctx, cacheKey+"_timestamp", fmt.Sprintf("%d", time.Now().Unix()), 3600)
	//if err != nil {
	//	return nil, fmt.Errorf("could not update cache timestamp: %v", err)
	//}

	// Return the videos
	return videos, nil
}

func UpdateNewRedisCache(ctx context.Context, cacheKey string, videoWithSuggestions []models.VideoMetadataWithSuggestion) error {
	suggestionsJSON, _ := json.Marshal(videoWithSuggestions)
	return cache.SetCache(ctx, cacheKey, string(suggestionsJSON), 3600) // TTL: 1 hour
}

func UpdateCacheTimestamp(ctx context.Context, cacheKey string) error {
	return cache.SetCache(ctx, cacheKey+"_timestamp", fmt.Sprintf("%d", time.Now().Unix()), 0)
}
