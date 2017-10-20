package cache

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"
	"time"
)

const (
	indexJSON      = "index.json"
	pagesDirectory = "pages"
	pageSuffix     = ".md"
	zipPath        = "/tldr.zip"
)

// Repository keeps a copy of the data from the remote location on the local
// filesystem. It implements the tldr.Repository to provide quick access
// to the requested markdown.
type Repository struct {
	directory string
	remote    string
	ttl       time.Duration
}

// NewRepository returns a new cache repository. The data is loaded from the
// remote if missing or stale.
func NewRepository(remote string, ttl time.Duration) (*Repository, error) {
	dir, err := cacheDir()
	if err != nil {
		return nil, fmt.Errorf("ERROR: getting cache directory: %s", err)
	}

	repo := &Repository{directory: dir, remote: remote, ttl: ttl}

	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		fmt.Println("creating cache dir ...")
		err = repo.makeCacheDir()
		if err != nil {
			return nil, fmt.Errorf("ERROR: creating cache directory: %s", err)
		}
		fmt.Println("fetch pages ...")
		err = repo.loadFromRemote()
		if err != nil {
			return nil, fmt.Errorf("ERROR: loading data from remote: %s", err)
		}
	} else if err != nil || info.ModTime().Before(time.Now().Add(-ttl)) {
		err = repo.Reload()
		if err != nil {
			return nil, fmt.Errorf("ERROR: reloading cache: %s", err)
		}
	}

	return repo, nil
}

// AvailablePlatforms returns all the availale platforms found in cache.
func (r *Repository) AvailablePlatforms() ([]string, error) {
	var platforms []string
	available, err := ioutil.ReadDir(path.Join(r.directory, pagesDirectory))
	if err != nil {
		return nil, err
	}

	for _, f := range available {
		platform := f.Name()
		if platform != indexJSON {
			platforms = append(platforms, platform)
		}
	}
	return platforms, nil
}

// Markdown pulls the markdown from the page in cache.
func (r *Repository) Markdown(platform, page string) (io.ReadCloser, error) {
	return os.Open(path.Join(r.directory, pagesDirectory, platform, page+pageSuffix))
}

// Pages returns all the pages for the given platform.
func (r *Repository) Pages(platform string) ([]string, error) {
	dir := path.Join(r.directory, pagesDirectory, platform)
	pages, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("ERROR: reading directory '%s': %s", dir, err)
	}

	names := make([]string, len(pages))
	for i, page := range pages {
		name := page.Name()
		names[i] = name[:len(name)-3]
	}
	return names, nil
}

// Reload removes the cache directory, recreates it, and saves the data from the remote
// to the local filesystem.
func (r *Repository) Reload() error {
	fmt.Println("removing old pages ...")
	err := os.RemoveAll(r.directory)
	if err != nil {
		return fmt.Errorf("ERROR: removing cache directory: %s", err)
	}

	fmt.Println("creating new cache dir ...")
	err = r.makeCacheDir()
	if err != nil {
		return fmt.Errorf("ERROR: creating cache directory: %s", err)
	}

	fmt.Println("fetch new pages ...")
	err = r.loadFromRemote()
	if err != nil {
		return fmt.Errorf("ERROR: loading data from remote: %s", err)
	}
	return nil
}

func (r *Repository) loadFromRemote() error {
	cache, err := os.Create(r.directory + zipPath)
	if err != nil {
		return fmt.Errorf("ERROR: creating cache: %s", err)
	}
	defer cache.Close()

	resp, err := http.Get(r.remote)
	if err != nil {
		return fmt.Errorf("ERROR: getting response from remote: %s", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(cache, resp.Body)
	if err != nil {
		return fmt.Errorf("ERROR: copying response body to cache: %s", err)
	}

	fmt.Println("unzipping pages ...")
	err = r.unzip()
	if err != nil {
		return fmt.Errorf("ERROR: unzipping pages: %s", err)
	}

	err = os.Remove(r.directory + zipPath)
	if err != nil {
		return fmt.Errorf("ERROR: removing zip: %s", err)
	}
	return nil
}

func (r *Repository) makeCacheDir() error {
	return os.MkdirAll(r.directory, 0755)
}

func (r *Repository) unzip() error {
	reader, err := zip.OpenReader(r.directory + zipPath)
	if err != nil {
		return fmt.Errorf("err opening zip: %s", err)
	}
	defer reader.Close()

	for _, f := range reader.File {
		zipFile, err := f.Open()
		if err != nil {
			return fmt.Errorf("err opening file '%s': %s", f.Name, err)
		}
		defer zipFile.Close()

		filepath := path.Join(r.directory, f.Name)
		if f.FileInfo().IsDir() {
			err := os.MkdirAll(filepath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("err making directory '%s': %s", filepath, err)
			}
			continue
		}

		var dirPath string
		if lastIndex := strings.LastIndex(filepath, string(os.PathSeparator)); lastIndex > -1 {
			dirPath = filepath[:lastIndex]
		}

		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("err making directories for '%s': %s", filepath, err)
		}

		file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("err opening file '%s': %s", filepath, err)
		}
		defer file.Close()

		_, err = io.Copy(file, zipFile)
		if err != nil {
			return fmt.Errorf("err copying file '%s': %s", file.Name(), err)
		}
	}
	return nil
}

func cacheDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("ERROR: getting current user: %s", err)
	}
	if usr.HomeDir == "" {
		return "", fmt.Errorf("ERROR: loading current user's home directory")
	}
	return path.Join(usr.HomeDir, ".tldr"), nil
}
