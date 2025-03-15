package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/db"
	"github.com/ffsgfy/hawloom/internal/ui"
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
		var vote int = -1
		var err error

		if authToken, _ := api.GetValidAuthToken(sc.Ctx); authToken != nil {
			vote = 0
			if row, err := sc.Queries.FindVerWithVote(
				sc.Ctx, params.VerID, authToken.AccountID,
			); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return api.ErrVerNotFound
				}
				return err
			} else {
				ver = &row.Ver
				if row.HasVote {
					vote = 1
				}
			}
		} else {
			if ver, err = sc.Queries.FindVer(sc.Ctx, params.VerID); err != nil {
				return err
			}
		}

		content, err := ui.Render(sc.Ctx, ui.VerFragment(ver, vote))
		if err != nil {
			return err
		}
		return c.HTML(http.StatusOK, content)
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

		content, err := ui.Render(sc.Ctx, ui.VerVoteUnvoteButton(params.VerID, vote))
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
