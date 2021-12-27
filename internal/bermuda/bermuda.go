package bermuda

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/client"
	"github.com/dustin/go-humanize"
	"github.com/jon4hz/bermuda/internal/config"
	"github.com/jon4hz/bermuda/internal/logger"
	"go.uber.org/zap"
)

func Run(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}
	if err := pruneContainer(cfg.Container); err != nil {
		logger.Log.Error("error while pruning containers", zap.Error(err))
		return err
	}
	if err := pruneImage(cfg.Image); err != nil {
		logger.Log.Error("error while pruning images", zap.Error(err))
		return err
	}
	return nil
}

func pruneContainer(cfg *config.PruneContainerConfig) error {
	if !cfg.Active {
		return nil
	}
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	report, err := c.ContainersPrune(
		context.Background(),
		cfg.ToFilterArgs(),
	)
	if err != nil {
		return err
	}
	if len(report.ContainersDeleted) > 0 {
		for _, r := range report.ContainersDeleted {
			var s strings.Builder
			s.WriteString("deleted container: ")
			s.WriteString(r)
			logger.Log.Info(s.String())
		}
		logger.Log.Info("claimed " + humanize.Bytes(report.SpaceReclaimed))
	}
	return nil
}

func pruneImage(cfg *config.PruneImageConfig) error {
	if !cfg.Active {
		return nil
	}
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	report, err := c.ImagesPrune(
		context.Background(),
		cfg.ToFilterArgs(),
	)
	if err != nil {
		return err
	}
	if len(report.ImagesDeleted) > 0 {
		for _, r := range report.ImagesDeleted {
			var s strings.Builder
			if r.Untagged != "" {
				s.WriteString("untagged image: ")
				s.WriteString(r.Untagged)
			}
			if r.Deleted != "" {
				s.WriteString("deleted image: ")
				s.WriteString(r.Deleted)
			}
			logger.Log.Info(s.String())
		}
		logger.Log.Info("claimed " + humanize.Bytes(report.SpaceReclaimed))
	}
	return nil
}
