package stdapp

import (
	"context"
	"fmt"
	"os"

	"go.ddollar.dev/errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
)

func dockerClient() (*docker.Client, error) {
	dc, err := docker.NewClientWithOpts(
		docker.FromEnv,
		docker.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return dc, nil
}

func dockerProjectContainers(dc *docker.Client) ([]types.Container, error) {
	ctx := context.Background()

	host, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err)
	}

	c, err := dc.ContainerInspect(ctx, host)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	project, ok := c.Config.Labels["com.docker.compose.project"]
	if !ok {
		return nil, errors.Errorf("could not find docker compose project name")
	}

	opts := types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", fmt.Sprintf("com.docker.compose.project=%s", project)),
		),
	}

	cs, err := dc.ContainerList(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return cs, nil
}
