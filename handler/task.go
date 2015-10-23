package handler

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/seesawlabs/ivan-kirichenko-exercise/model"
)

// GetGetTaskHandler creates HTTP handler for Get Task operation
func GetGetTaskHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c *echo.Context) error {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, NewApiError(err.Error()))
		}

		task := model.Task{}
		if err := db.First(&task, id).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		}

		return c.JSON(http.StatusOK, task)
	}
}

// GetGetTaskHandler creates HTTP handler for Create Task operation
func GetCreateTaskHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c *echo.Context) error {
		task := model.Task{}
		if err := c.Bind(&task); err != nil {
			return err
		}

		if err := db.Create(&task).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		}

		return c.JSON(http.StatusCreated, task)
	}
}

// GetGetTaskHandler creates HTTP handler for Update Task operation
func GetUpdateTaskHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c *echo.Context) error {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, NewApiError(err.Error()))
		}

		task := model.Task{}
		if err := c.Bind(&task); err != nil {
			return err
		}
		task.Id = id

		if err := db.Save(&task).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		}

		return c.NoContent(http.StatusNoContent)
	}
}

// GetGetTaskHandler creates HTTP handler for Delete Task operation
func GetDeleteTaskHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c *echo.Context) error {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, NewApiError(err.Error()))
		}

		task := model.Task{Id: id}
		if err := db.Delete(&task).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		}

		return c.JSON(http.StatusOK, task)
	}
}
