package types

import "media_tracker/internal/models"

type LayoutTmplData struct {
	Title           string
	ContentTemplate string
	Movies          []models.Movie
	TVShows         []models.TVShow
	ManhwaAndManga  []models.ManhwaAndManga
}
