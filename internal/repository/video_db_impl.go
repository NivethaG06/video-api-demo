package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"interview-project/pkg/models"
)

func (db *PostgresDB) InsertData(id string, metadata *models.VideoMetadata) error {
	tx, err := db.Sql_db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}

	insert(tx, id, metadata)

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}

func insert(tx *sql.Tx, id string, video *models.VideoMetadata) error {
	err := tx.QueryRow(`
		INSERT INTO video_metadata (id,video_name, description, modified_time, published_time, thumbnail_id)
		VALUES ($1, $2, $3, $4, $5 , $6) RETURNING id;
	`, id, video.VideoName, video.Description, video.ModifiedTime, video.PublishedTime, video.ThumbnailID)
	if err != nil {
		return fmt.Errorf("could not insert video_metadata: %v", err)
	}

	tagIDs := []int{}
	for _, tag := range video.Tags {
		var tagID int
		err = tx.QueryRow(`
			INSERT INTO tags (tag_name)
			VALUES ($1) ON CONFLICT (tag_name) DO NOTHING RETURNING id;
		`, tag)
		if err != nil {
			return fmt.Errorf("could not insert tag: %v", err)
		}
		tagIDs = append(tagIDs, tagID)
	}

	for _, tagID := range tagIDs {
		_, err1 := tx.Exec(`
			INSERT INTO post_tags (post_id, tag_id)
			VALUES ($1, $2);
		`, id, tagID)
		if err1 != nil {
			return fmt.Errorf("could not insert post_tags: %v", err1)
		}
	}
	return nil
}

func (db *PostgresDB) FetchPaginatedData(limitoffset *models.LimitOffset) ([]models.VideoMetadataWithSuggestion, error) {
	query := `
		WITH video_tags AS (
			SELECT vm.id AS video_id, vm.video_name, vm.description, vm.thumbnail_id, 
			       vt.tag_id, t.tag_name, vm.published_time
			FROM video_metadata vm
			LEFT JOIN post_tags vt ON vm.id = vt.video_id
			LEFT JOIN tags t ON vt.tag_id = t.id
		),
		ranked_videos AS (
			SELECT vt.video_id, vt.video_name, vt.description, vt.thumbnail_id, vt.tag_name, 
			       ROW_NUMBER() OVER (PARTITION BY vt.tag_name ORDER BY vt.published_time DESC) AS rank
			FROM video_tags vt
		),
		video_posts AS (
			SELECT video_id, video_name, description, thumbnail_id, tag_name
			FROM ranked_videos
			WHERE rank <= 2
			ORDER BY published_time DESC
			LIMIT $1 OFFSET $2
		)
		SELECT vp.video_id, vp.video_name, vp.description, vp.thumbnail_id, 
		       json_agg(json_build_object('video_name', rv.video_name, 'description', rv.description, 'thumbnail_id', rv.thumbnail_id)) AS suggestions
		FROM video_posts vp
		LEFT JOIN ranked_videos rv ON vp.tag_name = rv.tag_name
		GROUP BY vp.video_id, vp.video_name, vp.description, vp.thumbnail_id;
	`

	rows, err := db.Sql_db.Query(query, limitoffset.Limit, limitoffset.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch videos: %v", err)
	}
	defer rows.Close()

	return fetchHandler(rows)
}

func (db *PostgresDB) CreateIndexPost() error {
	query := `CREATE INDEX IF NOT EXISTS idx_post_id ON video_metadata(post_id);`
	_, err := db.Sql_db.Exec(query)
	return err
}

func (db *PostgresDB) CreateIndexTag() error {
	query := `CREATE INDEX IF NOT EXISTS idx_tag_id ON post_tags(tag_id);`
	_, err := db.Sql_db.Exec(query)
	return err
}

func fetchHandler(rows *sql.Rows) ([]models.VideoMetadataWithSuggestion, error) {
	var videos []models.VideoMetadataWithSuggestion
	for rows.Next() {
		var video models.VideoMetadataWithSuggestion
		var suggestionsJson string

		if err := rows.Scan(&video.VideoName, &video.Description, &video.ThumbnailID, &suggestionsJson); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		if err := json.Unmarshal([]byte(suggestionsJson), &suggestionsJson); err != nil {
			return nil, fmt.Errorf("failed to unmarshal suggestions: %v", err)
		}

		videos = append(videos, video)
	}

	return videos, nil
}
