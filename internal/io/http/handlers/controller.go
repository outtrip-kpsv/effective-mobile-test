package handlers

import (
	"em_test/internal/bl"
	"em_test/internal/db/repo"
	response "em_test/internal/io/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strconv"
)

type Controller struct {
	bl     *bl.BL
	logger *zap.Logger
}

func NewController(bl *bl.BL, logger *zap.Logger) *Controller {
	logger = logger.Named("Handler")

	return &Controller{bl: bl, logger: logger}
}

// GetPeople @Summary Получить список людей
// @Description Получить список людей с возможностью фильтрации и пагинации
// @Tags response.Peopdle
// @Accept json
// @Produce json
// @Param name query string false "Имя"
// @Param surname query string false "Фамилия"
// @Param patronymic query string false "Отчество"
// @Param age query int false "Возраст"
// @Param gender query string false "Пол"
// @Param country_id query string false "Идентификатор страны"
// @Param page query int false "Номер страницы"
// @Param page_size query int false "Размер страницы"
// @Success 200 {object} models.People
// @Failure 500 {object} models.ErrResp
// @Router /people [get]
func (ct *Controller) GetPeople(c *fiber.Ctx) error {
	ct.logger.Info("GetPeople", zap.String("fiberCtx", c.String()))
	var filter repo.Person

	filter.Name = c.Query("name")
	filter.Surname = c.Query("surname")
	filter.Patronymic = c.Query("patronymic")
	filter.Age = c.QueryInt("age", 0)
	filter.Gender = c.Query("gender")
	filter.CountryID = c.Query("country_id")

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	people, totalCount, err := ct.bl.Db.People.GetAllPeople(filter, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrResp{Err: "Ошибка при запросе к базе данных"})
	}
	prevPage := 0
	nextPage := 0
	if page > 1 {
		prevPage = page - 1
	}
	if page*pageSize < int(totalCount) {
		nextPage = page + 1
	}

	resp := response.People{
		Data:       people,
		TotalCount: int(totalCount),
		Limit:      pageSize,
		Page:       page,
		PrevPage:   prevPage,
		NextPage:   nextPage,
	}
	return c.JSON(resp)
}

// CreatePerson @Summary Создать нового человека
// @Description Создать нового человека с возможностью обогащения данных из API
// @Tags People
// @Accept json
// @Produce json
// @Param person body repo.Person true "Данные нового человека"
// @Success 200 {object} repo.Person
// @Failure 400 {object} models.ErrResp
// @Failure 500 {object} models.ErrResp
// @Router /create [post]
func (ct *Controller) CreatePerson(c *fiber.Ctx) error {
	ct.logger.Info("CreatePerson", zap.String("fiberCtx", c.String()))

	var person, personNil repo.Person
	if err := c.BodyParser(&person); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrResp{Err: "Неверный формат данных"})
	}
	//todo to bl
	if person == personNil && len(person.Name) != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrResp{Err: "пустой запрос"})
	}

	ct.bl.UpdateFromApi(&person)
	person.ID = 0
	err := ct.bl.Db.People.AddPerson(&person)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при добавлении человека в базу данных")
	}
	return c.JSON(person)
}

// DeletePerson @Summary Удалить человека по ID
// @Description Удалить человека по указанному идентификатору
// @Tags People
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор человека"
// @Success 200 {object} models.OkResp
// @Failure 400 {object} models.ErrResp
// @Failure 404 {object} models.ErrResp
// @Failure 500 {object} models.ErrResp
// @Router /del/{id} [delete]
func (ct *Controller) DeletePerson(c *fiber.Ctx) error {
	ct.logger.Info("DeletePeople", zap.String("fiberCtx", c.String()))

	id := c.Params("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrResp{Err: "Ошибка при преобразовании строки в int: " + err.Error()})
	}

	person, err := ct.bl.Db.People.GetPersonByID(uint(i))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrResp{Err: "Ошибка при получении информации о человеке: " + err.Error()})
	}

	if person.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrResp{Err: "Человек не найден"})
	}

	err = ct.bl.Db.People.DeletePersonByID(uint(i))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrResp{Err: "Ошибка при удалении человека из базы данных: " + err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(response.OkResp{Msg: fmt.Sprintf("Человек с ID %s успешно удален", id)})
}

// UpdatePerson @Summary Обновить информацию о человеке по ID
// @Description Обновить информацию о человеке по указанному идентификатору
// @Tags People
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор человека"
// @Param person body repo.Person true "Обновленные данные человека"
// @Success 200 {object} models.OkResp
// @Failure 400 {object} models.ErrResp
// @Failure 404 {object} models.ErrResp
// @Failure 500 {object} models.ErrResp
// @Router /update/{id} [patch]
func (ct *Controller) UpdatePerson(c *fiber.Ctx) error {
	ct.logger.Info("UpdatePeople", zap.String("fiberCtx", c.String()))

	id := c.Params("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrResp{Err: "Ошибка при преобразовании строки в int: " + err.Error()})
	}

	existingPerson, err := ct.bl.Db.People.GetPersonByID(uint(i))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrResp{Err: "Ошибка при получении информации о человеке: " + err.Error()})
	}

	if existingPerson.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrResp{Err: "Человек не найден"})
	}

	var updatedPerson repo.Person
	if err := c.BodyParser(&updatedPerson); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrResp{Err: "Ошибка при парсинге тела запроса: " + err.Error()})
	}

	ct.bl.PatchPerson(&updatedPerson, *existingPerson)

	if err := ct.bl.Db.People.UpdatePerson(&updatedPerson); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrResp{Err: "Ошибка при обновлении информации о человеке: " + err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(response.OkResp{Msg: fmt.Sprintf("Человек с ID %s успешно обновлен", id)})
}
