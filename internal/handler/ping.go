package handler

import (
	"database/sql"
	"net/http"

	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/validate"
	_ "github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
)

func (s Handler) HandlePing(echoCtx echo.Context) error {
	if echoCtx.Request().Method != http.MethodGet {
		return validate.AddError(s.ctx, echoCtx, validate.WrongTypeRequest, http.StatusBadRequest, 0)
	}

	flags := config.ConfigFromContext(s.ctx)
	db, err := sql.Open("pgx", flags.DataBaseDsn)
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, "", http.StatusInternalServerError, 0)
	}
	defer db.Close()

	return validate.AddMessageOK(s.ctx, echoCtx, "", http.StatusOK, 0)
}
