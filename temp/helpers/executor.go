package helpers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ishu17077/code_runner_backend/models"
	"github.com/ishu17077/code_runner_backend/services"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
)

var imageMap = map[string]string{
	"python":     "python:3.10-slim",
	"javascript": "node:18-alpine",
} //! just for now

func ExecuteWithInput(language, code, stdin string) (*models.ExecutionResult, error) {
	//TODO: Change timeout to a reasonable amount
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second)
	defer cancel()

	jobNameSpace := os.Getenv("JOB_NAMESPACE")
	if jobNameSpace == "" {
		jobNameSpace = "default"
	}

	configMapName := fmt.Sprintf("submission-%s", uuid.NewUUID())
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: jobNameSpace,
		},
		//! To be changed to whatever you like
		Data: map[string]string{"main.py": "print(\"Hiii\")"},
	}

	_, err := services.K8sClient.CoreV1().ConfigMaps(jobNameSpace).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to create config %s", err.Error())
	}
	cmd := []string{"python", filepath.Join("/app", "main.py")}
	defer services.K8sClient.CoreV1().ConfigMaps(jobNameSpace).Delete(context.Background(), configMapName, metav1.DeleteOptions{})
	//! To be changed
	var imageName = imageMap["python"]
	podName := fmt.Sprintf("job-%s", uuid.NewUUID())
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: jobNameSpace,
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name:    "runner",
					Image:   imageName,
					Command: cmd,
					Stdin:   true,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "code-volume",
							MountPath: "/app",
							ReadOnly:  true,
						},
					},
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("1"),
							corev1.ResourceMemory: resource.MustParse("256Mi"),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("0.5"),
							corev1.ResourceMemory: resource.MustParse("128Mi"),
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "code-volume",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: configMapName,
							},
						},
					},
				},
			},
			ActiveDeadlineSeconds: &[]int64{10}[0],
			SecurityContext: &corev1.PodSecurityContext{
				RunAsNonRoot: &[]bool{true}[0],
				RunAsUser:    &[]int64{1000}[0],
			},
			AutomountServiceAccountToken: &[]bool{false}[0],
		},
	}
	_, err = services.K8sClient.CoreV1().Pods(jobNameSpace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to create pod: %w", err)
	}
	defer services.K8sClient.CoreV1().Pods(jobNameSpace).Delete(context.Background(), podName, metav1.DeleteOptions{})
	if err := waitForPodRunning(ctx, podName, jobNameSpace); err != nil {
		return nil, err
	}

	attachReq := services.K8sClient.CoreV1().
		RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(jobNameSpace).
		SubResource("attach").
		VersionedParams(&corev1.PodAttachOptions{
			Stdin:  true,
			Stdout: true,
			Stderr: true,
		}, metav1.ParameterCodec)

	stream, err := attachReq.Stream(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to attach stream to pod: %w", err)
	}
	defer stream.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	readDone := make(chan error)
	writeDone := make(chan error)
	go func() {
		//? 1-> stdout , 2-> stderr
		_, err := demuxK8sStream(&stdoutBuf, &stderrBuf, stream)
		readDone <- err
	}()

	go func() {
		stdinWriter, ok := stream.(io.Writer)
		if !ok {
			writeDone <- fmt.Errorf("stream does not implement io.Writer")
			return
		}
		_, err := io.Copy(stdinWriter, strings.NewReader(stdin))
		writeDone <- err
	}()

	if err := waitForPodCompletion(ctx, podName, jobNameSpace); err != nil {
		return &models.ExecutionResult{
			Status: "Time Limit Exceeded: 10s",
			Stdout: "",
			Stderr: "",
		}, err
	}

	if err := <-readDone; err != nil && err != io.EOF {
		//! Implement logging for submission id if for any canse fails reasons beyond our control
		log.Printf("StdCopy error: %s for submission id", err)
	}

	finalPod, err := services.K8sClient.CoreV1().Pods(jobNameSpace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("Could not get final pod status: %w", err)
	}
	result := &models.ExecutionResult{
		Stdout: stdoutBuf.String(),
		Stderr: stderrBuf.String(),
	}

	if finalPod.Status.Phase == corev1.PodFailed {
		result.Status = "Runtime Error"
	} else {
		result.Status = "Completed"
	}

	return result, nil

}

func waitForPodRunning(ctx context.Context, podName, jobNameSpace string) error {
	watcher, err := services.K8sClient.CoreV1().Pods(jobNameSpace).Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", podName),
	})
	if err != nil {
		return fmt.Errorf("Failed to watch pod: %w", err)
	}
	defer watcher.Stop()

	for {
		select {
		case event := <-watcher.ResultChan():
			p, ok := event.Object.(*corev1.Pod)
			if !ok {
				continue
			}
			if p.Status.Phase == corev1.PodRunning {
				return nil
			}
			if p.Status.Phase == corev1.PodFailed || p.Status.Phase == corev1.PodSucceeded {
				return fmt.Errorf("pod finished before attach: %s", p.Status.Reason)
			}
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for pod to run: %w", ctx.Err())
		}
	}
}

func waitForPodCompletion(ctx context.Context, podName, namespace string) error {
	for {
		p, err := services.K8sClient.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if p.Status.Phase == corev1.PodSucceeded || p.Status.Phase == corev1.PodFailed {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func demuxK8sStream(stdout, stderr io.Writer, reader io.Reader) (int64, error) {
	// K8s stream format: [stream_id, 0, 0, 0, ...payload...] -> header
	// stream_id: 1=stdout, 2=stderr, stream_id change accordingly
	header := make([]byte, 8)
	var written int64
	for {
		n, err := io.ReadFull(reader, header)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return written, nil
		}
		if err != nil {
			return written, fmt.Errorf("Error reading stream header: %w", err)
		}
		if n < 8 {
			return written, fmt.Errorf("Invalid stream header")
		}
		// + 8 extra bytes performed by leftshift
		size := int64(header[4]) | int64(header[5])<<8 | int64(header[6])<<16 | int64(header[7])<<24

		if size == 0 {
			continue
		}

		dest := stdout
		if header[0] == 2 {
			dest = stderr
		}
		nbytes, err := io.CopyN(dest, reader, size)
		if err != nil {
			return written, fmt.Errorf("error copying payload: %w", err)
		}
		written += nbytes

	}

}
