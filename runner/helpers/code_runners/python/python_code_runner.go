package python

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
	var filePath = fmt.Sprintf("%s/main.py", dirPath)
	var outputPath = fmt.Sprintf("%s/main.pyc", dirPath)

	if err := coderunners.SaveFile(filePath, dirPath, submission.Code); err != nil {
		return "", dirPath, err
	}
	if err := compileCode(filePath, outputPath); err != nil {
		return "", dirPath, err
	}
	return outputPath, dirPath, nil
}

func PipeSubmission(filePath string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	runCmd := exec.CommandContext(ctx, "python", filePath)
	return coderunners.PipeCommand(runCmd)
}

func ExecuteSubmission(test models.TestCase, outputPath string) (string, error) {
	return executeCode(outputPath, test.Stdin)
}

func compileCode(filePath string, outputPath string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	py_compile_script := "import py_compile, sys; print(py_compile.compile(sys.argv[1], cfile=sys.argv[2]))"
	cmd := exec.CommandContext(ctx, "python", "-c", py_compile_script, filePath, outputPath)
	coderunners.SetPermissions(cmd)

	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}
	return nil
}

func executeCode(outputPath string, stdin string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	runCmd := exec.CommandContext(ctx, "python", outputPath)
	return coderunners.RunCommandWithInput(runCmd, ctx, stdin)
}
