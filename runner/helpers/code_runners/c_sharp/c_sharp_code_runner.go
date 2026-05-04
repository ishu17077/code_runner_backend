package c_sharp

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/ishu17077/code_runner_backend/models"
	coderunners "github.com/ishu17077/code_runner_backend/runner/helpers/code_runners"
)

//TODO: Compile with dotnet 10 sdk the 9 doens't really work

func PreCompilationTask(submission models.Submission) (string, string, error) {
	var dirPath = "/opt/dotnet-project"
	var filePath = fmt.Sprintf("%s/Program.cs", dirPath)
	var outputPath = fmt.Sprintf("%s/bin/release/net10.0/dotnet", dirPath)

	//TODO: ADD PRECOMPILATION Step with dotnet
	//! dotnet new console -n HelloWorld Replace Helloword with newId
	//? Thiw would automatically change dirs to filePath
	if err := coderunners.SaveFile(filePath, dirPath, submission.Code); err != nil {
		return "", dirPath, err
	}

	if err := compileCode(); err != nil {
		return filePath, dirPath, err
	}
	return outputPath, dirPath, nil
}

func CheckSubmission(test models.TestCase, binaryFile string) (string, error) {
	return executeCode(binaryFile, test.Stdin)
}

func compileCode() error {
	var ctx, cancel = context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "dotnet", "build", "-c", "release")
	coderunners.SetPermissions(cmd)

	if res, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Compilation Failed: %s %s", err.Error(), string(res))
	}
	return nil
}

func executeCode(binaryFilePath, stdin string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	runCmd := exec.CommandContext(ctx, binaryFilePath)
	return coderunners.RunCommandWithInput(runCmd, stdin)
}
