package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/dustin/go-humanize"
	"github.com/jon4hz/bermuda/cmd"
)

func main() {
	// test()
	cmd.Execute()
}

func test() {
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	/*
		containers, err := c.ContainerList(context.Background(), types.ContainerListOptions{
			All: true,
		})
		if err != nil {
			panic(err)
		}
		cm := dockerutils.ToContainerMap(containers)

		containerPruneReport, err := c.ContainersPrune(
			context.Background(),
			filters.NewArgs(
				filters.Arg("until", "1s"),
				filters.Arg("label!", "bermuda.exclude"),
			),
		)
		if err != nil {
			panic(err)
		}
		for _, r := range containerPruneReport.ContainersDeleted {
			fmt.Println("container")
			fmt.Println(cm[r].ID, cm[r].Names)
		}
		fmt.Println(containerPruneReport.SpaceReclaimed) */

	report, err := c.ImagesPrune(
		context.Background(),
		filters.NewArgs(
			filters.Arg("until", "0s"),
			filters.Arg("label!", "bermuda.exclude"),
			filters.Arg("dangling", "false"),
		),
	)
	if err != nil {
		panic(err)
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
}
