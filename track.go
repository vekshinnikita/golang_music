package golang_music

type AddTrackInput struct {
	Author      string  `json:"author" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	Public      *bool   `json:"public"`
}

type UpdateTrackInput struct {
	TrackId     int    `json:"trackId" binding:"required"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

type Track struct {
	TrackId        int    `db:"id" binding:"required"`
	UserId         int    `db:"user_id" binding:"required"`
	Author         string `db:"author" binding:"required"`
	Title          string `db:"title" binding:"required"`
	Description    string `db:"description"`
	Public         bool   `db:"public" binding:"required"`
	TrackFilename  string `db:"track_file_name" binding:"required"`
	PosterFilename string `db:"poster_file_name" binding:"required"`
	CreatedAt      string `db:"created_at" binding:"required"`
	UpdatedAt      string `db:"updated_at" binding:"required"`
}
