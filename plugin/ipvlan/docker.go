package ipvlan

import "github.com/samalba/dockerclient"

type dockerer struct {
	client *dockerclient.DockerClient
}
