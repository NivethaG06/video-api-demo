CREATE TABLE IF NOT EXISTS video_metadata (
                                              id SERIAL PRIMARY KEY,
                                              video_name TEXT NOT NULL,
                                              description TEXT,
                                              modified_time INT,
                                              published_time INT,
                                              thumbnail_id TEXT
);

CREATE TABLE IF NOT EXISTS tags (
                                    id SERIAL PRIMARY KEY,
                                    tag_name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS video_tags (
                                          id SERIAL PRIMARY KEY,
                                          post_id INT REFERENCES video_metadata(id),
    tag_id INT REFERENCES tags(id),
    creation_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX IF NOT EXISTS idx_video_metadata_id ON video_metadata(id);
CREATE INDEX IF NOT EXISTS idx_tags_id ON tags(id);
CREATE INDEX IF NOT EXISTS idx_video_tags_tag_id ON video_tags(tag_id);