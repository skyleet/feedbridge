package api

import (
	"fmt"

	"github.com/dewey/feedbridge/config"
	"github.com/dewey/feedbridge/plugin"
	"github.com/dewey/feedbridge/runner"
	"github.com/dewey/feedbridge/store"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Service provides access to the functions by the public API
type Service interface {
	// ServeFeed serves a feed based on the plugin and format
	ServeFeed(format string, plugin string) (string, error)

	// ListFeeds lists all available feed plugins
	ListFeeds() []string

	// RefreshFeed triggers a single plugin single scrape to force-refresh a feed based on the plugin name
	RefreshFeed(technicalName string) error
}

type service struct {
	l                 log.Logger
	cfg               config.Config
	storageRepository store.StorageRepository
	pluginRepository  plugin.Repository
	runner            *runner.Runner
}

// NewService initializes a new API service
func NewService(l log.Logger, cfg config.Config, sr store.StorageRepository, pr plugin.Repository, r *runner.Runner) *service {
	return &service{
		l:                 l,
		cfg:               cfg,
		storageRepository: sr,
		pluginRepository:  pr,
		runner:            r,
	}
}

// ServeFeed returns the generated feed from the store backend
func (s *service) ServeFeed(format string, plugin string) (string, error) {
	feed, err := s.storageRepository.Get(fmt.Sprintf("%s_%s", format, plugin))
	if err != nil {
		level.Error(s.l).Log("msg", err)
		return "", err
	}
	return feed, nil
}

func (s *service) ListFeeds() []plugin.PluginMetadata {
	var pp []plugin.PluginMetadata
	for _, p := range s.pluginRepository.All() {
		pp = append(pp, plugin.PluginMetadata{
			Name:          p.Info().Name,
			Description:   p.Info().Description,
			TechnicalName: p.Info().TechnicalName,
			Image:         p.Info().Image,
			Author:        p.Info().Author,
			AuthorURL:     p.Info().AuthorURL,
			SourceURL:     p.Info().SourceURL,
		})
	}
	return pp
}

func (s *service) RefreshFeed(technicalName string) error {
	cp, err := s.pluginRepository.Find(technicalName)
	if err != nil {
		return err
	}
	if err := s.runner.StartSingle(cp); err != nil {
		return err
	}
	return nil
}
