package models

type APIVideoRequest struct {
	video VideoMetadata
}

type APIVideoResponse struct {
	video []VideoMetadataWithSuggestion
}

type LimitOffset struct {
	Offset int
	Limit  int
}

type VideoMetadata struct {
	VideoName     string   `json:"name" binding:"required"`
	Description   string   `json:"description"`
	ModifiedTime  int32    `json:"modified_time"`
	PublishedTime int32    `json:"published_time"`
	ThumbnailID   string   `json:"thumbnail_id"`
	Tags          []string `json:"tags"`
}

type VideoMetadataWithSuggestion struct {
	VideoName   string `json:"name" binding:"required"`
	Description string
	ThumbnailID string
	Suggestions []Suggestion
}

type Suggestion struct {
	VideoName   string `json:"name" binding:"required"`
	Description string
	ThumbnailID string
}
