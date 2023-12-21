package service

import (
	"context"
	"fmt"

	"errors"

	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Pod pod

type pod struct{}

//获取 pod 列表 ，定义列表的返回类型

type PodsResp struct {
	Items []corev1.Pod `json:"items"`
	Total int          `json:"total"`
}

// 获取 pod 列表
// client 用于选择哪个集群
func (p *pod) GetPods(client *kubernetes.Clientset, fileName, nameSpace string, limit, page int) (podsResp *PodsResp, err error) {
	// context.TODO()用于声明一个空的context上下文，用于 List 方法内置这个请求的超时，这里常用方法
	// metav1.ListOptions{} 用于管理 list 数据，如使用 label，field 等
	pods, err := client.CoreV1().Pods(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取 pod列表失败，%v", err))
		return nil, errors.New(fmt.Sprintf("获取 pod列表失败，%v\n", err))
	}
	//
	selectData := &dataSelector{
		GenericDataList: p.toCells(pods.Items),
		dataSelectorQuery: &DataSelectorQuery{
			FilterQuery: &FilterQuery{Name: fileName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	filtered := selectData.Filter()
	total := len(filtered.GenericDataList)

	data := filtered.Sort().Paginate()
	// 将【】dataCell 类型的 pod 列表转成v1.pod 列表
	podList := p.fromCells(data.GenericDataList)
	return &PodsResp{Items: podList, Total: total}, nil
}

// 定义 DataCell 到 pod 类型转换方法
func (p *pod) toCells(std []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = podCell(std[i])
	}
	return cells
}

func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}
