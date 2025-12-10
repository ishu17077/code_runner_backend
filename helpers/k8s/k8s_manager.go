package k8s

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/ishu17077/code_runner_backend/models"
	currentstatus "github.com/ishu17077/code_runner_backend/models/enums/current_status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type K8sManager struct {
	clientSet *kubernetes.Clientset
	config    *rest.Config
}

var K8sMgr *K8sManager = NewK8sManager()

func NewK8sManager() *K8sManager {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to get in-cluster config: %v", err))
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create clientset: %v", err))
	}

	return &K8sManager{
		clientSet: clientSet,
		config:    config,
	}
}

func (k *K8sManager) RunOnPod(submission models.Submission) (models.Result, error) {
	//TODO: return models.Result instead of string
	executorPod, err := k.findRandomWarmPod()
	if err != nil {
		return emptyResult("Unable to execute the program."), err
	}

	fmt.Printf("Executing in pod: %s\n", executorPod)

	stdinPayload, err := json.Marshal(submission)

	stdinPayload = []byte(base64.StdEncoding.EncodeToString(stdinPayload))
	if err != nil {
		return emptyResult("Unable to execute the program."), fmt.Errorf("Conversion to JSON failed: %s", err.Error())
	}
	cmd := []string{"./runner"}

	output, stderr, execErr := k.execInPod(executorPod, cmd, stdinPayload)

	cleanupErr := k.clientSet.CoreV1().Pods("default").Delete(context.Background(), executorPod, metav1.DeleteOptions{
		GracePeriodSeconds: k.ptr(0),
	})

	if cleanupErr != nil {
		fmt.Printf("Warning: Failed to delete pod %s: %v\n", executorPod, cleanupErr)
	}

	if execErr != nil {
		if errors.Is(execErr, context.DeadlineExceeded) {
			return emptyResult("Time Limit Exceeded"), fmt.Errorf("Unable to complete execution")
		}
		return emptyResult("Memory Limit Exceeded"), fmt.Errorf("execution failed: %v | stderr: %s", execErr, stderr)
	}
	return extractJsonFromStdout(output), nil
}

func (k *K8sManager) findRandomWarmPod() (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pods, err := k.clientSet.CoreV1().Pods("default").List(ctx, metav1.ListOptions{
		//? warm-runner to be changed to language specific code
		LabelSelector: "app=warm-runner, status=ready",
		FieldSelector: "status.phase=Running",
	})

	if err != nil {
		return "", fmt.Errorf("Failed to list pods: %w", err)
	}

	if len(pods.Items) == 0 {
		return "", fmt.Errorf("no warm pods available in pool")
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	targetPod := pods.Items[rand.Intn(len(pods.Items))]
	if targetPod.DeletionTimestamp != nil {
		return "", fmt.Errorf("selected pod is terminating, retrying...")
	}
	payload := []byte(`{"metadata":{"labels":{"status":"busy"}}}`)
	_, err = k.clientSet.CoreV1().Pods("default").Patch(
		ctx,
		targetPod.Name,
		types.MergePatchType,
		payload,
		metav1.PatchOptions{},
	)
	if err != nil {
		fmt.Printf("Error claiming pod: %s", targetPod.Name)
	}
	return targetPod.Name, nil
}

func (k *K8sManager) execInPod(podName string, cmd []string, stdin []byte) (string, string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	req := k.clientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace("default").
		SubResource("exec")

	option := &corev1.PodExecOptions{
		Command:   cmd,
		Container: "warm-runner",
		Stdin:     len(stdin) > 0,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}
	req.VersionedParams(option, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(k.config, "POST", req.URL())
	if err != nil {
		return "", "", fmt.Errorf("Failed to init code: %w", err)
	}

	var stdout, stderr bytes.Buffer

	streamOpts := remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	}
	if len(stdin) > 0 {
		streamOpts.Stdin = bytes.NewReader(stdin)
	}
	//* Add context timeout option is available
	err = exec.StreamWithContext(ctx, streamOpts)

	return stdout.String(), stderr.String(), err
}

func extractJsonFromStdout(res string) models.Result {

	var regExpMatch = regexp.MustCompile(`(?s)---JSON_START---(.*?)---JSON_END---`)
	matches := regExpMatch.FindStringSubmatch(res)

	if len(matches) < 2 {
		return emptyResult("Error consuming too much resources")
	}
	res = strings.TrimSpace(matches[1])

	var result models.Result
	jsonData, err := base64.StdEncoding.DecodeString(res)
	if err != nil {
		return emptyResult("Unable to execute the program.")

	}
	if err = json.Unmarshal(jsonData, &result); err != nil {
		return emptyResult("Unable to execute the program.")
	}

	return result

}

func emptyResult(err string) models.Result {
	return models.Result{
		Status:  currentstatus.FAILED.ToString(),
		Results: []models.ExecResult{},
		Error:   err,
	}
}

// escapeForShell safeguards string inputs against shell injection.
// func (k *K8sManager) escapeForShell(s string) string {
// 	return strings.ReplaceAll(s, "'", "'\\''")
// }

func (k *K8sManager) ptr(i int64) *int64 {
	return &i
}
