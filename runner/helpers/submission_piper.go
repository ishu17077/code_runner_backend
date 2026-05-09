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
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/perl"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/python"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/rust"
)

func ExecuteSubmissionAndPipe(submission models.Submission) error {
	lang := language.LanguageParser(submission.Language)
	switch lang {
	case language.C:
		err := executeCCode(submission)
		if err != nil {
			return err
		}
		return nil
	case language.Cpp:
		err := executeCppCode(submission)
		if err != nil {
			return err
		}
		return nil

	case language.Python:
		err := executePythonCode(submission)
		if err != nil {
			return err
		}
		return nil
	case language.Java:
		err := executeJavaCode(submission)
		if err != nil {
			return err
		}
		return nil
	case language.Cs:
		err := executeCSharpCode(submission)
		if err != nil {
			return err
		}
		return nil
	case language.Rust:
		err := executeRustCode(submission)
		if err != nil {
			return err
		}
		return nil
	case language.Go:
		err := executeGoCode(submission)
		if err != nil {
			return err
		}
		return nil
	case language.Perl:
		err := executePerlCode(submission)
		if err != nil {
			return err
		}
		return nil

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
