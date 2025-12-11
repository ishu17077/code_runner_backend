package golang

// import (
// 	"fmt"
// 	"os/exec"

// 	"github.com/ishu17077/code_runner_backend/models"
// 	coderunners "github.com/ishu17077/code_runner_backend/runner/helpers/code_runners"
// 	"go.mongodb.org/mongo-driver/v2/bson"
// )

// func PreCompilationTask(submission models.Submission) (string, string, error) {
// 	newId := bson.NewObjectID().Hex()
// 	var dirPath = fmt.Sprintf("/temp/%s", newId)
// 	var filePath = fmt.Sprintf("%s/main.go", dirPath)
// 	var outputPath = fmt.Sprintf("%s/main", dirPath)

// 	if err := coderunners.SaveFile(filePath, dirPath, submission.Code); err != nil {
// 		return "", dirPath, err
// 	}

// }

// func compileCode(filePath, outputPath string) error {
// 	cmd := exec.Command("go", "build", ".", "-o")

// }
