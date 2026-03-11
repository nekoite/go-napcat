package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultBotConfigSetsDefaultMaxReconnect(t *testing.T) {
	cfg := DefaultBotConfig(123456, "token")

	require.Equal(t, 3, cfg.Ws.MaxReconnect)
}

func TestBotConfigFromYamlParsesMaxReconnect(t *testing.T) {
	yamlConfig := []byte("ws:\n" +
		"  host: example.com\n" +
		"  port: 3001\n" +
		"  endpoint: /\n" +
		"  token: token\n" +
		"  maxReconnect: 7\n")

	cfg, err := BotConfigFromYaml(yamlConfig)

	require.NoError(t, err)
	require.Equal(t, 7, cfg.Ws.MaxReconnect)
	require.Equal(t, 10000, cfg.Ws.Timeout)
	require.Equal(t, 60000, cfg.Ws.PongTimeout)
}
