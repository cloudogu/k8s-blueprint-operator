package serializer

import (
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_toDomainTargetState(t *testing.T) {
	assert.Equal(t, domain.TargetState(domain.TargetStatePresent), ToDomainTargetState(false))
	assert.Equal(t, domain.TargetState(domain.TargetStateAbsent), ToDomainTargetState(true))
}

func Test_ToSerializerAbsentState(t *testing.T) {
	assert.False(t, ToSerializerAbsentState(domain.TargetStatePresent))
	assert.True(t, ToSerializerAbsentState(domain.TargetState(-1)))
	assert.True(t, ToSerializerAbsentState(domain.TargetStateAbsent))
}
