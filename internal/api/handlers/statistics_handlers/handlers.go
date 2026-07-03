package statisticshandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ByWorker — статистика по конкретному исполнителю. Только владелец.
func (h *StatisticsHandler) ByWorker(c *gin.Context) {
	workerID, err := uuid.Parse(c.Param("workerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid worker id"})
		return
	}

	from := c.Query("from")
	to := c.Query("to")

	ownerID, _ := uuid.Parse(c.GetString("UserId"))

	stats, err := h.statisticsService.ByWorker(c.Request.Context(), ownerID, workerID, from, to)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// AllWorkers — статистика по всем. Только владелец.
func (h *StatisticsHandler) AllWorkers(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")

	ownerID, _ := uuid.Parse(c.GetString("UserId"))

	stats, err := h.statisticsService.AllWorkers(c.Request.Context(), ownerID, from, to)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
