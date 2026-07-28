package perl

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/ishu17077/code_runner_backend/models"
	coderunners "github.com/ishu17077/code_runner_backend/runner/helpers/code_runners"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func PreCompilationTask(submission models.Submission) (string, string, error) {
	newId := bson.NewObjectID().Hex()
	var dirPath = fmt.Sprintf("/temp/%s", newId)
	var filePath = fmt.Sprintf("%s/main.pl", dirPath)

	if err := coderunners.SaveFile(filePath, dirPath, submission.Code); err != nil {
		return "", dirPath, err
	}
	return filePath, dirPath, nil
}

func PipeSubmission(filePath string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	runCmd := exec.CommandContext(ctx, "perl", filePath)
	return coderunners.PipeCommand(runCmd)
}

func ExecuteSubmission(testcase models.TestCase, filePath string) (string, error) {
	return executeCode(filePath, testcase.Stdin)

}

func executeCode(filePath string, stdin string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	runCmd := exec.CommandContext(ctx, "perl", filePath)
	return coderunners.RunCommandWithInput(runCmd, ctx, stdin)
}
