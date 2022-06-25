package helpers

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

var (
	InvalidNsm = []corev1.NodeSelectorRequirement{
		{
			Key:      "os",
			Values:   []string{"LINUX"},
			Operator: "PIWPIW",
		},
	}
	EmptyNsm = []corev1.NodeSelectorRequirement{}
	Node     = getBasicNode()
)

func TestValidNodeSelectorRequirementsAsSelector(t *testing.T) {
	inOperators := []corev1.NodeSelectorOperator{
		corev1.NodeSelectorOpIn,
		corev1.NodeSelectorOpNotIn,
	}
	existOperatos := []corev1.NodeSelectorOperator{
		corev1.NodeSelectorOpExists,
		// Excluded for now
		// corev1.NodeSelectorOpDoesNotExist,
	}
	comparatifOperators := []corev1.NodeSelectorOperator{
		corev1.NodeSelectorOpGt,
		corev1.NodeSelectorOpLt,
	}
	for _, op := range inOperators {
		t.Run(string(op), func(t *testing.T) {
			validNsm := []corev1.NodeSelectorRequirement{
				{
					Key:      "os",
					Operator: op,
					Values:   []string{"LINUX"},
				},
			}
			selector, err := NodeSelectorRequirementsAsSelector(validNsm)
			assertLabelsAndRequirementsEquality(t, selector, validNsm, err, true)
		})
	}
	for _, op := range existOperatos {
		t.Run(string(op), func(t *testing.T) {
			validNsm := []corev1.NodeSelectorRequirement{
				{
					Key:      "ip-address",
					Operator: op,
				},
			}
			selector, err := NodeSelectorRequirementsAsSelector(validNsm)
			assertLabelsAndRequirementsEquality(t, selector, validNsm, err, false)
		})
	}
	for _, op := range comparatifOperators {
		t.Run(string(op), func(t *testing.T) {
			validNsm := []corev1.NodeSelectorRequirement{
				{
					Key:      "number-of-years",
					Operator: op,
					Values:   []string{"2"},
				},
			}
			selector, err := NodeSelectorRequirementsAsSelector(validNsm)
			assertLabelsAndRequirementsEquality(t, selector, validNsm, err, true)
		})
	}
}

func assertLabelsAndRequirementsEquality(t *testing.T, selector labels.Selector, nsm []corev1.NodeSelectorRequirement, err error, valuesExist bool) {
	assert.NoError(t, err)
	requirements, _ := selector.Requirements()
	keysEquality := strings.EqualFold(requirements[0].Key(), nsm[0].Key)
	operatorsEquality := strings.EqualFold(string(requirements[0].Operator()), string(nsm[0].Operator))
	fmt.Println(string(requirements[0].Operator()))
	fmt.Println(string(nsm[0].Operator))
	if valuesExist == true {
		valuesEquality := reflect.DeepEqual(requirements[0].Values().List(), nsm[0].Values)
		assert.True(t, valuesEquality)
	}
	assert.True(t, keysEquality)
	assert.True(t, operatorsEquality)
}

func TestInvalidNodeSelectorRequirementsAsSelector(t *testing.T) {
	labels, err := NodeSelectorRequirementsAsSelector(InvalidNsm)
	assert.Error(t, err)
	assert.Nil(t, labels)
}

func TestEmptyNodeSelectorRequirementsAsSelector(t *testing.T) {
	res, err := NodeSelectorRequirementsAsSelector(EmptyNsm)

	assert.NoError(t, err)
	assert.Equal(t, res, labels.Nothing())
}

func getBasicNode() *corev1.Node {
	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "basic-node",
			Labels: map[string]string{
				"os":              "linux",
				"number-of-years": "2",
				"ip-address":      "127.0.0.1",
			},
		},
		Spec: corev1.NodeSpec{
			Taints: []corev1.Taint{
				{
					Key:    "key1",
					Value:  "value1",
					Effect: corev1.TaintEffectNoExecute,
				},
			},
		},
	}
}

func TestNodesMatchingWithInvalidSelector(t *testing.T) {
	invalidNodeSelectorTerms := []corev1.NodeSelectorTerm{
		{
			MatchExpressions: InvalidNsm,
		},
	}
	isMatch := NodeMatchesNodeSelectorTerms(Node, invalidNodeSelectorTerms)
	assert.False(t, isMatch)
}
func TestNodesMatchingWithEmptyTerms(t *testing.T) {
	invalidNodeSelectorTerms := []corev1.NodeSelectorTerm{}
	isMatch := NodeMatchesNodeSelectorTerms(Node, invalidNodeSelectorTerms)
	assert.False(t, isMatch)
}

func TestNodesMatchingWithValidTerms(t *testing.T) {
	matchingRequirements := []corev1.NodeSelectorTerm{
		{
			MatchExpressions: []corev1.NodeSelectorRequirement{
				{
					Key:      "os",
					Values:   []string{"linux"},
					Operator: corev1.NodeSelectorOpIn,
				},
			},
		},
		{
			MatchExpressions: []corev1.NodeSelectorRequirement{
				{
					Key:      "ip-address",
					Values:   []string{"172.0.10.1"},
					Operator: corev1.NodeSelectorOpIn,
				},
			},
		},
	}
	isMatch := NodeMatchesNodeSelectorTerms(Node, matchingRequirements)
	assert.True(t, isMatch)

	nonMatchingRequirements := []corev1.NodeSelectorTerm{
		{
			MatchExpressions: []corev1.NodeSelectorRequirement{
				{
					Key:      "os",
					Values:   []string{"linux"},
					Operator: corev1.NodeSelectorOpIn,
				},
				{
					Key:      "ip-address",
					Values:   []string{"172.0.10.1"},
					Operator: corev1.NodeSelectorOpIn,
				},
			},
		},
		{
			MatchExpressions: []corev1.NodeSelectorRequirement{
				{
					Key:      "ip-address",
					Values:   []string{"172.0.10.1"},
					Operator: corev1.NodeSelectorOpIn,
				},
			},
		},
	}

	isMatch = NodeMatchesNodeSelectorTerms(Node, nonMatchingRequirements)
	assert.False(t, isMatch)
}

// func TestNodeListsMerging(t *testing.T) {

// }
