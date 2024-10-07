package kubernetes

import (
	"context"
	"fmt"
	"io"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"strings"
	"time"

	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// IKubernetesClient defines the interface for interacting with Kubernetes Jobs.
type IKubernetesClient interface {
	CreateJob(imageName string, cmd []string, env []string) (string, error)
	WaitForJobCompletion(jobName string) error
	GetJobLogs(jobName string) (string, error)
	DeleteJob(jobName string) error
}

// KubernetesClient is a client for creating and managing Kubernetes Jobs.
type KubernetesClient struct {
	clientset *kubernetes.Clientset
	namespace string
}

// NewKubernetesClient creates a new Kubernetes client.
func NewKubernetesClient(namespace string) (*KubernetesClient, error) {
	// Attempt to use in-cluster config (for running inside a cluster)
	config, err := rest.InClusterConfig()
	if err != nil {
		// Use the local kubeconfig when running locally
		kubeconfig := filepath.Join(homeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to get Kubernetes config: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	return &KubernetesClient{clientset: clientset, namespace: namespace}, nil
}

// homeDir returns the home directory for the executing user.
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // Windows compatibility
}

// CreateJob creates a Kubernetes Job to execute the given command in the provided container image.
func (k *KubernetesClient) CreateJob(imageName string, cmd []string, env []string) (string, error) {
	jobName := fmt.Sprintf("code-job-%d", time.Now().Unix())

	job := &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
			Labels: map[string]string{
				"frontend": "code-execution",
			},
		},
		Spec: v1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:    "code-runner",
							Image:   imageName,
							Command: cmd,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "code-volume",
									MountPath: "/code",
								},
							},
							Env: createEnvVars(env),
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("100m"),
									corev1.ResourceMemory: resource.MustParse("128Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("500m"),
									corev1.ResourceMemory: resource.MustParse("256Mi"),
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "code-volume",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}

	_, err := k.clientset.BatchV1().Jobs(k.namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create Kubernetes Job: %v", err)
	}

	return jobName, nil
}

// WaitForJobCompletion waits for the Kubernetes Job to complete successfully.
func (k *KubernetesClient) WaitForJobCompletion(jobName string) error {
	timeout := time.After(10 * time.Minute) // Set a timeout to avoid indefinite waiting
	for {
		select {
		case <-timeout:
			return fmt.Errorf("timed out waiting for job %s to complete", jobName)
		default:
			job, err := k.clientset.BatchV1().Jobs(k.namespace).Get(context.TODO(), jobName, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get Job status: %v", err)
			}

			if job.Status.Succeeded > 0 {
				return nil // Job completed successfully
			}

			if job.Status.Failed > 0 {
				return fmt.Errorf("job %s failed", jobName)
			}

			time.Sleep(2 * time.Second) // Poll every 2 seconds
		}
	}
}

// GetJobLogs retrieves the logs from the Pod created by the Job.
func (k *KubernetesClient) GetJobLogs(jobName string) (string, error) {
	// List Pods created by the Job, using job-name as the label selector.
	podList, err := k.clientset.CoreV1().Pods(k.namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", jobName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get Pod list for Job %s: %v", jobName, err)
	}

	if len(podList.Items) == 0 {
		return "", fmt.Errorf("no Pods found for Job %s", jobName)
	}

	// Get logs from the first Pod
	podName := podList.Items[0].Name
	logs, err := k.clientset.CoreV1().Pods(k.namespace).GetLogs(podName, &corev1.PodLogOptions{}).Stream(context.TODO())
	if err != nil {
		return "", fmt.Errorf("failed to get logs for Pod %s: %v", podName, err)
	}
	defer logs.Close()

	var logData strings.Builder
	_, err = io.Copy(&logData, logs)
	if err != nil {
		return "", fmt.Errorf("failed to read logs from Pod %s: %v", podName, err)
	}

	return logData.String(), nil
}

// DeleteJob deletes the Job and its associated Pods.
func (k *KubernetesClient) DeleteJob(jobName string) error {
	deletePolicy := metav1.DeletePropagationForeground
	err := k.clientset.BatchV1().Jobs(k.namespace).Delete(context.TODO(), jobName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		return fmt.Errorf("failed to delete Job %s: %v", jobName, err)
	}
	return nil
}

// createEnvVars creates a list of Kubernetes environment variables from a string slice.
func createEnvVars(env []string) []corev1.EnvVar {
	envVars := []corev1.EnvVar{}
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			envVars = append(envVars, corev1.EnvVar{
				Name:  parts[0],
				Value: parts[1],
			})
		}
	}
	return envVars
}
