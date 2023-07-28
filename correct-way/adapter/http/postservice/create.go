package postservice

import (
	"net/http"

	"github.com/gdgpf/posts-api/adapter/validate"
	"github.com/gdgpf/posts-api/core/dto"
	"github.com/gin-gonic/gin"
)

func (service service) Create(c *gin.Context) {
	postRequest, err := dto.FromJsonPostRequestBody(c.Request.Body, service.validator)

	if err != nil {
		if _, ok := err.(validate.RequestError); ok {
			c.JSON(http.StatusUnprocessableEntity, err.(validate.RequestError))
			return
		}

		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	account, err := service.usecase.Create(postRequest)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, account)
}
