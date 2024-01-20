package db

import (
  "em_test/internal/db/repo"
  "gorm.io/gorm"
)

type DBRepo struct {
  db     *gorm.DB
  People repo.PersonRepository
}

func NewDBRepo(db *gorm.DB) *DBRepo {
  return &DBRepo{
    db:     db,
    People: repo.NewPeopleRepository(db),
  }
}
