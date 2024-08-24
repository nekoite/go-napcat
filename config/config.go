package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

type WsConfig struct {
	Host     string
	Port     int
	Endpoint string
	Token    string
}

type BotConfig struct {
	Ws           WsConfig
	Id           int64
	Debug        bool
	UseGoroutine bool
	ApiTimeout   int
}

type LogConfig struct {
	Level string
	Paths []string
	Debug bool
}

var defaultBotCfg = BotConfig{
	Ws: WsConfig{
		Host:     "localhost",
		Port:     3001,
		Endpoint: "/",
	},
	ApiTimeout: 30000,
}

func BotConfigFromYamlFile(path string) (*BotConfig, error) {
	s, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return BotConfigFromYaml(s)
}

func BotConfigFromYaml(s []byte) (*BotConfig, error) {
	cfg := defaultBotCfg
	err := yaml.Unmarshal(s, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func DefaultBotConfig(id int64, token string) *BotConfig {
	cfg := defaultBotCfg
	cfg.Id = id
	cfg.Ws.Token = token
	return &cfg
}

func (c *BotConfig) WithWs(host string, port int, endpoint string) *BotConfig {
	c.Ws.Host = host
	c.Ws.Port = port
	c.Ws.Endpoint = endpoint
	return c
}

func (c *BotConfig) DebugMode(debug bool) *BotConfig {
	c.Debug = debug
	return c
}

func (c *BotConfig) GoroutineMode(useGoroutine bool) *BotConfig {
	c.UseGoroutine = useGoroutine
	return c
}

func (c *BotConfig) WithApiTimeout(timeout int) *BotConfig {
	c.ApiTimeout = timeout
	return c
}

func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		Level: "info",
	}
}

func (c *LogConfig) WithStderr() *LogConfig {
	c.Paths = append(c.Paths, "stderr")
	return c
}

func (c *LogConfig) WithStdout() *LogConfig {
	c.Paths = append(c.Paths, "stdout")
	return c
}

func (c *LogConfig) WithLevel(level string) *LogConfig {
	c.Level = level
	return c
}

func (c *LogConfig) WithPaths(path ...string) *LogConfig {
	c.Paths = append(c.Paths, path...)
	return c
}
