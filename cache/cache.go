package cache

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	indexJSON      = "index.json"
	pagesDirectory = "pages"
	historyPath    = "/history"
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

// HistoryRecord represent the search history of certain page
type HistoryRecord struct {
	page  string
	count int
}

func (h HistoryRecord) String() string {
	return fmt.Sprintf("%s %d", h.page, h.count)
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
		if repo.isReachable() {
			err = repo.Reload()
			if err != nil {
				return nil, fmt.Errorf("ERROR: reloading cache: %s", err)
			}
		} else {
			fmt.Println("INFO: remote is not reachable, reload skipped")
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
func (r *Repository) Pages() ([]string, error) {
	dir := path.Join(r.directory, pagesDirectory)

	pages := []os.FileInfo{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
			pages = append(pages, f)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("ERROR: can't read pages")
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
	err := os.RemoveAll(r.directory)
	if err != nil {
		return fmt.Errorf("ERROR: removing cache directory: %s", err)
	}

	err = r.makeCacheDir()
	if err != nil {
		return fmt.Errorf("ERROR: creating cache directory: %s", err)
	}

	err = r.loadFromRemote()
	if err != nil {
		return fmt.Errorf("ERROR: loading data from remote: %s", err)
	}
	return nil
}

func (r *Repository) copyZipFile(f *zip.File) error {
	zipFile, err := f.Open()
	if err != nil {
		return fmt.Errorf("ERROR: opening file '%s': %s", f.Name, err)
	}
	defer zipFile.Close()

	filepath := path.Join(r.directory, f.Name)
	if f.FileInfo().IsDir() {
		err := os.MkdirAll(filepath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("ERROR: making directory '%s': %s", filepath, err)
		}
		return nil
	}

	var dirPath string
	if lastIndex := strings.LastIndex(filepath, string(os.PathSeparator)); lastIndex > -1 {
		dirPath = filepath[:lastIndex]
	}

	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("ERROR: making directories for '%s': %s", filepath, err)
	}

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return fmt.Errorf("ERROR: opening file '%s': %s", filepath, err)
	}
	defer file.Close()

	_, err = io.Copy(file, zipFile)
	if err != nil {
		return fmt.Errorf("ERROR: copying file '%s': %s", file.Name(), err)
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
	if err := os.MkdirAll(r.directory, 0755); err != nil {
		return fmt.Errorf("ERROR: creating directory %s: %s", r.directory, err)
	}
	// touching the history file.
	historyFile := path.Join(r.directory, historyPath)
	return touchFile(historyFile)
}

func touchFile(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("ERROR: creating file %s: %s", fileName, err)
	}
	defer file.Close()
	return nil
}

func (r *Repository) unzip() error {
	reader, err := zip.OpenReader(r.directory + zipPath)
	if err != nil {
		return fmt.Errorf("ERROR: opening zip: %s", err)
	}
	defer reader.Close()

	for _, f := range reader.File {
		err = r.copyZipFile(f)
		if err != nil {
			return fmt.Errorf("ERROR: copying zip file: %s", err)
		}
	}
	return nil
}

func cacheDir() (string, error) {
	XDG_CACHE_HOME := os.Getenv("XDG_CACHE_HOME")

	// Use the XDG_CACHE_HOME environment variable if possible
	if XDG_CACHE_HOME != "" {
		return path.Join(XDG_CACHE_HOME, "tldr"), nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("ERROR: getting current user: %s", err)
	}

	homeDir, err := filepath.Abs(usr.HomeDir)

	if usr.HomeDir == "" && err != nil {
		return "", fmt.Errorf("ERROR: loading current user's home directory: %s", err)
	}

	XDG_CACHE_HOME_DEFAULT := ".cache/"

	return path.Join(homeDir, XDG_CACHE_HOME_DEFAULT, "tldr"), nil
}

func (r Repository) isReachable() bool {
	u, err := url.Parse(r.remote)
	if err != nil {
		return false
	}

	seconds := 5
	timeout := time.Duration(seconds) * time.Second

	_, err = net.DialTimeout("tcp", u.Hostname()+":"+u.Port(), timeout)
	return err == nil
}

func (r Repository) RecordHistory(page string) error {
	records, err := r.LoadHistory()
	if err != nil {
		return fmt.Errorf("ERROR: loading history failed %s", err)
	}

	newRecord := HistoryRecord{
		page:  page,
		count: 1,
	}

	foundIdx := -1
	for idx, r := range *records {
		if r.page == page {
			newRecord.count = r.count + 1
			foundIdx = idx
			break
		}
	}

	if foundIdx != -1 { //found in history, we want to put the last search at the end of the history.
		newRecords := append((*records)[:foundIdx], (*records)[foundIdx+1:]...)
		records = &newRecords
	}

	newRecords := append(*records, newRecord)
	return r.saveHistory(&newRecords)
}

func (r Repository) saveHistory(history *[]HistoryRecord) error {
	hisFile := path.Join(r.directory, historyPath)
	inFile, err := os.Create(hisFile)
	if err != nil {
		return fmt.Errorf("ERROR: opening history file %s", hisFile)
	}
	defer inFile.Close()

	for _, his := range *history {
		fmt.Fprintln(inFile, fmt.Sprintf("%s,%d", his.page, his.count))
	}
	return nil
}

func (r Repository) LoadHistory() (*[]HistoryRecord, error) {
	// read the history file line by line, into a map.
	history := path.Join(r.directory, historyPath)
	//if it is not exist, touch it.
	_, err := os.Stat(history)

	if os.IsNotExist(err) {
		if err := touchFile(history); err != nil {
			return nil, fmt.Errorf("ERROR: cannot create the history file %s, %s", history, err)
		}
	}

	inFile, err := os.Open(history)
	if err != nil {
		return nil, fmt.Errorf("ERROR: opening history file %s", history)

	}

	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	historyRecords := make([]HistoryRecord, 0, 10)

	for scanner.Scan() {
		line := scanner.Text()
		lineParts := strings.Split(line, ",")
		count, err := strconv.Atoi(lineParts[1])

		if err != nil {
			return nil, err
		}

		historyRecords = append(historyRecords, HistoryRecord{
			page:  lineParts[0],
			count: count,
		})
	}

	return &historyRecords, nil
}
