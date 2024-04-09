package handler

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vekshinnikita/golang_music"
	"github.com/vekshinnikita/golang_music/pkg/tools"
)

func (h *Handler) AddTrack(c *gin.Context) {
	var input golang_music.AddTrackInput

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	trackId, err := h.services.Track.AddTrack(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, &map[string]interface{}{
		"id": trackId,
	})
}

func (h *Handler) UploadTrack(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	trackIdForm := c.PostForm("trackId")
	trackId, err := strconv.Atoi(trackIdForm)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	contentRangeHeader := strings.Split(c.Request.Header.Get("Content-Range"), " ")[1]

	rangeAndSize := strings.Split(contentRangeHeader, "/")

	rangeBytes, err := tools.Map(strings.Split(rangeAndSize[0], "-"), func(v string) (int64, error) {
		value, _ := strconv.ParseInt(v, 10, 64)
		return value, nil
	})

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Missing file size in Content-Range header")
		return
	}

	fileSize, err := strconv.ParseInt(rangeAndSize[1], 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Missing file size in Content-Range header")
		return
	}

	err = h.services.Track.UploadTrack(userId, trackId, file, rangeBytes, fileSize)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if rangeBytes[1] >= fileSize-1 {
		c.JSON(http.StatusOK, &map[string]any{
			"status": "uploaded",
		})
		return
	}

	percentage := (float64(rangeBytes[1]) + 1) / (float64(fileSize) / 100)

	c.JSON(http.StatusOK, &map[string]any{
		"status":     "uploading",
		"percentage": math.Ceil(percentage),
	})
}

func (h *Handler) UploadPoster(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if file.Size > 10<<20 {
		newErrorResponse(c, http.StatusBadRequest, "file size grater then 10MB")
		return
	}

	trackIdForm := c.PostForm("trackId")
	trackId, err := strconv.Atoi(trackIdForm)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Track.UploadPoster(userId, trackId, file)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) UpdateTrack(c *gin.Context) {
	var input golang_music.UpdateTrackInput

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Track.UpdateTrack(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteTrack(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	trackId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Track.DeleteTrack(userId, trackId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetPoster(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	trackId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	buf, err := h.services.Track.GetPoster(userId, trackId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	contentType := http.DetectContentType(buf.Bytes())

	c.Writer.Header().Add("Content-Type", contentType)
	c.Writer.Header().Add("Content-Length", strconv.Itoa(len(buf.Bytes())))

	c.Writer.Write(buf.Bytes())
	c.Status(http.StatusOK)
}

func (h *Handler) GetTrackInfo(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	trackId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	track, err := h.services.Track.GetTrackInfo(userId, trackId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	domain := tools.GetDomain(c)

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":          track.TrackId,
		"user_id":     track.UserId,
		"author":      track.Author,
		"title":       track.Title,
		"description": track.Description,
		"public":      track.Public,
		"track_url":   fmt.Sprintf("%s/api/track/%d/streaming", domain, track.TrackId),
		"poster_url":  fmt.Sprintf("%s/api/track/%d/poster", domain, track.TrackId),
		"created_at":  track.CreatedAt,
		"updated_at":  track.UpdatedAt,
	})
}

func (h *Handler) StreamingTrack(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	trackId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var bytesRange []*int64 = nil

	bytesRangeForm := c.GetHeader("Range")
	if bytesRangeForm != "" {
		bytesRangeForm = strings.Split(bytesRangeForm, "=")[1]
		bytesRange, err = tools.Map(strings.Split(bytesRangeForm, "-"), func(v string) (*int64, error) {
			fmt.Println(v)
			fmt.Println(v != "")
			if v != "" {
				value, _ := strconv.ParseInt(v, 10, 64)
				return &value, nil
			}
			return nil, nil

		})
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

	}

	buf, rangeStart, rangeEnd, contentType, fullSize, err := h.services.Track.StreamingTrack(userId, trackId, bytesRange)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	contentLen := len(buf.Bytes())
	if fullSize == int64(contentLen) {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusPartialContent)
	}
	c.Writer.Header().Add("Content-Type", contentType)
	c.Writer.Header().Add("Content-Length", strconv.Itoa(contentLen))
	c.Writer.Header().Add("Accept-Ranges", "bytes")
	c.Writer.Header().Add("Content-Range", fmt.Sprintf("bytes %d-%d/%d", rangeStart, rangeEnd, fullSize))
	c.Writer.Write(buf.Bytes())
}
