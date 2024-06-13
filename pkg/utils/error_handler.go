package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type ErrResponse struct {
	Message  string      `json:"message"`
	Metadata interface{} `json:"metadata,omitempty"`
}

func HttpErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	if m, ok := err.(*ValidationError); ok {
		c.JSON(http.StatusBadRequest, m)
	} else if m, ok := err.(*echo.HTTPError); ok {
		switch mType := m.Message.(type) {
		case string:
			c.JSON(m.Code, ErrResponse{Message: mType})
		case error:
			c.JSON(m.Code, ErrResponse{Message: mType.Error()})
		}
	} else {
		log.Err(err).Msg("http error")
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, ErrResponse{
				Message: "Resource not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrResponse{
			Message: http.StatusText(http.StatusInternalServerError),
		})
	}
}
