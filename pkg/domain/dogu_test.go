package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

func Test_TargetDogu_validate_errorOnMissingDoguName(t *testing.T) {
	dogus := []Dogu{
		{Version: version3212, TargetState: TargetStatePresent},
	}
	blueprint := Blueprint{Dogus: dogus}

	err := blueprint.Validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "dogu name must not be empty")
}

func Test_TargetDogu_validate_errorOnMissingVersionForPresentDogu(t *testing.T) {
	dogu := Dogu{Name: officialDogu1, TargetState: TargetStatePresent}

	err := dogu.validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "dogu version must not be empty")
}

func Test_TargetDogu_validate_missingVersionOkayForAbsentDogu(t *testing.T) {
	dogu := Dogu{Name: officialDogu1, TargetState: TargetStateAbsent}

	err := dogu.validate()

	require.Nil(t, err)
}

func Test_TargetDogu_validate_defaultToPresentState(t *testing.T) {
	dogu := Dogu{Name: officialDogu1, Version: version1_2_3}

	err := dogu.validate()

	require.Nil(t, err)
	assert.Equal(t, TargetState(TargetStatePresent), dogu.TargetState)
}

func Test_TargetDogu_validate_errorOnUnknownTargetState(t *testing.T) {
	dogu := Dogu{Name: officialDogu1, TargetState: -1}

	err := dogu.validate()

	require.Error(t, err)
	require.ErrorContains(t, err, "dogu target state is invalid: official/dogu1")
}

func Test_TargetDogu_validate_ProxySizeFormat(t *testing.T) {
	t.Run("error on invalid proxy body size format", func(t *testing.T) {
		// given
		parse := resource.MustParse("1Mi")
		dogu := Dogu{Name: officialDogu1, ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: &parse}}
		// when
		err := dogu.validate()
		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "dogu proxy body size is not in Decimal SI: official/dogu1")
	})

	t.Run("error on invalid volume size format", func(t *testing.T) {
		// given
		parse := resource.MustParse("1M")
		dogu := Dogu{Name: officialDogu1, MinVolumeSize: &parse}
		// when
		err := dogu.validate()
		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "dogu minimum volume size is not in Binary SI: official/dogu1")
	})

	t.Run("no error on empty quantity", func(t *testing.T) {
		// given
		dogu := Dogu{Name: officialDogu1, Version: version1_2_3}
		// when
		err := dogu.validate()
		// then
		require.NoError(t, err)
	})

	t.Run("no error on zero size quantity", func(t *testing.T) {
		// given
		zeroQuantity := resource.MustParse("0")
		dogu := Dogu{Name: officialDogu1, Version: version1_2_3, ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: &zeroQuantity}}
		// when
		err := dogu.validate()
		// then
		require.NoError(t, err)
		assert.Equal(t, resource.DecimalSI, dogu.ReverseProxyConfig.MaxBodySize.Format)
	})
}
