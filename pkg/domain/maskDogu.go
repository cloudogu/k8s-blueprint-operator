package domain

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
)

// MaskDogu defines a Dogu, its version, and the installation state in which it is supposed to be after a blueprint
// was applied for a blueprintMask.
type MaskDogu struct {
	// Name is the qualified name of the dogu.
	Name cescommons.QualifiedName
	// Version defines the version of the dogu that is to be installed. This version is optional and overrides
	// the version of the dogu from the blueprint.
	Version core.Version
	// Absent defines if the dogu should be absent in the ecosystem. Defaults to false.
	Absent bool
}

func (dogu MaskDogu) validate() error {
	return dogu.Name.Validate()
}
