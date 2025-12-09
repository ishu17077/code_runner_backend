package k8s

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/ishu17077/code_runner_backend/worker-node/models"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (k *K8sManager) RunOnPod(submission models.Submission, testCases []models.TestCase) (string, error) {
	//TODO: return models.Result instead of string
	executorPod, err := k.findRandomWarmPod()
	if err != nil {
		return "", err
	}

	fmt.Printf("Executing in pod: %s\n", executorPod)

	var payload models.Payload = models.Payload{
		Submission: submission,
		TestCases:  testCases,
	}

	stdinPayload, err := json.Marshal(payload)

	stdinPayload = []byte(base64.StdEncoding.EncodeToString(stdinPayload))
	if err != nil {
		return "", fmt.Errorf("Conversion to JSON failed: %s", err.Error())
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
		return "", fmt.Errorf("execution failed: %v | stderr: %s", execErr, stderr)
	}
	return output, nil
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

	return targetPod.Name, nil
}

func (k *K8sManager) execInPod(podName string, cmd []string, stdin []byte) (string, string, error) {
	req := k.clientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace("default").
		SubResource("exec")

	option := &corev1.PodExecOptions{
		Command:   cmd,
		Container: "code-runner",
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
	err = exec.Stream(streamOpts)

	return stdout.String(), stderr.String(), err
}

// escapeForShell safeguards string inputs against shell injection.
func (k *K8sManager) escapeForShell(s string) string {
	return strings.ReplaceAll(s, "'", "'\\''")
}

func (k *K8sManager) ptr(i int64) *int64 {
	return &i
}
