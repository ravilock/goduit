package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/articlePublisher/assemblers"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
	"github.com/ravilock/goduit/internal/articlePublisher/requests"
	profileManagerAssembler "github.com/ravilock/goduit/internal/profileManager/assemblers"
)

type articleWriter interface {
	WriteArticle(ctx context.Context, article *models.Article) error
}

type writeArticleHandler struct {
	service        articleWriter
	profileManager profileGetter
}

func (h *writeArticleHandler) WriteArticle(c echo.Context) error {
	authorUsername := c.Request().Header.Get("Goduit-Client-Username")
	request := new(requests.WriteArticle)
	if err := c.Bind(request); err != nil {
		return api.CouldNotUnmarshalBodyError
	}

	if err := request.Validate(); err != nil {
		return err
	}

	article := request.Model(authorUsername)

	ctx := c.Request().Context()

	err := h.service.WriteArticle(ctx, article)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.ConflictErrorCode:
				return api.ConfictError
			}
		}
		return err
	}

	authorProfile, err := h.profileManager.GetProfileByUsername(ctx, authorUsername)
	if err != nil {
		return api.UserNotFound(authorUsername)
	}

	profileResponse, err := profileManagerAssembler.ProfileResponse(authorProfile, false)
	if err != nil {
		return err
	}

	response := assemblers.ArticleResponse(article, profileResponse)
	return c.JSON(http.StatusCreated, response)
}