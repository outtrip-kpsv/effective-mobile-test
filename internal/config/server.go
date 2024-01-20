package config

import (
  "github.com/jessevdk/go-flags"
  "go.uber.org/zap"
)

type OptionsSrv struct {
  Host   string `short:"h" long:"host" description:"хост" default:"localhost" env:"HOST"`
  Port   string `short:"p" long:"port" description:"порт" default:"8765" env:"PORT"`
  Log    string `long:"logger-create" description:"logger-create output" default:"debug" env:"LOG"`
  DbHost string `long:"dbhost" description:"the db server host" default:"localhost" env:"DB_HOST"`
  DbPort string `long:"dbport" description:"the db server port" default:"5432" env:"DB_PORT"`
  PgUser string `long:"pguser" description:"the db user" default:"postgres" env:"POSTGRES_USER"`
  PgPass string `long:"pgpass" description:"the db pass" default:"postgres" env:"POSTGRES_PASSWORD"`
  DbName string `long:"dbname" description:"the db name" default:"namedb" env:"POSTGRES_DB"`
}

type ConfSrv struct {
  Options OptionsSrv
  Logger  *zap.Logger
}

func InitConfServ() (*ConfSrv, error) {
  var conf ConfSrv
  var opts OptionsSrv
  parser := flags.NewParser(&opts, flags.Default)
  _, err := parser.Parse()
  if err != nil {
    return nil, err
  }
  logger := initLogger(opts.Log)
  conf.Options = opts
  conf.Logger = logger
  return &conf, nil
}
