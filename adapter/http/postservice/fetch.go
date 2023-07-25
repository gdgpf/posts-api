package postservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (service service) Fetch(c *gin.Context) {
	posts, err := service.usecase.Fetch()

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, posts)
}
