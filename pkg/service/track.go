package service

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"github.com/spf13/viper"
	"github.com/vekshinnikita/golang_music"
	"github.com/vekshinnikita/golang_music/pkg/repository"
	"github.com/vekshinnikita/golang_music/pkg/tools"
)

type TrackService struct {
	repo repository.Track
}

func NewTrackService(repo repository.Track) *TrackService {
	return &TrackService{
		repo,
	}
}

func (s *TrackService) AddTrack(userId int, track golang_music.AddTrackInput) (int, error) {
	return s.repo.AddTrack(userId, track)
}

func (s *TrackService) UploadTrack(userId int, trackId int, fileHeader *multipart.FileHeader, rangeBytes []int64, fileSize int64) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	} else {
		defer file.Close()
	}

	var filename string

	if rangeBytes[0] == 0 {
		//get extension of file
		fileExtension, err := tools.GetFileExtension(&file)
		if err != nil {
			return err
		}

		filename = uuid.New().String() + fileExtension

		//checking file mimetype
		if err := tools.VerifyMIMEType(&file, "audio"); err != nil {
			return err
		}
	} else {
		// get filename from DB
		track, err := s.repo.GetTrack(trackId)
		if err != nil {
			return err
		}
		filename = track.TrackFilename
	}

	folderPath := strings.Replace(viper.GetString("media.tracks_folder"), "{user_id}", strconv.Itoa(userId), 1)
	filepath := filepath.Join(folderPath, filename)

	err = tools.SaveFile(file, filepath, rangeBytes[0])
	if err != nil {
		return err
	}

	if rangeBytes[0] == 0 {
		return s.repo.UpdateTrackFilename(filename, trackId, userId)
	}

	return nil
}

func (s *TrackService) UploadPoster(userId int, trackId int, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	//checking file mime type
	if err := tools.VerifyMIMEType(&file, "image"); err != nil {
		return err
	}

	//get extension of file
	fileExtension, err := tools.GetFileExtension(&file)
	if err != nil {
		return err
	}

	// resize image
	buff := new(bytes.Buffer)
	imgJpg, _, _ := image.Decode(file)
	file.Seek(0, io.SeekStart)

	imgJpg = resize.Resize(600, 600, imgJpg, resize.Bicubic)
	jpeg.Encode(buff, imgJpg, nil)

	reader := bytes.NewReader(buff.Bytes())

	filename := uuid.New().String() + fileExtension

	folderPath := strings.Replace(viper.GetString("media.posters_folder"), "{user_id}", strconv.Itoa(userId), 1)
	filepath := filepath.Join(folderPath, filename)

	err = tools.SaveFile(reader, filepath, 0)
	if err != nil {
		return err
	}

	return s.repo.UpdatePosterFilename(filename, trackId, userId)
}

func (s *TrackService) UpdateTrack(userId int, track golang_music.UpdateTrackInput) error {
	return s.repo.UpdateTrack(userId, track)
}

func (s *TrackService) DeleteTrack(userId int, trackId int) error {
	track, err := s.repo.GetTrack(trackId)
	if err != nil {
		return err
	}

	if track.UserId != userId {
		return errors.New("permission denied")
	}

	return s.repo.DeleteTrack(trackId)
}

func (s *TrackService) GetPoster(userId int, trackId int) (*bytes.Buffer, error) {
	track, err := s.repo.GetTrack(trackId)
	if err != nil {
		return nil, err
	}

	if !track.Public && track.UserId != userId {
		return nil, errors.New("permission denied")
	}

	var buf bytes.Buffer

	folderPath := strings.Replace(viper.GetString("media.posters_folder"), "{user_id}", strconv.Itoa(userId), 1)
	filePath := filepath.Join(folderPath, track.PosterFilename)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(&buf, file)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

func (s *TrackService) GetTrackInfo(userId int, trackId int) (*golang_music.Track, error) {
	track, err := s.repo.GetTrack(trackId)
	if err != nil {
		return nil, err
	}

	if !track.Public && track.UserId != userId {
		return nil, errors.New("permission denied")
	}

	return track, nil
}

func (s *TrackService) StreamingTrack(userId int, trackId int, bytesRange []*int64) (*bytes.Buffer, int64, int64, string, int64, error) {
	track, err := s.repo.GetTrack(trackId)
	if err != nil {
		return nil, 0, 0, "", 0, err
	}

	if !track.Public && track.UserId != userId {
		return nil, 0, 0, "", 0, errors.New("permission denied")
	}

	var buf bytes.Buffer
	var content_type string
	var full_size int64

	folderPath := strings.Replace(viper.GetString("media.tracks_folder"), "{user_id}", strconv.Itoa(userId), 1)
	filePath := filepath.Join(folderPath, track.TrackFilename)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, 0, "", 0, err
	}
	fileInfo, _ := os.Stat(filePath)
	defer file.Close()

	mtype, err := mimetype.DetectReader(file)
	if err != nil {
		return nil, 0, 0, "", 0, err
	}
	content_type = mtype.String()
	full_size = fileInfo.Size()

	var rangeStart int64
	var rangeEnd int64
	if bytesRange[1] == nil {
		rangeEnd = full_size
	} else {
		rangeEnd = *bytesRange[1]
	}

	if bytesRange[0] == nil {
		rangeStart = 0
	} else {
		fmt.Println(bytesRange[0])
		rangeStart = *bytesRange[0]
	}

	N := rangeEnd - rangeStart

	file.Seek(rangeStart, io.SeekStart)
	_, err = io.CopyN(&buf, file, N)
	if err != nil {
		return nil, 0, 0, "", 0, err
	}

	return &buf, rangeStart, rangeEnd, content_type, full_size, nil
}
