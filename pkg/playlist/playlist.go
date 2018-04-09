package playlist

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"strings"

	"broadcastle.co/code/icii/pkg/config"
	"broadcastle.co/code/icii/pkg/util"
)

var playlist []string
var idx int
var np string

// First returns the first file in the playlist.
func First() string {
	if len(playlist) > 0 {
		return playlist[0]
	}
	return ""
}

// Next returns the next file to be played from the playlist.
func Next() string {
	//save_idx;

	// get_next_file := pl.Strings[idx];
	if idx > len(playlist)-1 {
		idx = 0
	}
	np = playlist[idx]
	Load()
	if idx > len(playlist)-1 {
		idx = 0
	}
	for (np == playlist[idx]) && (len(playlist) > 1) {
		if !config.Cfg.PlayRandom {
			idx = idx + 1
			if idx > len(playlist)-1 {
				idx = 0
			}
		} else {
			idx = rand.Intn(len(playlist))
		}
	}
	return playlist[idx]
}

// Load gets the playlist file.
func Load() error {
	if ok := util.FileExists(config.Cfg.Playlist); !ok {
		return errors.New("Playlist file doesn't exist")
	}

	content, err := ioutil.ReadFile(config.Cfg.Playlist)
	if err != nil {
		return err
	}
	playlist = strings.Split(string(content), "\n")

	i := 0
	for i < len(playlist) {
		playlist[i] = strings.Replace(playlist[i], "\r", "", -1)
		if ok := util.FileExists(playlist[i]); !ok && !strings.HasPrefix(playlist[i], "http") {
			playlist = append(playlist[:i], playlist[i+1:]...)
			continue
		}
		i++
	}
	if len(playlist) < 1 {
		return errors.New("Error: all files in the playlist do not exist")
	}

	return nil
}
