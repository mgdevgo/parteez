package config

import (
	"log"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/hjson"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	AppEnv      string     `koanf:"env"`
	HTTPServer  HTTPServer `koanf:"http"`
	DatabaseURL string     `koanf:"database_url"`
}

type HTTPServer struct {
	Address     string
	Timeout     time.Duration
	IdleTimeout time.Duration `koanf:"idle_timeout"`
}

func New() *Config {
	k := koanf.New(".")

	if err := k.Load(file.Provider(".env"), dotenv.ParserEnv("", ".", func(s string) string {
		//return strings.Replace(strings.ToLower(strings.TrimPrefix(s, "")), "_", ".", -1)
		return strings.ToLower(s)
	})); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	path := k.String("config_path")
	if path == "" {
		log.Fatal("cannot read config: env variable CONFIG_PATH is not set")
	}

	if err := k.Load(file.Provider(path), hjson.Parser()); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	var conf Config
	if err := k.Unmarshal("", &conf); err != nil {
		log.Fatalf("cannot unmarshal config: %s", err)
	}
	if err := k.UnmarshalWithConf("database.url", &conf.DatabaseURL, koanf.UnmarshalConf{FlatPaths: true}); err != nil {
		log.Fatalf("cannot unmarshal config: %s", err)
	}

	return &conf
}
