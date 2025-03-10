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

func prepareVerRows(vers []*db.FindVerListRow, votes []uuid.UUID) []ui.VerRow {
	voteMap := map[uuid.UUID]bool{}
	for _, vote := range votes {
		voteMap[vote] = true
	}

	rows := make([]ui.VerRow, 0, len(vers))
	for _, ver := range vers {
		rows = append(rows, ui.VerRow{
			ID:      ver.ID,
			Votes:   strconv.FormatInt(int64(ver.Votes), 10),
			Author:  ver.Author,
			Summary: ver.Summary,
			HasVote: voteMap[ver.ID],
		})
	}
	return rows
}

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
		c.Response().Header().Set(HXTriggerAfterSwap, "update-view-alternate")
		return c.HTML(http.StatusOK, content)
	}
}
