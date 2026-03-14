package pipes

import (
	"gost/src/common/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateBody[T any](c *gin.Context) (*T, error) {
	var dto T
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.FormattedErrorGenerator(c, http.StatusBadRequest, "Bad Request", err.Error())
		return nil, err
	}
	return &dto, nil
}
