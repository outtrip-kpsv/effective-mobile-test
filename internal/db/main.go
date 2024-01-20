package db

import (
  "em_test/internal/config"
  "em_test/internal/db/repo"
  "errors"
  "fmt"
  "go.uber.org/zap"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
  "moul.io/zapgorm2"
  "strings"
  "time"
)

func NewDB(logger *zap.Logger, settings config.OptionsSrv) (*gorm.DB, error) {
  logger = logger.Named("PostgreSQL")
  dataSourceName := makeDsnStr(settings)
  db, err := getClient(dataSourceName, logger)
  if err != nil {
    return nil, err
  }
  sqlDB, err := db.DB()
  if err != nil {
    return nil, err
  }
  sqlDB.SetConnMaxLifetime(3 * time.Minute)
  db = db.Debug()

  err = db.AutoMigrate(&repo.Person{})
  if err != nil {
    logger.Error(err.Error())
  } else {
    logger.Info("migrate ok")
  }

  return db, nil
}

func makeDsnStr(settings config.OptionsSrv) string {
  parameters := map[string]string{
    "host":     settings.DbHost,
    "port":     settings.DbPort,
    "user":     settings.PgUser,
    "password": settings.PgPass,
    "dbname":   settings.DbName,
    "sslmode":  "disable",
  }
  var pairs []string
  for key, value := range parameters {
    pairs = append(pairs, fmt.Sprintf("%s=%s", key, value))
  }
  return strings.Join(pairs, " ")
}

func getClient(connectionString string, logger *zap.Logger) (*gorm.DB, error) {
  ticker := time.NewTicker(1 * time.Nanosecond)
  timeout := time.After(15 * time.Minute)
  seconds := 1
  gormLogger := zapgorm2.New(logger)

  for {
    select {
    case <-ticker.C:
      ticker.Stop()

      client, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{Logger: gormLogger})

      if err != nil {
        logger.With(zap.Error(err)).Warn("не удалось установить соединение с PostgreSQL")
        ticker = time.NewTicker(time.Duration(seconds) * time.Second)
        seconds *= 2
        if seconds > 60 {
          seconds = 60
        }
        continue
      }
      return client, nil
    case <-timeout:
      return nil, errors.New("PostgreSQL: connection timeout")
    }
  }
}
