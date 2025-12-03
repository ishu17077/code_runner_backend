package c_sharp

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	coderunners "github.com/ishu17077/code_runner_backend/worker-node/helpers/code_runners"
	"github.com/ishu17077/code_runner_backend/worker-node/models"
	currentstatus "github.com/ishu17077/code_runner_backend/worker-node/models/enums/current_status"
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

func CheckSubmission(submission models.Submission, test models.TestCase, filePath string) (currentstatus.CurrentStatus, error) {
	res, err := executeCode(filePath, test.Stdin)
	if err != nil {
		return currentstatus.FAILED, err
	}
	return coderunners.CheckOutput(res, test.ExpectedOutput)
}

func executeCode(filepath, stdin string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	runCmd := exec.CommandContext(ctx, "dotnet-script", filepath)
	return coderunners.RunCommandWithInput(runCmd, stdin)
}
