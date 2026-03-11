package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

type WsConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Endpoint     string `yaml:"endpoint"`
	Token        string `yaml:"token"`
	Timeout      int    `yaml:"timeout"`     // in milliseconds
	PingPeriod   int    `yaml:"pingPeriod"`  // in milliseconds
	PongTimeout  int    `yaml:"pongTimeout"` // in milliseconds
	MaxReconnect int    `yaml:"maxReconnect"`
}

type BotConfig struct {
	Ws           WsConfig `yaml:"ws"`
	Id           int64    `yaml:"id"`
	Debug        bool     `yaml:"debug"`
	UseGoroutine bool     `yaml:"useGoroutine"`
	ApiTimeout   int      `yaml:"apiTimeout"`
}

type LogConfig struct {
	Level string   `yaml:"level"`
	Paths []string `yaml:"paths"`
	Debug bool     `yaml:"debug"`
}

var defaultBotCfg = BotConfig{
	Ws: WsConfig{
		Host:         "localhost",
		Port:         3001,
		Endpoint:     "/",
		Timeout:      10000,
		PingPeriod:   54000,
		PongTimeout:  60000,
		MaxReconnect: 3,
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

func (c *BotConfig) WithWsMaxReconnect(maxReconnect int) *BotConfig {
	c.Ws.MaxReconnect = maxReconnect
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
