package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := getKubeconfig()
	clientset, err := getClientset(kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	pods, err := listPods(clientset)
	if err != nil {
		panic(err.Error())
	}

	restartDatabasePods(clientset, pods)
}

func getKubeconfig() string {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	if _, err := os.Stat(*kubeconfig); os.IsNotExist(err) {
		fmt.Printf("Kubeconfig file not found: %s\n", *kubeconfig)
		os.Exit(1)
	}

	return *kubeconfig
}

func getClientset(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func listPods(clientset *kubernetes.Clientset) (*corev1.PodList, error) {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return pods, nil
}

func restartDatabasePods(clientset *kubernetes.Clientset, pods *corev1.PodList) {
	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, "database") {
			fmt.Printf("Restarting pod: %s\n", pod.Name)

			if podOwner := metav1.GetControllerOf(&pod); podOwner != nil {
				var err error
				switch podOwner.Kind {
				case "Deployment":
					err = rolloutRestartDeployment(clientset, pod.Namespace, podOwner.Name)
				case "StatefulSet":
					err = rolloutRestartStatefulSet(clientset, pod.Namespace, podOwner.Name)
				default:
					fmt.Printf("Skipping %s: unsupported controller kind %s\n", pod.Name, podOwner.Kind)
					continue
				}
				if err != nil {
					fmt.Printf("Error restarting %s: %v\n", pod.Name, err)
				}
			} else {
				fmt.Printf("Pod %s is not controlled by a deployment or statefulset\n", pod.Name)
			}
		}
	}
}

func rolloutRestartDeployment(clientset *kubernetes.Clientset, namespace, name string) error {
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	deployment, err := deploymentsClient.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = map[string]string{}
	}
	deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	_, err = deploymentsClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})
	return err
}

func rolloutRestartStatefulSet(clientset *kubernetes.Clientset, namespace, name string) error {
	statefulSetsClient := clientset.AppsV1().StatefulSets(namespace)
	statefulSet, err := statefulSetsClient.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if statefulSet.Spec.Template.Annotations == nil {
		statefulSet.Spec.Template.Annotations = map[string]string{}
	}
	statefulSet.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	_, err = statefulSetsClient.Update(context.TODO(), statefulSet, metav1.UpdateOptions{})
	return err
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}
