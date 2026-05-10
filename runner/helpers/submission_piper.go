package helpers

import (
	"fmt"

	"github.com/ishu17077/code_runner_backend/models"
	"github.com/ishu17077/code_runner_backend/models/enums/language"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/c"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/c_sharp"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/cpp"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/golang"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/java"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/javascript"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/perl"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/python"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/rust"
)

func ExecuteSubmissionAndPipe(submission models.Submission) error {
	lang := language.LanguageParser(submission.Language)
	switch lang {
	case language.C:
		return executeCCode(submission)
	case language.Cpp:
		return executeCppCode(submission)
	case language.Python:
		return executePythonCode(submission)
	case language.Java:
		return executeJavaCode(submission)
	case language.Cs:
		return executeCSharpCode(submission)
	case language.Rust:
		return executeRustCode(submission)
	case language.Go:
		return executeGoCode(submission)
	case language.Perl:
		return executePerlCode(submission)
	case language.Javascript:
		return executeJavascriptCode(submission)
	case language.Undefined:
		return fmt.Errorf("INVALID_LANGUAGE_PROVIDED")
	}
	return nil
}

func executeCCode(submission models.Submission) error {
	outputPath, dirPath, err := c.PreCompilationTask(submission)
	defer cleanUp(dirPath)
	if err != nil {
		return fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	return c.PipeSubmission(outputPath)
}
func executeCppCode(submission models.Submission) error {
	outputPath, dirPath, err := cpp.PreCompilationTask(submission)
	defer cleanUp(dirPath)

	if err != nil {
		return fmt.Errorf("Error compiling the file: %s", err.Error())
	}

	return cpp.PipeSubmission(outputPath)
}

func executePythonCode(submission models.Submission) error {
	outputPath, dirPath, err := python.PreCompilationTask(submission)

	defer cleanUp(dirPath)

	if err != nil {
		return fmt.Errorf("Error compiling the file: %s", err.Error())
	}

	return python.PipeSubmission(outputPath)
}

func executeRustCode(submission models.Submission) error {
	outputPath, dirPath, err := rust.PreCompilationTask(submission)

	defer cleanUp(dirPath)

	if err != nil {
		return fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	return rust.PipeSubmission(outputPath)
}

func executeCSharpCode(submission models.Submission) error {
	outputPath, dirPath, err := c_sharp.PreCompilationTask(submission)

	defer cleanUp(dirPath)

	if err != nil {
		return fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	return c_sharp.PipeSubmission(outputPath)
}

func executeJavaCode(submission models.Submission) error {
	className, classPath, dirPath, err := java.PreCompilationTask(submission)

	defer cleanUp(dirPath)

	if err != nil {
		return fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	return java.PipeSubmission(classPath, className)
}

func executePerlCode(submission models.Submission) error {
	outputPath, dirPath, err := perl.PreCompilationTask(submission)

	defer cleanUp(dirPath)

	if err != nil {
		return fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	return perl.PipeSubmission(outputPath)
}

func executeGoCode(submission models.Submission) error {
	outputPath, dirPath, err := golang.PreCompilationTask(submission)

	defer cleanUp(dirPath)

	if err != nil {
		return fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	return golang.PipeSubmission(outputPath)
}
func executeJavascriptCode(submission models.Submission) error {
	filePath, dirPath, err := javascript.PreCompilationTask(submission)
	defer cleanUp(dirPath)

	if err != nil {
		return fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	return javascript.PipeSubmission(filePath)
}
