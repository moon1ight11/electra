package statisticshandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *StatisticsHandler) ByWorker(c *gin.Context) {
	workerID, err := uuid.Parse(c.Param("workerId"))
	if err != nil {
		h.logger.Error("failed to parse worker id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid worker id"})
		return
	}

	from := c.Query("from")
	to := c.Query("to")

	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("failed to parse owner id", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid owner Id"})
		return
	}

	stats, err := h.statisticsService.ByWorker(c.Request.Context(), ownerID, workerID, from, to)
	if err != nil {
		h.logger.Error("failed to get stats by worker", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *StatisticsHandler) AllWorkers(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")

	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("failed to parse owner id", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid owner Id"})
		return
	}

	stats, err := h.statisticsService.AllWorkers(c.Request.Context(), ownerID, from, to)
	if err != nil {
		h.logger.Error("failed to get stats all workers", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *StatisticsHandler) Summary(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")

	stats, err := h.statisticsService.Summary(c.Request.Context(), from, to)
	if err != nil {
		h.logger.Error("failed to get summary stats", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}