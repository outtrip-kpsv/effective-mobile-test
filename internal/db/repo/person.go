package repo

import (
  "errors"
  "fmt"
  "gorm.io/gorm"
)

type PersonRepositoryImpl struct {
  db *gorm.DB
}

func NewPeopleRepository(db *gorm.DB) *PersonRepositoryImpl {
  return &PersonRepositoryImpl{db: db}
}

type PersonRepository interface {
  GetAllPeople(filter Person, page, pageSize int) ([]Person, int64, error)
  GetPersonByID(id uint) (*Person, error)
  DeletePersonByID(id uint) error
  UpdatePerson(person *Person) error
  AddPerson(person *Person) error
}

type Person struct {
  ID         uint   `gorm:"primaryKey"`
  Name       string `json:"name" validate:"required"`
  Surname    string `json:"surname"`
  Patronymic string `json:"patronymic"`
  Age        int    `json:"age"`
  Gender     string `json:"gender"`
  CountryID  string `json:"country_id"`
}

func (r *PersonRepositoryImpl) GetAllPeople(filter Person, page, pageSize int) ([]Person, int64, error) {
  query := r.db.Model(&Person{})

  if filter.Name != "" {
    query = query.Where("name = ?", filter.Name)
  }

  if filter.Surname != "" {
    query = query.Where("surname = ?", filter.Surname)
  }

  if filter.Patronymic != "" {
    query = query.Where("patronymic = ?", filter.Patronymic)
  }

  if filter.Age > 0 {
    query = query.Where("age = ?", filter.Age)
  }

  if filter.Gender != "" {
    query = query.Where("gender = ?", filter.Gender)
  }

  if filter.CountryID != "" {
    query = query.Where("country_id = ?", filter.CountryID)
  }

  var totalCount int64
  var people []Person

  // Получаем общее количество записей
  if err := query.Count(&totalCount).Error; err != nil {
    return nil, 0, err
  }

  // Применяем пагинацию
  offset := (page - 1) * pageSize
  query = query.Order("id ASC")
  query = query.Offset(offset).Limit(pageSize)

  // Выполняем запрос и получаем данные
  if err := query.Find(&people).Error; err != nil {
    return nil, 0, err
  }

  return people, totalCount, nil
}

func (r *PersonRepositoryImpl) GetPersonByID(id uint) (*Person, error) {
  var person Person
  if err := r.db.First(&person, id).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      // Запись не найдена
      return &person, nil
    }
    // Другая ошибка
    fmt.Println("Ошибка при запросе к базе данных:", err)
    return nil, err
  }
  return &person, nil
}

func (r *PersonRepositoryImpl) DeletePersonByID(id uint) error {
  return r.db.Delete(&Person{}, id).Error
}

func (r *PersonRepositoryImpl) UpdatePerson(person *Person) error {
  return r.db.Save(person).Error
}

func (r *PersonRepositoryImpl) AddPerson(person *Person) error {
  return r.db.Create(person).Error
}
