package bermuda

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/client"
	"github.com/dustin/go-humanize"
	"github.com/jon4hz/bermuda/internal/config"
)

func Run(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}
	if err := pruneContainer(cfg.Container); err != nil {
		return err
	}
	if err := pruneImage(cfg.Image); err != nil {
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
	var s strings.Builder
	for _, r := range report.ContainersDeleted {
		s.WriteString("deleted container: ")
		s.WriteString(r)
		s.WriteByte('\n')
	}
	s.WriteString(humanize.Bytes(report.SpaceReclaimed))

	fmt.Println(s.String())

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
	var s strings.Builder
	for _, r := range report.ImagesDeleted {
		if r.Untagged != "" {
			s.WriteString("untagged image: ")
			s.WriteString(r.Untagged)
			s.WriteByte('\n')
		}
		if r.Deleted != "" {
			s.WriteString("deleted image: ")
			s.WriteString(r.Deleted)
			s.WriteByte('\n')
		}
	}
	s.WriteString(humanize.Bytes(report.SpaceReclaimed))

	fmt.Println(s.String())

	return nil
}
