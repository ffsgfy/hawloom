package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/ui"
)

func HandleMain(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		sc := s.Ctx(c.Request().Context())
		docs, err := sc.Queries.FindAllPublicDocList(sc.Ctx)
		if err != nil {
			return err
		}

		docRows := make([]*ui.DocRow, 0, len(docs))
		for _, doc := range docs {
			docRows = append(docRows, &ui.DocRow{
				ID:          doc.ID,
				Title:       doc.Title,
				Description: doc.Description,
				Author:      doc.Author,
			})
		}

		content, err := ui.Render(c.Request().Context(), ui.MainPage(docRows))
		if err != nil {
			return err
		}
		return c.HTML(http.StatusOK, content)
	}
}
