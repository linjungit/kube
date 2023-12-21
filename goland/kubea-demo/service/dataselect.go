package service

import (
	"strings"
	"time"

	"sort"

	corev1 "k8s.io/api/core/v1"
)

type dataSelector struct {
	GenericDataList   []DataCell
	dataSelectorQuery *DataSelectorQuery
}

// 接口 用于各种资源 list 的类型转换，转换后可以使用 dataSelector 的排序、过滤 分页
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

// DataSelectorQuery 定义过滤和分页的属性，过滤：name，分页 limit和 page
type DataSelectorQuery struct {
	FilterQuery   *FilterQuery
	PaginateQuery *PaginateQuery
}

type FilterQuery struct {
	Name string
}

type PaginateQuery struct {
	Limit int
	Page  int
}

// 排序，实现自定义结构的排序，需要重写 len \swap\less
// len 方法用于获取数组长度
func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

// Swap 方法用于数组中元素在比较大小后怎么交换位置，可定义升降序
// i,j 是切片的下标
func (d *dataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

// Less 方法用于比较两个元素的大小，
// 返回 true 表示 i 小于 j，false 表示 i 大于 j
func (d *dataSelector) Less(i, j int) bool {
	return d.GenericDataList[i].GetCreation().Before(d.GenericDataList[j].GetCreation())
}

func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

// 过滤 Filter 方法用于过滤元素，比较元素的 name 属性。若包含则返回
func (d *dataSelector) Filter() *dataSelector {

	if d.dataSelectorQuery.FilterQuery.Name == "" {
		return d
	}
	// 若 name 的传参不为空，则返回元素中包含 name 的所有元素
	filteredList := []DataCell{}
	for _, value := range d.GenericDataList {
		matched := true
		objName := value.GetName()
		//如果 objname 不包含name
		if !strings.Contains(objName, d.dataSelectorQuery.FilterQuery.Name) {
			// 就是不匹配的，进入下一轮的循环
			matched = false
			continue
		}
		if matched {
			filteredList = append(filteredList, value)
		}
	}
	d.GenericDataList = filteredList
	return d
}

// Paginate方法用于数组分页，根据 limit 和 Page 的传参，返回数据
func (d *dataSelector) Paginate() *dataSelector {
	limit := d.dataSelectorQuery.PaginateQuery.Limit //10
	page := d.dataSelectorQuery.PaginateQuery.Page   //2

	//验证参数合法，若参数不合法，则返回所有数据
	if limit <= 0 || page <= 0 {
		return d
	}
	// 计算开始下标
	// 举例 25个元素的切片，limit10
	// page1  start0 end 10
	// page2  start10 end 20
	startIndex := limit * (page - 1) //10
	endIndex := limit * page         //20
	if len(d.GenericDataList) < endIndex {
		endIndex = len(d.GenericDataList)
	}
	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d
}

// 定义 podCell 类型，实现两个方法，可以进行类型转换 getName 可类型转换
type podCell corev1.Pod // 可以增加方法

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}
