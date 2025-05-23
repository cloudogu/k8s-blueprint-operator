package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

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
	dogu := Dogu{Name: officialDogu1, Version: version123}

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
		require.ErrorContains(t, err, "dogu proxy body size is not in Decimal SI (\"M\" or \"G\"): official/dogu1")
	})

	t.Run("error on invalid volume size format", func(t *testing.T) {
		// given
		parse := resource.MustParse("1M")
		dogu := Dogu{Name: officialDogu1, MinVolumeSize: &parse}
		// when
		err := dogu.validate()
		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "dogu minimum volume size is not in Binary SI (\"Mi\" or \"Gi\"): official/dogu1")
	})

	t.Run("no error on empty quantity", func(t *testing.T) {
		// given
		dogu := Dogu{Name: officialDogu1, Version: version123}
		// when
		err := dogu.validate()
		// then
		require.NoError(t, err)
	})

	t.Run("no error on zero size quantity", func(t *testing.T) {
		// given
		zeroQuantity := resource.MustParse("0")
		dogu := Dogu{Name: officialDogu1, Version: version123, ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: &zeroQuantity}}
		// when
		err := dogu.validate()
		// then
		require.NoError(t, err)
		assert.Equal(t, resource.DecimalSI, dogu.ReverseProxyConfig.MaxBodySize.Format)
	})
}

func Test_TargetDogu_validate_AdditionalMounts(t *testing.T) {
	t.Run("additionalMounts ok", func(t *testing.T) {
		// given
		dogu := Dogu{Name: nginxStatic, Version: version123, AdditionalMounts: []ecosystem.AdditionalMount{
			{
				SourceType: ecosystem.DataSourceConfigMap,
				Name:       "html-config",
				Volume:     "customhtml",
				Subfolder:  "test",
			},
		}}
		// when
		err := dogu.validate()
		// then
		require.NoError(t, err)
	})

	t.Run("unknown sourceType", func(t *testing.T) {
		// given
		dogu := Dogu{Name: nginxStatic, Version: version123, AdditionalMounts: []ecosystem.AdditionalMount{
			{
				SourceType: "unsupportedType",
				Name:       "html-config",
				Volume:     "customhtml",
				Subfolder:  "test",
			},
		}}
		// when
		err := dogu.validate()
		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "dogu is invalid: dogu additional mounts sourceType must be one of 'ConfigMap', 'Secret': k8s/nginx-static")
	})

	t.Run("subfolder is no relative path", func(t *testing.T) {
		// given
		dogu := Dogu{Name: nginxStatic, Version: version123, AdditionalMounts: []ecosystem.AdditionalMount{
			{
				SourceType: ecosystem.DataSourceConfigMap,
				Name:       "html-config",
				Volume:     "customhtml",
				Subfolder:  "/test",
			},
		}}
		// when
		err := dogu.validate()
		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "dogu is invalid: dogu additional mounts Subfolder must be a relative path : k8s/nginx-static")
	})

}
