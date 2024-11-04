package domain

import (
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
	"slices"
)

// EffectiveBlueprint describes what the wanted state after evaluating the blueprint and the blueprintMask is.
// This is still a static description, so no actual state of the ecosystem is taken into consideration here.
type EffectiveBlueprint struct {
	// Dogus contains a set of exact dogu versions which should be present or absent in the CES instance after which this
	// blueprint was applied. Optional.
	Dogus []Dogu
	// Components contains a set of exact components versions which should be present or absent in the CES instance after which
	// this blueprint was applied. Optional.
	Components []Component
	// Config contains all config entries to set via blueprint.
	Config Config
}

// GetWantedDogus returns a list of all dogus which should be installed
func (effectiveBlueprint *EffectiveBlueprint) GetWantedDogus() []Dogu {
	var wantedDogus []Dogu
	for _, dogu := range effectiveBlueprint.Dogus {
		if dogu.TargetState == TargetStatePresent {
			wantedDogus = append(wantedDogus, dogu)
		}
	}
	return wantedDogus
}

// validateOnlyConfigForDogusInBlueprint checks that there is only config for dogus to install in the blueprint
func (effectiveBlueprint *EffectiveBlueprint) validateOnlyConfigForDogusInBlueprint() error {
	wantedDogus := util.Map(effectiveBlueprint.GetWantedDogus(), func(dogu Dogu) common.SimpleDoguName {
		return dogu.Name.SimpleName
	})
	var errorList []error
	for doguInConfig := range effectiveBlueprint.Config.Dogus {
		if !slices.Contains(wantedDogus, doguInConfig) {
			errorList = append(errorList, fmt.Errorf("setting config for dogu %q is not allowed as it will not be installed with the blueprint", doguInConfig))
		}
	}
	return errors.Join(errorList...)
}
