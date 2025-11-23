package helpers

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/ishu17077/code_runner_backend/models"
	"github.com/ishu17077/code_runner_backend/services"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var rabbitMQ = services.QueuerChannel
var submissionsCollection = services.OpenCollection(services.MongoClient, "submissions")
var updateOptions = options.UpdateOne().SetUpsert(true)

func ExecuteJobs() {
	msgs, err := rabbitMQ.Consume(
		"jobs",
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Failed to register a consumer/ listener %s", err.Error())
		return
	}
	for i := 0; i < 10; i++ {
		go func(workerId int) {
			for d := range msgs {
				log.Printf("Worker %d received a job.", workerId)
				var job models.Job
				json.Unmarshal(d.Body, &job)

				err := processSubmission(job)

				if err != nil {
					log.Printf("Failed to process job %s: %s", job.SubmissionID, err)
					// Nack the message to requeue it
					d.Nack(false, true)
				} else {
					d.Ack(false)
				}

			}
		}(i)
	}
}

func processSubmission(job models.Job) error {
	var submission models.Submission
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "submission_id", Value: job.SubmissionID}}}, {Key: "$limit", Value: 1}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "tests"},
		{Key: "localField", Value: "problem_id"},
		{Key: "foreignField", Value: "problem_id"},
		{Key: "as", Value: "tests"},
	}}}
	// mongoDBFixArrayStage := bson.D{{Key: "$set", Value: bson.D{{Key: "tests", Value: bson.D{{Key: "$arrayElemAt", Value: []any{"$tests", 0}}}}}}}

	var tests []models.Test

	type AggregationResult struct {
		models.Submission `bson:",inline"` // Embeds all fields from Submission
		Tests             []models.Test    `bson:"tests"` // Catches the 'tests' array from $lookup
	}

	result, err := submissionsCollection.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage})
	if err != nil {
		log.Printf("Unable to process submission id: %s", job.SubmissionID)
		return err
	}
	defer result.Close(ctx)
	// if err != mongo.ErrNoDocuments {
	// 	log.Printf("Cannot process the submission id: %s", job.SubmissionID)
	// 	return nil
	// }
	if !result.Next(ctx) {
		// Check if the cursor is empty
		if err := result.Err(); err != nil {
			log.Printf("Cursor error: %v", err)
			return err
		}
		// No document was found
		log.Printf("No submission found with id: %s", job.SubmissionID)
		return mongo.ErrNoDocuments
	}
	var aggResult AggregationResult

	if err := result.Decode(&aggResult); err != nil {
		log.Printf("Decoding error: %s", err.Error())
		return err
	}

	submission = aggResult.Submission
	tests = aggResult.Tests
	log.Printf("Executing submission for Submission Id: %s, for %d tests", submission.Submission_id, len(tests))
	for _, test := range tests {
		log.Printf("Executing test %s for Submission Id: %s", test.Test_id, submission.Submission_id)

		execResults, err := ExecuteWithInput(
			submission.Language,
			submission.Code,
			test.Stdin,
		)

		if err != nil {
			if err == context.DeadlineExceeded {
				updateStatus(submission.Submission_id, "Time Limit Exceeded: 10s", "", err.Error(), "FAILED")
				return err
			}
			updateStatus(submission.Submission_id, "Error Occured", "", err.Error(), "FAILED")
			return nil
		}
		if execResults.Stderr != "" {
			updateStatus(submission.Submission_id, "Runtime Error", "", execResults.Stderr, "FAILED")
			return nil
		}

		if !compareOutput(execResults.Stdout, test.ExpectedOutput) {
			updateStatus(aggResult.Submission_id, "Wrong Answer", execResults.Stdout, execResults.Stderr, "FAILED")
			return nil
		}

		updateStatus(submission.Submission_id, "Accepted", execResults.Stdout, execResults.Stderr, "SUCCESS")
	}
	return nil
}

func updateStatus(submission_id, message, stdout, stderr, currStatus string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var filter = bson.M{"submission_id": submission_id}
	var status models.Status = models.Status{
		Current_status: currStatus,
		Message:        message,
		Stdout:         stdout,
		Stderr:         stderr,
		Completed_At:   time.Now(),
	}
	defer cancel()
	var updateObj = bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: status}}}}
	_, err := submissionsCollection.UpdateOne(ctx, filter, updateObj, updateOptions)
	if err != nil {
		log.Printf("Error updating status for submission_id: %s", err.Error())
	}
}

func compareOutput(stdout, expected string) bool {
	return strings.TrimSpace(stdout) == strings.TrimSpace(expected)
}
