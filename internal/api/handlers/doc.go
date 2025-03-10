package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/db"
	"github.com/ffsgfy/hawloom/internal/ui"
)

func HandleNewDoc(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		if _, err := api.GetValidAuthToken(c.Request().Context()); err != nil {
			return handleRedirect(c, "/")
		}

		content, err := ui.Render(c.Request().Context(), ui.NewDocPage())
		if err != nil {
			return err
		}
		return c.HTML(http.StatusOK, content)
	}
}

type newDocParams struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	Content     string `form:"content"`
	VDuration   int32  `form:"vduration"`
	VMode       string `form:"vmode"`
	Public      bool   `form:"public"`
	Majority    bool   `form:"majority"`
}

func HandleNewDocPost(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		var params newDocParams
		if err := c.Bind(&params); err != nil {
			return err
		}

		var flags api.DocFlags
		if params.Public {
			flags |= api.DocFlagPublic
		}
		if params.Majority {
			flags |= api.DocFlagMajority
		}
		if params.VMode == "approval" {
			flags |= api.DocFlagApproval
		}

		sc := s.Ctx(c.Request().Context())
		doc, _, err := sc.CreateDoc(&api.CreateDocParams{
			Title:        params.Title,
			Description:  params.Description,
			Content:      params.Content,
			Flags:        flags,
			VordDuration: params.VDuration,
		})
		if err != nil {
			return err
		}

		return handleRedirect(c, "/doc/"+doc.ID.String())
	}
}

type docParams struct {
	DocID   uuid.UUID  `param:"doc"`
	VerID   *uuid.UUID `query:"ver"`
	VordNum *int32     `param:"vord"`
}

func HandleDoc(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		var params docParams
		if err := c.Bind(&params); err != nil {
			return err
		}

		sc := s.Ctx(c.Request().Context())
		var doc *db.Doc
		var ver *db.Ver // current
		var err error

		if params.VordNum == nil {
			if row, e := sc.Queries.FindCurrentVer(sc.Ctx, params.DocID); e == nil {
				doc = &row.Doc
				ver = &row.Ver
			} else {
				err = e
			}
		} else {
			if row, e := sc.Queries.FindWinningVer(sc.Ctx, params.DocID, *params.VordNum-1); e == nil {
				doc = &row.Doc
				ver = &row.Ver
			} else {
				err = e
			}
		}
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return api.ErrDocNotFound
			}
			return err
		}

		content, err := ui.Render(c.Request().Context(), ui.DocPage(doc, ver, params.VerID))
		if err != nil {
			return err
		}
		return c.HTML(http.StatusOK, content)
	}
}
