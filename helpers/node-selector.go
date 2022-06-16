package helpers

import(
	"k8s.io/api/core/v1"
)
func HandleMatchExpressions(matchExpressions []v1.NodeSelectorRequirement) ([]int) {
	return []int{1}
}
func HandleMatchFields(matchFields []v1.NodeSelectorRequirement) ([]int) {
	return []int{1}
}