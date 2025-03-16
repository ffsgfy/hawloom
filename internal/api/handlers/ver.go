package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/db"
	"github.com/ffsgfy/hawloom/internal/ui"
	"github.com/ffsgfy/hawloom/internal/utils"
)

type verParams struct {
	VerID uuid.UUID `param:"ver"`
}

func HandleVer(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		var params verParams
		if err := c.Bind(&params); err != nil {
			return err
		}

		sc := s.Ctx(c.Request().Context())
		var ver *db.Ver
		var author string
		var hasVote, canVote bool

		if authToken, _ := api.GetValidAuthToken(sc.Ctx); authToken != nil {
			if row, err := sc.Queries.FindVerWithVote(
				sc.Ctx, params.VerID, authToken.AccountID,
			); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return api.ErrVerNotFound
				}
				return err
			} else {
				ver = &row.Ver
				author = row.VerAuthor
				hasVote = row.VerVoteExists
				if utils.TestFlags(row.DocFlags, api.DocFlagApproval) || !row.DocVoteExists {
					canVote = true
				}
			}
		} else {
			if row, err := sc.Queries.FindVer(sc.Ctx, params.VerID); err != nil {
				return err
			} else {
				ver = &row.Ver
				author = row.Author
			}
		}

		content, err := ui.Render(sc.Ctx, ui.VerFragment(ver, author, hasVote, canVote))
		if err != nil {
			return err
		}
		return c.HTML(http.StatusOK, content)
	}
}

func HandleVerDelete(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		var params verParams
		if err := c.Bind(&params); err != nil {
			return err
		}

		sc := s.Ctx(c.Request().Context())
		if docID, err := sc.DeleteVer(params.VerID); err != nil {
			return err
		} else {
			return handleRedirect(c, fmt.Sprintf("/doc/%v", docID))
		}
	}
}

func HandleVerVoteUnvote(s *api.State, vote bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		var params verParams
		if err := c.Bind(&params); err != nil {
			return err
		}

		sc := s.Ctx(c.Request().Context())
		var err error

		if vote {
			err = sc.CreateVote(params.VerID)
		} else {
			err = sc.DeleteVote(params.VerID)
		}
		if err != nil {
			return err
		}

		content, err := ui.Render(sc.Ctx, ui.VerVoteUnvoteButton(params.VerID, vote, true))
		if err != nil {
			return err
		}
		return c.HTML(http.StatusOK, content)
	}
}

func prepareVerRows(vers []*db.FindVerListRow, votes []uuid.UUID) []*ui.VerRow {
	voteMap := map[uuid.UUID]bool{}
	for _, vote := range votes {
		voteMap[vote] = true
	}

	rows := make([]*ui.VerRow, 0, len(vers))
	for _, ver := range vers {
		rows = append(rows, &ui.VerRow{
			ID:      ver.ID,
			Author:  ver.Author,
			Summary: ver.Summary,
			Votes:   strconv.FormatInt(int64(ver.Votes), 10),
			HasVote: voteMap[ver.ID],
		})
	}
	return rows
}

type verListParams struct {
	DocID   uuid.UUID `query:"doc-id"`
	VordNum int32     `query:"vord-num"`
}

func HandleVerList(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		var params verListParams
		if err := c.Bind(&params); err != nil {
			return err
		}

		sc := s.Ctx(c.Request().Context())
		vers, err := sc.Queries.FindVerList(sc.Ctx, params.DocID, params.VordNum)
		if err != nil {
			return err
		}

		var votes []uuid.UUID
		if len(vers) > 0 {
			if authToken, _ := api.GetValidAuthToken(sc.Ctx); authToken != nil {
				if votes, err = sc.Queries.FindVoteList(sc.Ctx, &db.FindVoteListParams{
					Doc:     params.DocID,
					VordNum: params.VordNum,
					Account: authToken.AccountID,
				}); err != nil {
					return err
				}
			}
		}

		verRows := prepareVerRows(vers, votes)
		content, err := ui.Render(sc.Ctx, ui.VerList(verRows))
		if err != nil {
			return err
		}
		return c.HTML(http.StatusOK, content)
	}
}

type newVerParams struct {
	DocID uuid.UUID `query:"doc"`
}

func HandleNewVer(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		var params newVerParams
		if err := c.Bind(&params); err != nil {
			return err
		}

		sc := s.Ctx(c.Request().Context())
		row, err := sc.Queries.FindCurrentVer(sc.Ctx, params.DocID)
		if err != nil {
			return err
		}

		content, err := ui.Render(sc.Ctx, ui.NewVerPage(&row.Doc, row.DocAuthor, &row.Ver, row.VerAuthor))
		if err != nil {
			return err
		}
		return c.HTML(http.StatusOK, content)
	}
}

type newVerPostParams struct {
	DocID   uuid.UUID `form:"doc-id"`
	Content string    `form:"content"`
	Summary string    `form:"summary"`
}

func HandleNewVerPost(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		var params newVerPostParams
		if err := c.Bind(&params); err != nil {
			return err
		}

		sc := s.Ctx(c.Request().Context())
		if ver, err := sc.CreateVer(&api.CreateVerParams{
			DocID:   params.DocID,
			Summary: params.Summary,
			Content: params.Content,
		}); err != nil {
			return err
		} else {
			return handleRedirect(c, fmt.Sprintf("/doc/%v?ver=%v", ver.Doc, ver.ID))
		}
	}
}
