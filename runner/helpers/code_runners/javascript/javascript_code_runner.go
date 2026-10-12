package javascript

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/ishu17077/code_runner_backend/models"
	coderunners "github.com/ishu17077/code_runner_backend/runner/helpers/code_runners"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func PreCompilationTask(submission models.Submission) (filePath, dirPath string, err error) {
	newId := bson.NewObjectID().Hex()
	dirPath = fmt.Sprintf("/temp/%s", newId)
	filePath = fmt.Sprintf("%s/main.js", dirPath)
	if err := coderunners.SaveFile(filePath, dirPath, submission.Code); err != nil {
		return "", dirPath, err
	}
	return filePath, dirPath, err
}

func PipeSubmission(filePath string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	runCmd := exec.CommandContext(ctx, "node", filePath)
	return coderunners.PipeCommand(runCmd)
}

func ExecuteSubmission(testCase models.TestCase, filePath string) (string, error) {
	return executeCode(filePath, testCase.Stdin)
}

func executeCode(filePath, stdin string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	runCmd := exec.CommandContext(ctx, "node", filePath)
	return coderunners.RunCommandWithInput(runCmd, ctx, stdin)
}
