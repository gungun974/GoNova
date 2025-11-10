package watch

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gohugoio/hugo/watcher/filenotify"
	"github.com/gungun974/gonova/internal/logger"

	"github.com/bmatcuk/doublestar/v4"
)

type WatcherConfig struct {
	Extensions []string
	Filter     []string
	Ignore     []string
}

func WatchNewFiles(config WatcherConfig, updateFile func()) {
	watcher, err := newWatcher()
	if err != nil {
		logger.WatcherLogger.Fatal(err)
	}

	err = filepath.Walk(".", func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() && !isValidDirectory(config, path) {
			return filepath.SkipDir
		}

		if !info.IsDir() && !isValidFile(config, path) {
			return nil
		}

		err = watcher.Add(path)
		if err != nil {
			logger.WatcherLogger.Errorf(
				"Fail to add path to watch %s : %v",
				path,
				err,
			)
		}

		return nil
	})
	if err != nil {
		logger.WatcherLogger.Errorf(
			"Fail walking path for watching : %v",
			err,
		)
	}

	for {
		select {
		case ev := <-watcher.Events():

			if ev.Has(fsnotify.Remove) {
				_ = watcher.Remove(ev.Name)
				if isValidFile(config, ev.Name) {
					updateFile()
				}
			}
			if ev.Has(fsnotify.Create) || ev.Has(fsnotify.Write) {
				info, err := os.Stat(ev.Name)
				if err != nil {
					if os.IsNotExist(err) {
						_ = watcher.Remove(ev.Name)
						continue
					}
					logger.WatcherLogger.Fatal(err)
				}

				if !info.IsDir() && !isValidFile(config, ev.Name) {
					continue
				}

				if ev.Has(fsnotify.Create) {
					err = watcher.Add(ev.Name)
					if err != nil {
						logger.WatcherLogger.Errorf(
							"Fail to add path to watch %s : %v",
							ev.Name,
							err,
						)
					}
					if !info.IsDir() {
						updateFile()
					}
				}
				if ev.Has(fsnotify.Write) && !info.IsDir() {
					updateFile()
				}

			}
		case err := <-watcher.Errors():
			logger.WatcherLogger.Error(err)
		}
	}
}

func isValidFile(config WatcherConfig, path string) bool {
	if len(config.Extensions) != 0 &&
		!slices.Contains(
			config.Extensions,
			strings.TrimPrefix(filepath.Ext(filepath.Base(path)), "."),
		) {
		return false
	}

	if len(config.Filter) != 0 {
		for _, filter := range config.Filter {
			match, err := doublestar.PathMatch(filter, path)
			if err == nil && match {
				return true
			}
			match, err = doublestar.PathMatch(filter, filepath.Base(path))
			if err == nil && match {
				return true
			}
		}
	}

	if len(config.Ignore) != 0 {
		for _, ignore := range config.Ignore {
			match, err := doublestar.PathMatch(ignore, path)
			if err == nil && match {
				return false
			}
			match, err = doublestar.PathMatch(ignore, filepath.Base(path))
			if err == nil && match {
				return false
			}
		}
	}

	return true
}

func isValidDirectory(config WatcherConfig, path string) bool {
	if len(config.Filter) != 0 {
		for _, filter := range config.Filter {
			match, err := doublestar.PathMatch(filter, path)
			if err == nil && match {
				return true
			}
			match, err = doublestar.PathMatch(filter, filepath.Base(path))
			if err == nil && match {
				return true
			}
		}
	}

	if len(config.Ignore) != 0 {
		for _, ignore := range config.Ignore {
			match, err := doublestar.PathMatch(ignore, path)
			if err == nil && match {
				return false
			}
			match, err = doublestar.PathMatch(ignore, filepath.Base(path))
			if err == nil && match {
				return false
			}
		}
	}

	return true
}

func newWatcher() (filenotify.FileWatcher, error) {
	return filenotify.NewEventWatcher()

	// if interval < 500 {
	// 	interval = 500
	// }
	// pollInterval := time.Duration(interval) * time.Millisecond
	//
	// return filenotify.NewPollingWatcher(pollInterval), nil
}
