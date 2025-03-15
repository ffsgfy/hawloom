package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/db"
	"github.com/ffsgfy/hawloom/internal/ui"
)

type userParams struct {
	Name string `param:"name"`
}

func HandleUser(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		var params userParams
		if err := c.Bind(&params); err != nil {
			return err
		}

		sc := s.Ctx(c.Request().Context())
		account, err := sc.FindAccount(nil, &params.Name)
		if err != nil {
			return err
		}

		self := false
		if authToken, _ := api.GetValidAuthToken(sc.Ctx); authToken != nil {
			if authToken.AccountID == account.ID {
				self = true
			}
		}

		var docs []*db.Doc
		if self {
			if docs, err = sc.Queries.FindDocList(sc.Ctx, account.ID); err != nil {
				return err
			}
		} else {
			if docs, err = sc.Queries.FindPublicDocList(sc.Ctx, account.ID); err != nil {
				return err
			}
		}

		docRows := make([]*ui.DocRow, 0, len(docs))
		for _, doc := range docs {
			docRows = append(docRows, &ui.DocRow{
				ID:          doc.ID,
				Title:       doc.Title,
				Description: doc.Description,
			})
		}

		content, err := ui.Render(c.Request().Context(), ui.UserPage(account, self, docRows))
		if err != nil {
			return err
		}
		return c.HTML(http.StatusOK, content)
	}
}
