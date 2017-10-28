package tldr

import "io"

// Repository is used to abstract where the pages are stored.
type Repository interface {
	AvailablePlatforms() ([]string, error)
	Markdown(platform, page string) (io.ReadCloser, error)
	Pages() ([]string, error)
}
