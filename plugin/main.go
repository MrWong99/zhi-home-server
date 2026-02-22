// zhi-config-homeserver is a zhi configuration plugin that provides default
// configuration values for a home server Docker Compose deployment.
//
// It serves values for the following components: core settings, PiHole,
// Plex, Nextcloud, MariaDB, Redis, and Nginx Proxy Manager.
package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"

	"github.com/MrWong99/zhi/pkg/zhiplugin"
	"github.com/MrWong99/zhi/pkg/zhiplugin/config"
)

func main() {
	level := hclog.LevelFromString(os.Getenv("ZHI_LOG_LEVEL"))
	if level == hclog.NoLevel {
		level = hclog.Info
	}
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "zhi-config-homeserver",
		Level:  level,
		Output: os.Stderr,
	})
	logger.Info("starting homeserver config plugin")

	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: zhiplugin.Handshake,
		Plugins: map[string]goplugin.Plugin{
			"config": &config.GRPCPlugin{Impl: newHomeserverPlugin()},
		},
		GRPCServer: goplugin.DefaultGRPCServer,
		Logger:     logger,
	})
}
