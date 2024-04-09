package service

import (
	"bytes"
	"mime/multipart"

	"github.com/vekshinnikita/golang_music"
	"github.com/vekshinnikita/golang_music/pkg/repository"
)

type Authorization interface {
	CreateUser(user golang_music.User) (int, error)
	GenerateToken(username string, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Track interface {
	AddTrack(userId int, track golang_music.AddTrackInput) (int, error)
	GetTrackInfo(userId int, trackId int) (*golang_music.Track, error)
	UpdateTrack(userId int, track golang_music.UpdateTrackInput) error
	DeleteTrack(userId int, trackId int) error
	UploadTrack(userId int, trackId int, file *multipart.FileHeader, rangeBytes []int64, fileSize int64) error
	UploadPoster(userId int, trackId int, file *multipart.FileHeader) error
	GetPoster(userId int, trackId int) (*bytes.Buffer, error)
	StreamingTrack(userId int, trackId int, bytesRange []*int64) (*bytes.Buffer, int64, int64, string, int64, error)
}

type Service struct {
	Track
	Authorization
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Track:         NewTrackService(repo.Track),
	}
}
