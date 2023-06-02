package stdapp

import (
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
)

func dockerClient() (*docker.Client, error) {
	dc, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, errors.Wrapf(err, "could not initialize docker client")
	}

	return dc, nil
}

func hostContainer(dc *docker.Client) (*docker.Container, error) {
	// host, err := os.Hostname()
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "could not get hostname")
	// }

	host := "network-cron-1"

	c, err := dc.InspectContainer(host)
	if err != nil {
		return nil, errors.Wrap(err, "could not inspect container")
	}

	return c, nil
}
