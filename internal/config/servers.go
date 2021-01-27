package config

import "github.com/steve-nzr/goff-server/internal/domain/objects"

const ClusterAddress = "127.0.0.1"
const WorldAddress = "127.0.0.1"

var Servers = []*objects.Server{
	{
		Name: "Server 1",
		IP:   ClusterAddress,
		Channels: []*objects.Channel{
			{
				Name:      "Channel 1",
				IP:        WorldAddress,
				MaxPlayer: 500,
			},
		},
	},
}
