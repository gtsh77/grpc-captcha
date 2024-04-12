package config

import (
	"time"
)

const ServicePrefix = "CAPTCHA"

type TLS struct {
	Enabled    bool
	DomainName string
	CrtData    string
	KeyData    string
	CrtCAData  string
}

type Config struct {
	Runtime struct {
		ID         string `conf:"-"` // авто-генерация
		Name       string `conf:"-"` // LD флаг
		Version    string `conf:"-"` // LD флаг
		CompiledAt string `conf:"-"` // LD флаг
		IsDevMode  bool   `conf:"default:false"`
	}
	Log struct {
		Level  int  `conf:"default:0"`
		AsJSON bool `conf:"default:false"`
	}
	HTTP struct {
		DomainNames []string      `conf:"default:*"`
		Host        string        `conf:"default:0.0.0.0"`
		Port        string        `conf:"default:1111"`
		Timeout     time.Duration `conf:"default:30s"`

		EnablePprof bool   `conf:"default:false"`
		EnableProm  bool   `conf:"default:true"`
		MetricPath  string `conf:"default:/metrics"`
		HealthPath  string `conf:"default:/health/check"`
		ReadyPath   string `conf:"default:/health/operable"`

		TLS *TLS
	}
	GRPC struct {
		Host       string        `conf:"default:0.0.0.0"`
		Port       string        `conf:"default:1111"`
		Timeout    time.Duration `conf:"default:30s"`
		XApiKey    string        `conf:"required"`
		EnableProm bool          `conf:"default:true"`

		TLS *TLS
	}
	Redis struct {
		Host        string `conf:"required"`
		Port        string `conf:"required"`
		DB          int    `conf:"required"`
		User        string
		Pass        string
		MaxExecTime time.Duration `conf:"default:3s"`
		MaxRetries  int           `conf:"default:2"`
		MaxIdleConn int           `conf:"default:5"`
		MaxOpenConn int           `conf:"default:25"`
		MaxConnTTL  time.Duration `conf:"default:1h"`

		TLS *TLS
	}
	Render struct {
		TTL    time.Duration `conf:"default:120s"`
		DigCnt int           `conf:"default:5"`
		Width  int           `conf:"default:180"`
		Height int           `conf:"default:80"`
	}
}

func New(name, version, compiledAt string) *Config {
	var config Config

	config.Runtime.Name = name
	config.Runtime.Version = version
	config.Runtime.CompiledAt = compiledAt

	return &config
}
