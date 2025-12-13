package java

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	models "github.com/ishu17077/code_runner_backend/models"
	currentstatus "github.com/ishu17077/code_runner_backend/models/enums/current_status"
	coderunners "github.com/ishu17077/code_runner_backend/runner/helpers/code_runners"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// const outputPath = "/temp/Solution.class"

func PreCompilationTask(submission models.Submission) (string, string, string, error) {
	newId := bson.NewObjectID().Hex()
	var dirPath = fmt.Sprintf("/temp/%s", newId)
	var filePath = fmt.Sprintf("%s/Solution.java", dirPath)
	const javaExcutorPath = "./JavaExecutor.jar"
	var classPath = fmt.Sprintf("%s:%s", javaExcutorPath, dirPath)
	const className = "code_runner_backend.JavaExecutor"

	if err := coderunners.SaveFile(filePath, dirPath, submission.Code); err != nil {
		return "", classPath, dirPath, err
	}

	if err := compileCode(filePath, dirPath); err != nil {
		return "", classPath, dirPath, err
	}
	return className, classPath, dirPath, nil
}

func CheckSubmission(payload models.Payload, className, classPath string) (bool, models.Result, error) {
	var jsonPayload, err = json.Marshal(payload)

	if err != nil {
		return false, models.Result{
			Status:  currentstatus.INTERNAL_ERROR.ToString(),
			Results: []models.ExecResult{},
			Error:   err.Error(),
		}, err
	}

	var base64EncPayload string = base64.StdEncoding.EncodeToString(jsonPayload)

	res, err := executeCode(classPath, className, base64EncPayload)
	if err != nil {
		return false, models.Result{
			Status:  currentstatus.INTERNAL_ERROR.ToString(),
			Results: []models.ExecResult{},
			Error:   err.Error(),
		}, err
	}
	return coderunners.CheckJavaOutput(res, payload.Tests)
}

func compileCode(filepath, outputDir string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "javac", "-d", outputDir, filepath)
	//TODO: Check why???
	// coderunners.SetPermissions(cmd)
	res, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("Compilation Failed: %s %s", err.Error(), string(res))
	}
	return nil
}

func executeCode(classPath, className, stdin string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	runCmd := exec.CommandContext(ctx, "java",
		"-XX:+UseSerialGC",        //? Lightweight GC
		"-XX:TieredStopAtLevel=1", //? Fast startup, less optimization
		"-Xshare:on",              //? Use shared class data if available
		"-Xss256k",                //? Lower stack memory per thread
		"-Xms64m",                 //? Initial Heap
		"-Xmx128m",                //? Max Heap (Must be < Pod Limit mem)
		"-XX:-UsePerfData",
		"-cp", classPath, className)

	return coderunners.RunCommandWithInput(runCmd, stdin)
}
