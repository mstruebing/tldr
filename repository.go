package tldr

import "io"

type Repository interface {
	Markdown(platform, page string) (io.ReadCloser, error)
	Pages(platform string) ([]string, error)
}
