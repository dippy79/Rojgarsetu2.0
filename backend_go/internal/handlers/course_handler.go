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

// @Summary Get courses
// @Description Retrieve paginated list of courses with optional filters
// @Tags courses
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param provider query string false "Filter by provider"
// @Param mode query string false "Filter by mode"
// @Param level query string false "Filter by level"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/courses [get]
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

	respCourses := make([]CourseResponse, 0, len(courses))
	for _, c := range courses {
		respCourses = append(respCourses, toCourseResponse(c))
	}

	c.JSON(http.StatusOK, db.SuccessResponse(respCourses, &pagination))
}

// @Summary Get course providers
// @Description Retrieve list of course providers
// @Tags courses
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/courses/providers [get]
func (h *CourseHandler) GetCourseProviders(c *gin.Context) {
	providers, err := h.service.GetCourseProviders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to fetch course providers"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(providers, nil))
}

func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	id := c.Param("id")

	course, err := h.service.GetCourseByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, db.ErrorResponse(404, "Course not found"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(toCourseByIDResponse(*course), nil))
}
