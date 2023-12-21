package service

import (
	"errors"

	"fmt"

	"encoding/json"

	"kube-demo-fe/goland/kubea-demo/config"

	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var K8S k8s

type k8s struct {
	//提供多集群 client
	ClientMap map[string]*kubernetes.Clientset
	// 提供集群列表功能
	kubeConfMap map[string]string
}

// 根据集群名称获取 client
func (k *k8s) GetClient(clusterName string) (*kubernetes.Clientset, error) {
	client, ok := k.ClientMap[clusterName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("集群:%s 不存在，无法获取 client\n"))
	}
	return client, nil
}

//初始化 client

func (k *k8s) Init() {
	mp := make(map[string]string)
	k.ClientMap = make(map[string]*kubernetes.Clientset, 0)
	//反序列化
	if err := json.Unmarshal([]byte(config.Kubeconfigs), &mp); err != nil {
		panic(err)
	}
	k.kubeConfMap = mp

	//初始化 client
	for clusterName, kubeconfig := range mp {
		logger.Info(clusterName, kubeconfig)
		conf, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			logger.Warn(fmt.Sprintf("获取 %s;创建配置失败%s", clusterName, err.Error()))
		}
		clientSet, err := kubernetes.NewForConfig(conf)
		if err != nil {
			logger.Warn(fmt.Sprintf("获取 %s;创建 clientset 失败%s", clusterName, err.Error()))
		}
		k.ClientMap[clusterName] = clientSet
		logger.Info(fmt.Sprintf("集群:%s: 获取client 成功", clusterName))
	}

}
