package c_sharp

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/ishu17077/code_runner_backend/models"
	coderunners "github.com/ishu17077/code_runner_backend/runner/helpers/code_runners"
	"go.mongodb.org/mongo-driver/v2/bson"
)

//TODO: Compile with dotnet 10 sdk the 9 doens't really work

func PreCompilationTask(submission models.Submission) (string, string, error) {
	newId := bson.NewObjectID().Hex()
	var dirPath = fmt.Sprintf("/temp/%s", newId)
	var filePath = fmt.Sprintf("%s/Program.cs", dirPath)

	//TODO: ADD PRECOMPILATION Step with dotnet
	//! dotnet new console -n HelloWorld Replace Helloword with newId

	if err := coderunners.SaveFile(filePath, dirPath, submission.Code); err != nil {
		return "", dirPath, err
	}
	return filePath, dirPath, nil
}

func CheckSubmission(test models.TestCase, filePath string) (string, error) {
	return executeCode(filePath, test.Stdin)
}

func executeCode(filepath, stdin string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	runCmd := exec.CommandContext(ctx, "dotnet-script", filepath)
	return coderunners.RunCommandWithInput(runCmd, stdin)
}
