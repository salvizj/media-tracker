package types

import "media_tracker/internal/models"

type LayoutTmplData struct {
	Title           string
	Message         string
	ContentTemplate string
	Movies          []models.Movie
	TvShows         []models.TVShow
	ManhwaAndManga  []models.ManhwaAndManga
}
