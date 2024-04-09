package repository

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/jmoiron/sqlx"
	"github.com/vekshinnikita/golang_music"
	"github.com/vekshinnikita/golang_music/pkg/tools"
)

type TrackPostgres struct {
	db *sqlx.DB
}

func NewTrackPostgres(db *sqlx.DB) *TrackPostgres {
	return &TrackPostgres{db}
}

func (r *TrackPostgres) AddTrack(userId int, track golang_music.AddTrackInput) (int, error) {
	var id int
	setKeys := []string{"title", "author", "user_id"}
	setValues := []string{"$1", "$2", "$3"}
	args := []interface{}{track.Title, track.Author, userId}

	if track.Description != nil {
		setKeys = append(setKeys, "description")
		setValues = append(setValues, fmt.Sprintf("$%d", len(setValues)+1))
		args = append(args, *track.Description)
	}

	if track.Public != nil {
		setKeys = append(setKeys, "public")
		setValues = append(setValues, fmt.Sprintf("$%d", len(setValues)+1))
		args = append(args, *track.Public)
	}

	setQueryValues := strings.Join(setValues, ", ")
	setQueryKeys := strings.Join(setKeys, ", ")

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id", tracksTable, setQueryKeys, setQueryValues)

	row := r.db.QueryRow(query, args...)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *TrackPostgres) UpdateTrackFilename(filename string, trackId int, userId int) error {
	query := fmt.Sprintf("UPDATE %s SET track_file_name=$1, updated_at=CURRENT_TIMESTAMP WHERE user_id=$2 AND id=$3", tracksTable)

	_, err := r.db.Exec(query, filename, userId, trackId)
	if err != nil {
		return err
	}

	return nil
}

func (r *TrackPostgres) UpdatePosterFilename(filename string, trackId int, userId int) error {
	query := fmt.Sprintf("UPDATE %s SET poster_file_name=$1, updated_at=CURRENT_TIMESTAMP WHERE user_id=$2 AND id=$3", tracksTable)

	_, err := r.db.Exec(query, filename, userId, trackId)
	if err != nil {
		return err
	}

	return nil
}

// func (r *TrackPostgres) GetTrackFilename(trackId int, userId int) (string, error) {
// 	var track_file_name string
// 	query := fmt.Sprintf("SELECT track_file_name FROM %s WHERE user_id=$1 AND id=$2", tracksTable)
// 	err := r.db.Get(&track_file_name, query, userId, trackId)

// 	if err != nil {
// 		return "", err
// 	}
// 	return track_file_name, nil
// }

func (r *TrackPostgres) UpdateTrack(userId int, input golang_music.UpdateTrackInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argsCounter := 1

	inputMap := structs.Map(input)
	delete(inputMap, "TrackId")

	for key, value := range inputMap {
		field, ok := reflect.TypeOf(&input).Elem().FieldByName(key)
		if !ok {
			message := fmt.Sprintf("field '%s' not found", key)
			return errors.New(message)
		}
		fmt.Println(value)
		if value != "" {
			setValues = append(setValues, fmt.Sprintf("%s=$%d", tools.GetStructTag(field, "json"), argsCounter))
			args = append(args, value)
			argsCounter++
		}
	}

	args = append(args, userId, input.TrackId)

	setQueries := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE user_id=$%d AND id=$%d", tracksTable, setQueries, argsCounter, argsCounter+1)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

// func (r *TrackPostgres) GetUserIdByTrackId(trackId int) (int, error) {
// 	fmt.Println(trackId)
// 	var user_id int
// 	query := fmt.Sprintf("SELECT user_id FROM %s WHERE id=$1", tracksTable)

// 	row := r.db.QueryRow(query, trackId)
// 	if err := row.Scan(&user_id); err != nil {
// 		return 0, err
// 	}

// 	return user_id, nil
// }

func (r *TrackPostgres) DeleteTrack(trackId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", tracksTable)

	_, err := r.db.Exec(query, trackId)
	if err != nil {
		return err
	}

	return nil
}

func (r *TrackPostgres) GetTrack(trackId int) (*golang_music.Track, error) {
	var track golang_music.Track
	query := fmt.Sprintf("SELECT id, user_id, title, description, author, track_file_name, poster_file_name, created_at, updated_at FROM %s WHERE id=$1", tracksTable)

	err := r.db.Get(&track, query, trackId)
	if err != nil {
		return nil, err
	}

	return &track, nil
}
