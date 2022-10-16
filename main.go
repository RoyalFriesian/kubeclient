package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	typev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	NSL int = 20 //namespace length for cmd ui
	PNL int = 40 //podName length for cmd ui
	CL  int = 10 //count length for cmd ui
	AL  int = 20 //age length for cmd ui
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "kube config file path - string")
	namespace := flag.String("namespace", "", "namespace - string")
	podName := flag.String("name", "", "pod name - string")
	podAge := flag.Int("age", 0, "pod age - Int")
	flag.Parse()
	if *kubeconfig != "" {
		output := getPods(*kubeconfig, *namespace, *podName, *podAge)
		if output == "" {
			fmt.Println("No result found")
		} else {
			fmt.Println(ms("NameSpace", NSL), ms("PodName", PNL), ms("Count", CL), ms("Age", AL))
			fmt.Println(output)
		}
	} else {
		fmt.Println("Error")
	}
}

func ms(data string, length int) string {
	diff := length - len([]rune(data))
	if diff > 0 {
		for diff != 0 {
			data = data + " "
			diff--
		}
	} else {
		minus := len(data) + (diff - 3)
		data = data[:minus]
		data = data + "..."
	}
	return data
}

func getPods(kubeconfig string, namespaceFilter string, podName string, podAge int) string {
	output := ""
	k8sClient, err := getClient(kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	svcs, err := getAllService(k8sClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
	found := make(map[string]bool)
	for _, svc := range svcs.Items {
		pods, err := getPodsForSvc(&svc, namespaceFilter, k8sClient)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(2)
		}
		count := 1
		for _, pv := range pods.Items {

			if !found[pv.ObjectMeta.Name] {
				age := time.Now().Sub(pv.Status.StartTime.Time)
				if podAge > 0 {
					if int(age.Seconds()) < podAge {
						continue
					}
				}
				if podName != "" {
					if podName == pv.ObjectMeta.Name {
						output += ms(pv.ObjectMeta.Namespace, NSL) + " " + ms(pv.ObjectMeta.Name, PNL) + " " + ms(strconv.Itoa(count), CL) + " " + ms(age.String(), AL) + "\n"
						break
					}
				} else {
					showTotalCount := ""
					if podAge <= 0 {
						showTotalCount = "/" + strconv.Itoa(len(pods.Items))
					} else {
						showTotalCount = ""
					}
					output += ms(pv.ObjectMeta.Namespace, NSL) + " " + ms(pv.ObjectMeta.Name, PNL) + " " + ms(strconv.Itoa(count)+showTotalCount, CL) + " " + ms(age.String(), AL) + "\n"
					count++
				}
				found[pv.ObjectMeta.Name] = true
			}
		}
	}
	return output
}

func getClient(configLocation string) (typev1.CoreV1Interface, error) {
	kubeconfig := filepath.Clean(configLocation)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1(), nil
}

func getAllService(k8sClient typev1.CoreV1Interface) (*corev1.ServiceList, error) {
	listOptions := metav1.ListOptions{}
	ctx := context.TODO()
	svcs, err := k8sClient.Services("").List(ctx, listOptions)
	if err != nil {
		log.Fatal(err)
	}
	return svcs, nil
}

func getPodsForSvc(svc *corev1.Service, namespace string, k8sClient typev1.CoreV1Interface) (*corev1.PodList, error) {
	set := labels.Set(svc.Spec.Selector)
	listOptions := metav1.ListOptions{LabelSelector: set.AsSelector().String()}
	ctx := context.TODO()
	pods, err := k8sClient.Pods(namespace).List(ctx, listOptions)
	/* for _, pod := range pods.Items {
		fmt.Fprintf(os.Stdout, "pod name: %v\n", pod.Name)
	} */
	return pods, err
}
