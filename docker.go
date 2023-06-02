package stdapp

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func dockerClient() (*docker.Client, error) {
	dc, err := docker.NewClientWithOpts(
		docker.FromEnv,
		docker.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not initialize docker client")
	}

	return dc, nil
}

func dockerProjectContainers(dc *docker.Client) ([]types.Container, error) {
	ctx := context.Background()

	host, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrapf(err, "could not get hostname")
	}

	c, err := dc.ContainerInspect(ctx, host)
	if err != nil {
		return nil, errors.Wrap(err, "could not inspect container")
	}

	project, ok := c.Config.Labels["com.docker.compose.project"]
	if !ok {
		return nil, fmt.Errorf("could not find docker compose project name")
	}

	opts := types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", fmt.Sprintf("com.docker.compose.project=%s", project)),
		),
	}

	cs, err := dc.ContainerList(ctx, opts)
	if err != nil {
		return nil, err
	}

	return cs, nil
}
