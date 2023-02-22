package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	//获取kubeconfig配置文件地址
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse() 
	//buildkubeconfig从给定的kubeconfig路径构建k8s的clientset
	//如果kubeconfig为空则返回一个用于在集群内工作的clientset
	var config *rest.Config
	var err error
	if *kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		panic(err)
	}
	kubeClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	//这个namespace啥意思wc
	var pod *v1.Pod
	pod = &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sjs-1",
			Namespace: "default",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "web",
					Image: "nginx:1.12",
					Ports: []v1.ContainerPort{
						{
							Name:          "http",
							Protocol:      v1.ProtocolTCP,
							ContainerPort: 80,
						},
					},
				},
			},
		},
	}
	pod, err = kubeClientSet.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})

	if err != nil {
		panic(err)
	}
	fmt.Println("pod 的相关信息为 %v", pod)

}
