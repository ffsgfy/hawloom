package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/db"
	"github.com/ffsgfy/hawloom/internal/ui"
)

func HandleNewDoc(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		content, err := ui.Render(c.Request().Context(), ui.NewDocPage())
		if err != nil {
			return err
		}
		return c.HTML(http.StatusOK, content)
	}
}

type newDocPostParams struct {
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
		var params newDocPostParams
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
		if doc, _, err := sc.CreateDoc(&api.CreateDocParams{
			Title:        params.Title,
			Description:  params.Description,
			Content:      params.Content,
			Flags:        flags,
			VordDuration: params.VDuration,
		}); err != nil {
			return err
		} else {
			return handleRedirect(c, "/doc/"+doc.ID.String())
		}
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
		var ver *db.Ver // current ver
		var vord *db.Vord
		var docAuthor, verAuthor string

		if params.VordNum == nil {
			if row, err := sc.Queries.FindCurrentVer(sc.Ctx, params.DocID); err == nil {
				doc = &row.Doc
				ver = &row.Ver
				vord = &row.Vord
				docAuthor = row.DocAuthor
				verAuthor = row.VerAuthor
			} else if errors.Is(err, sql.ErrNoRows) {
				return api.ErrDocNotFound
			} else {
				return err
			}
		} else {
			*params.VordNum = max(*params.VordNum, 0)
			if row, err := sc.Queries.FindWinningVer(sc.Ctx, &db.FindWinningVerParams{
				Doc:         params.DocID,
				VordNum:     max(*params.VordNum-1, 0),
				VordNumJoin: *params.VordNum,
			}); err == nil {
				doc = &row.Doc
				ver = &row.Ver
				vord = &row.Vord
				docAuthor = row.DocAuthor
				verAuthor = row.VerAuthor
			} else if errors.Is(err, sql.ErrNoRows) {
				return handleRedirect(c, fmt.Sprintf("/doc/%v", params.DocID))
			} else {
				return err
			}
		}

		content, err := ui.Render(c.Request().Context(), ui.DocPage(doc, docAuthor, ver, verAuthor, params.VerID, vord))
		if err != nil {
			return err
		}
		return c.HTML(http.StatusOK, content)
	}
}

type docTimerParams struct {
	DocID uuid.UUID `param:"doc"`
	Start int64     `query:"start"`
}

func HandleDocTimer(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		var params docTimerParams
		if err := c.Bind(&params); err != nil {
			return err
		}

		sc := s.Ctx(c.Request().Context())
		if vord, err := sc.Queries.FindVord(sc.Ctx, params.DocID, -1); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return api.ErrNoVordExists
			}
			return err
		} else {
			if vord.StartAt.Unix() > params.Start {
				return handleRedirect(c, fmt.Sprintf("/doc/%v", vord.Doc))
			}

			content, err := ui.Render(c.Request().Context(), ui.DocVordTimer(vord))
			if err != nil {
				return err
			}
			return c.HTML(http.StatusOK, content)
		}
	}
}
