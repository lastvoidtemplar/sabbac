CREATE TABLE fingerprints (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hash_key INTEGER NOT NULL,         
    song_id TEXT NOT NULL,           
    video_timestamp INTEGER NOT NULL
);

CREATE INDEX idx_hash_key ON fingerprints(hash_key);


SELECT COUNT(*) FROM fingerprints