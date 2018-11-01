package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	op "github.com/lrx0014/service-catalog-test/pkg/operation"
)

func HandleSuccess(c *gin.Context, result interface{}, httpmethod string) {
	c.JSON(200, result)
}

func HandleError(c *gin.Context, request string, result map[string]interface{}, err error, httpmethod string) {
	response := make(map[string]interface{})
	if err.Error() == "exist" {
		c.JSON(409, response)
	} else if err.Error() == "noexist" {
		c.JSON(410, response)
	} else {
		response["error"] = request
		response["description"] = fmt.Sprintf("Error: %v", err)
		c.JSON(500, response)
	}
}

func AddBroker(c *gin.Context) {
	sta := op.AddBroker(op.NewClient())
	HandleSuccess(c, sta, "AddBroker")
}

func GetBroker(c *gin.Context) {
	sta := op.GetBroker(op.NewClient(), c.Param("broker"))
	HandleSuccess(c, sta, "GetBroker")
}
