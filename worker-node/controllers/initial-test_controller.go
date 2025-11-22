package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ishu17077/code_runner_backend/worker-node/helpers"
	"github.com/ishu17077/code_runner_backend/worker-node/models"
)

// ! Simple odd even Test
var init_tests [2]models.TestCase = [2]models.TestCase{
	{
		Problem_id:     "69",
		Is_public:      true,
		Stdin:          "12\n",
		ExpectedOutput: "Yes",
		Test_id:        "1",
	},
	{
		Problem_id:     "28",
		Is_public:      true,
		Stdin:          "11\n",
		ExpectedOutput: "No",
		Test_id:        "1",
	},
}

//TODO: Implement sync.Mutex to handle process flow

func InitialTest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var submission models.Submission

		if err := c.ShouldBind(&submission); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}
		allOkay, execResults := helpers.AnalyzeSubmission(submission, init_tests[:])
		c.JSON(http.StatusAccepted, gin.H{"All tests passed": allOkay, "Execution Result": execResults})
	}
}
