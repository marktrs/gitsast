package git

import (
	"io"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/rs/zerolog/log"
)

type IClient interface {
	GetPathsFromRemoteURL(tmpDir string, remoteURL string) ([]string, error)
}

type client struct {
}

func NewClient() IClient {
	return &client{}
}

// getPathsFromRemoteURL - get all file paths from remote url except ignored file types
func (c *client) GetPathsFromRemoteURL(tmpDir string, remoteURL string) ([]string, error) {
	// local clone
	r, err := git.PlainClone(tmpDir, false, &git.CloneOptions{URL: remoteURL})
	if err != nil {
		return nil, err
	}

	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	return excludeIgnorePathsFromTree(object.NewTreeWalker(tree, true, nil), tmpDir), nil
}

var ignoredFileTypes = []string{
	".aac", ".aiff", ".ape", ".au", ".flac", ".gsm", ".it", ".m3u", ".m4a", ".mid", ".mod", ".mp3", ".mpa", ".pls", ".ra", ".s3m", ".sid", ".wav", ".wma", ".xm", ".7z",
	".a", ".ar", ".bz2", ".cab", ".cpio", ".deb", ".dmg", ".egg", ".gz", ".iso", ".lha", ".mar", ".pea", ".rar", ".rpm", ".s7z", ".shar", ".tar", ".tbz2", ".tgz", ".tlz",
	".whl", ".xpi", ".deb", ".rpm", ".xz", ".pak", ".crx", ".exe", ".msi", ".bin", ".eot", ".otf", ".ttf", ".woff", ".woff2", ".3dm", ".3ds", ".max", ".bmp", ".dds", ".gif",
	".jpg", ".jpeg", ".png", ".psd", ".xcf", ".tga", ".thm", ".tif", ".tiff", ".yuv", ".ai", ".eps", ".ps", ".svg", ".dwg", ".dxf", ".gpx", ".kml", ".kmz", ".ods", ".xls",
	".xlsx", ".csv", ".ics", ".vcf", ".ppt", ".odp", ".3g2", ".3gp", ".aaf", ".asf", ".avchd", ".avi", ".drc", ".flv", ".m2v", ".m4p", ".m4v", ".mkv", ".mng", ".mov", ".mp2",
	".mp4", ".mpe", ".mpeg", ".mpg", ".mpv", ".mxf", ".nsv", ".ogg", ".ogv", ".ogm", ".qt", ".rm", ".rmvb", ".roq", ".srt", ".svi", ".vob", ".webm", ".wmv", ".yuv",
}

func isFileTypeIgnored(filename string) bool {
	var isIgnored bool
	for _, ignoredFileType := range ignoredFileTypes {
		if strings.HasSuffix(filename, ignoredFileType) {
			isIgnored = true
		}
	}

	return isIgnored
}

func excludeIgnorePathsFromTree(treeWalker *object.TreeWalker, tmpDir string) []string {
	filepaths := make([]string, 0)
	for {
		name, _, err := treeWalker.Next()
		if err == io.EOF {
			break
		}

		isIgnored := isFileTypeIgnored(name)
		if isIgnored {
			continue
		}

		info, err := os.Stat(path.Join(tmpDir, name))
		if err != nil {
			log.Err(err)
			break
		}

		if info.IsDir() {
			continue
		}

		filepaths = append(filepaths, path.Join(tmpDir, name))
	}
	defer treeWalker.Close()

	return filepaths
}
