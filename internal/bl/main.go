package bl

import (
	"em_test/internal/bl/models"
	"em_test/internal/db"
	"em_test/internal/db/repo"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"reflect"
	"sync"
)

type BL struct {
	Db     *db.DBRepo
	Logger *zap.Logger
}

func NewBL(repo *db.DBRepo, logger *zap.Logger) *BL {
	logger = logger.Named("Bl")

	return &BL{
		Db:     repo,
		Logger: logger,
	}
}
func (b *BL) PatchPerson(newPerson *repo.Person, oldPerson repo.Person) {
	b.Logger.Info("patch person", zap.Reflect("old", oldPerson), zap.Reflect("new", newPerson))

	if len(newPerson.Surname) == 0 {
		newPerson.Surname = oldPerson.Surname
	}
	if len(newPerson.Patronymic) == 0 {
		newPerson.Patronymic = oldPerson.Patronymic
	}
	if len(newPerson.Name) == 0 {
		newPerson.Name = oldPerson.Name
	}
	b.UpdateFromApi(newPerson)
	newPerson.ID = oldPerson.ID
}

func (b *BL) UpdateFromApi(person *repo.Person) {
	b.Logger.Info("Get data from api")
	var wg sync.WaitGroup
	wg.Add(3)
	out := make(chan interface{}, 3)
	go b.GetData(&wg, "api.agify.io", person, &models.Agify{}, out)
	go b.GetData(&wg, "api.genderize.io", person, &models.Genderize{}, out)
	go b.GetData(&wg, "api.nationalize.io", person, &models.Nationalize{}, out)
	wg.Wait()

	for data := range out {
		if len(out) == 0 {
			close(out)
		}
		dataType := reflect.TypeOf(data)
		switch dataType {
		case reflect.TypeOf(&models.Agify{}):
			agifyData, ok := data.(*models.Agify)
			if ok {
				person.Age = agifyData.Age
				b.Logger.Info("Agify", zap.String("name", agifyData.Name), zap.Int("age", agifyData.Age))

			}
		case reflect.TypeOf(&models.Genderize{}):
			genderizeData, ok := data.(*models.Genderize)
			if ok {
				person.Gender = genderizeData.Gender
				b.Logger.Info("Genderize", zap.String("name", genderizeData.Name), zap.String("Gender", genderizeData.Gender))

			}
		case reflect.TypeOf(&models.Nationalize{}):
			nationalizeData, ok := data.(*models.Nationalize)
			if ok {
				if len(nationalizeData.Country) > 0 {
					person.CountryID = nationalizeData.Country[0].CountryId
					b.Logger.Info("Nationalize", zap.String("name", nationalizeData.Name), zap.String("CountryID", person.CountryID))

				}
			}
		default:
			fmt.Println("Неизвестный тип данных")
		}
	}
}

func (b *BL) GetData(wg *sync.WaitGroup, url string, person *repo.Person, data interface{}, out chan interface{}) {
	defer wg.Done()

	url = fmt.Sprintf("https://%s/?name=%s", url, person.Name)
	b.Logger.Info("Get data from api", zap.String("url", url))
	response, err := http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		b.Logger.Warn("error decode", zap.Error(err))
		return
	}
	b.Logger.Info("response", zap.Reflect("data", data))

	out <- data
}
