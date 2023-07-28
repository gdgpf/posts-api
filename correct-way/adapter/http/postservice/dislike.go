package postservice

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (service service) Dislike(c *gin.Context) {
	postIDString := c.Param("id")

	postID, err := strconv.Atoi(postIDString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Post inválido",
		})
		return
	}

	post, err := service.usecase.Dislike(postID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, post)
}
