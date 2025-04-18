package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/lastvoidtemplar/sabbac/internal/fingerprint"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	logger *slog.Logger
	db     *sql.DB
}

func New(dbPath string, logger *slog.Logger) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		logger.With(slog.String("db_path", dbPath), slog.String("err", err.Error())).Error("Error while creating a db")
		return nil, err
	}

	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		panic(err)
	}

	logger.With(slog.String("db_path", dbPath)).Info("DB is created successfully")
	return &DB{
		logger: logger,
		db:     db,
	}, nil
}

func (db *DB) InsertFingerprint(hash uint32, timestamps []fingerprint.VideoTimestamp) error {
	str := fmt.Sprintf("INSERT INTO fingerprints (hash_key, song_id, video_timestamp) VALUES (%d, '%s', %d)",
		hash, timestamps[0].SongId, timestamps[0].AnchorTime)

	for i := 1; i < len(timestamps); i++ {
		str += fmt.Sprintf(", (%d, '%s', %d)", hash, timestamps[i].SongId, timestamps[i].AnchorTime)
	}

	_, err := db.db.Exec(str)

	if err != nil {
		db.logger.With(
			slog.String("hash", fmt.Sprintf("%x", hash)),
			slog.String("err", err.Error()),
		).Error("Error while inserting fingerprint")
		return err
	}

	return nil
}
