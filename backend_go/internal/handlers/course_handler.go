package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/services"
)

type CourseHandler struct {
	service *services.CourseService
}

func NewCourseHandler(service *services.CourseService) *CourseHandler {
	return &CourseHandler{service: service}
}

func (h *CourseHandler) GetCourses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := db.CourseFilter{
		Provider: c.Query("provider"),
		Mode:     c.Query("mode"),
		Level:    c.Query("level"),
	}

	courses, total, err := h.service.GetCourses(filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to fetch courses"))
		return
	}

	pagination := db.NewPagination(page, limit, total)
	c.JSON(http.StatusOK, db.SuccessResponse(courses, &pagination))
}

func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	id := c.Param("id")

	course, err := h.service.GetCourseByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, db.ErrorResponse(404, "Course not found"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(course, nil))
}
