package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/vekshinnikita/golang_music"
)

type Authorization interface {
	CreateUser(user golang_music.User) (int, error)
	GetUser(username string, password string) (golang_music.User, error)
}

type Track interface {
	AddTrack(userId int, track golang_music.AddTrackInput) (int, error)
	GetTrack(trackId int) (*golang_music.Track, error)
	DeleteTrack(trackId int) error
	UpdateTrack(userId int, track golang_music.UpdateTrackInput) error
	UpdateTrackFilename(filename string, trackId int, userId int) error
	UpdatePosterFilename(filename string, trackId int, userId int) error
}

type Repository struct {
	Track
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Track:         NewTrackPostgres(db),
	}
}
