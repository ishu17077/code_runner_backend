package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	c_runner "github.com/ishu17077/code_runner_backend/worker-node/helpers/code_runners/c"
	"github.com/ishu17077/code_runner_backend/worker-node/models"
)

// ! Simple odd even Test
var init_tests [2]models.TestCase = [2]models.TestCase{
	{
		Problem_id:     "69",
		Is_public:      true,
		Stdin:          "12",
		ExpectedOutput: "Yes",
		Test_id:        "1",
	},
	{
		Problem_id:     "69",
		Is_public:      true,
		Stdin:          "11",
		ExpectedOutput: "No",
		Test_id:        "1",
	},
}

func InitialTest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var submission models.Submission

		if err := c.ShouldBind(&submission); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}

		for _, test := range init_tests {
			res, err := c_runner.CheckSubmission(submission, test)
			if err != nil || res != "SUCCESS" {
				//TODO Fix this whatever
				var err = fmt.Sprintf("Error: %s\nTest Result: %s", err, res)
				c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
				return
			}
		}
		c.JSON(http.StatusAccepted, gin.H{"msg": "All Okay!!"})
	}
}
