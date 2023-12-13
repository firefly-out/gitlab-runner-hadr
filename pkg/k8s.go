package pkg

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// CreateK8sClient creates a *kubernetes.Clientset to access the kubernetes API from inside the cluster.
// Please make sure the pod has a RoleBinding with the ability to read pods.
func CreateK8sClient() (client *kubernetes.Clientset, err error) {
	// Use the in-cluster configuration.
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("Error creating in-cluster configuration: ", err.Error())
		return nil, err
	}

	// Create a Kubernetes client.
	client, err = kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error creating Kubernetes client: ", err.Error())
		return nil, err
	}

	return client, nil
}

// ListPods listing pods under the received deployment in the given namespace.
func ListPods(client *kubernetes.Clientset, namespace, deployment string) (*v1.PodList, error) {
	// Set up a context
	ctx := context.TODO()

	// Get the pods associated with the deployment.
	podList, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", deployment),
	})
	if err != nil {
		fmt.Println("Error getting pods: ", err.Error())
		return nil, err
	}

	// Print pod names.
	fmt.Println("Pods in deployment", deployment, "are:")
	for _, pod := range podList.Items {
		fmt.Println("-", pod.Name)
	}

	return podList, nil
}
