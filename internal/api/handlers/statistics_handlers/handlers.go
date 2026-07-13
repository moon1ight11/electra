package statisticshandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// получение статистики по конкретному исполнителю
func (h *StatisticsHandler) ByWorker(c *gin.Context) {
	workerID, err := uuid.Parse(c.Param("workerId"))
	if err != nil {
		h.logger.Error("error in parce id in stat by worker", "error", err, "worker_id", workerID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid worker id"})
		return
	}

	from := c.Query("from")
	to := c.Query("to")

	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("error in parce id in stat by worker", "error", err, "owner_id", ownerID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid owner Id"})
		return
	}

	stats, err := h.statisticsService.ByWorker(c.Request.Context(), ownerID, workerID, from, to)
	if err != nil {
		h.logger.Error("error in service in stat by worker", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// получение статистики по всем исполнителям
func (h *StatisticsHandler) AllWorkers(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")

	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("error in parce id in stats by all workers", "error", err, "owner_id", ownerID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid owner Id"})
		return
	}

	stats, err := h.statisticsService.AllWorkers(c.Request.Context(), ownerID, from, to)
	if err != nil {
		h.logger.Error("error in service in stat by all workers", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// получение общей статистики
func (h *StatisticsHandler) Summary(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")

	stats, err := h.statisticsService.Summary(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
