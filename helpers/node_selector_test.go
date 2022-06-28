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
	// We test whether we obtain the same semantic when transforming from Node Selector Requirement to Labels Selectors
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
	/* In order to make sure that the semantic is the same we should simply compare
	the labels' keys values and operators with node selector requirement ones
	*/
	assert.NoError(t, err)
	requirements, _ := selector.Requirements()
	keysEquality := strings.EqualFold(requirements[0].Key(), nsm[0].Key)
	operatorsEquality := strings.EqualFold(string(requirements[0].Operator()), string(nsm[0].Operator))
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
	// case1: One requirement matches and the other is not
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

	// case 2: Both requirement do not match even if in one of them there is an expression that matches
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

func TestNodeListsMerging(t *testing.T) {
	// case 1: init not empty, filtered empty
	initNodes := corev1.NodeList{
		Items: []corev1.Node{
			*Node,
		},
	}

	filteredNodes := corev1.NodeList{}

	res := MergeNodes(initNodes, filteredNodes)
	assert.Equal(t, len(res.Items), 1)

	// case 2: init not empty, filtered not empty but the same node as init
	filteredNodes = corev1.NodeList{
		Items: []corev1.Node{
			*Node,
		},
	}

	res = MergeNodes(initNodes, filteredNodes)
	assert.Equal(t, len(res.Items), 1)

	// case 3: init not empty, filtered not empty
	filteredNodes = corev1.NodeList{
		Items: []corev1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "second-node",
				},
			},
		},
	}

	res = MergeNodes(initNodes, filteredNodes)
	assert.Equal(t, len(res.Items), 2)

	// case 4: init empty, filtered not empty
	initNodes = corev1.NodeList{}
	res = MergeNodes(initNodes, filteredNodes)
	assert.Equal(t, len(res.Items), 1)

	// case 5: both empty
	filteredNodes = corev1.NodeList{}
	res = MergeNodes(initNodes, filteredNodes)
	assert.Equal(t, len(res.Items), 0)
}

func TestFilterNodesByRegex(t *testing.T) {
	nodes := *generateNodeListForRegexTest()

	filteredNodes := FilterByRegex(&nodes,
		[]string{"[a-zA-Z0-9]*minikube[a-zA-Z0-9]*"},
	)
	fmt.Print(filteredNodes.Items)
	expectedResult := &corev1.NodeList{
		Items: []corev1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "minikube1",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "2minikube2",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "3minikube3",
				},
			},
		},
	}

	assert.Equal(t, len(filteredNodes.Items), len(expectedResult.Items))
	assert.Equal(t, filteredNodes.Items[0].Name, expectedResult.Items[0].Name)
	assert.Equal(t, filteredNodes.Items[1].Name, expectedResult.Items[1].Name)
	assert.Equal(t, filteredNodes.Items[2].Name, expectedResult.Items[2].Name)
}
