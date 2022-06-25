package helpers

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
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
		validNsm := []corev1.NodeSelectorRequirement{
			{
				Key:      "os",
				Operator: op,
				Values:   []string{"LINUX"},
			},
		}
		selector, err := NodeSelectorRequirementsAsSelector(validNsm)
		assertLabelsAndRequirementsEquality(t, selector, validNsm, err, true)
	}
	for _, op := range existOperatos {
		validNsm := []corev1.NodeSelectorRequirement{
			{
				Key:      "ip-address",
				Operator: op,
			},
		}
		selector, err := NodeSelectorRequirementsAsSelector(validNsm)
		assertLabelsAndRequirementsEquality(t, selector, validNsm, err, false)
	}
	for _, op := range comparatifOperators {
		validNsm := []corev1.NodeSelectorRequirement{
			{
				Key:      "number-of-years",
				Operator: op,
				Values:   []string{"2"},
			},
		}
		selector, err := NodeSelectorRequirementsAsSelector(validNsm)
		assertLabelsAndRequirementsEquality(t, selector, validNsm, err, true)
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
	invalidNsm := []corev1.NodeSelectorRequirement{
		{
			Key:      "os",
			Values:   []string{"LINUX"},
			Operator: "PIWPIW",
		},
	}
	labels, err := NodeSelectorRequirementsAsSelector(invalidNsm)
	assert.Error(t, err)
	assert.Nil(t, labels)
}

func TestEmptyNodeSelectorRequirementsAsSelector(t *testing.T) {
	emptyNsm := []corev1.NodeSelectorRequirement{}
	res, err := NodeSelectorRequirementsAsSelector(emptyNsm)

	assert.NoError(t, err)
	assert.Equal(t, res, labels.Nothing())
}
