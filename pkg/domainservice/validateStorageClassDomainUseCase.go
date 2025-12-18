package domainservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type ValidateStorageClassDomainUseCase struct {
	doguRepository DoguInstallationRepository
}

func NewValidateStorageClassDomainUseCase(doguRepository DoguInstallationRepository) *ValidateStorageClassDomainUseCase {
	return &ValidateStorageClassDomainUseCase{
		doguRepository,
	}
}

func (useCase *ValidateStorageClassDomainUseCase) ValidateDoguStorageClass(ctx context.Context, effectiveBlueprint domain.EffectiveBlueprint) error {
	logger := log.FromContext(ctx).WithName("ValidateStorageClassDomainUseCase.ValidateDoguStorageClass")

	logger.V(2).Info("load installed dogus")
	installedDogus, err := useCase.doguRepository.GetAll(ctx)
	if err != nil {
		return &InternalError{WrappedError: err, Message: "cannot get installed dogus for storage class validation"}
	}
	logger.V(2).Info("installed dogus loaded", "dogus", installedDogus)

	var errs []error
	for _, wantedDogu := range effectiveBlueprint.GetWantedDogus() {
		installedDogu, exists := installedDogus[wantedDogu.Name.SimpleName]
		if exists && !equalPtr(installedDogu.StorageClassName, wantedDogu.StorageClassName) {
			errs = append(errs,
				fmt.Errorf("wanted dogu %s's storage class differs from installed dogu: %s != %s",
					wantedDogu.Name.SimpleName,
					valueOrNil(wantedDogu.StorageClassName),
					valueOrNil(installedDogu.StorageClassName),
				),
			)
		}
	}

	err = errors.Join(errs...)
	if err != nil {
		return &domain.InvalidBlueprintError{WrappedError: err,
			Message: "storage classes are invalid in effective blueprint"}
	}

	return nil
}

func equalPtr[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b // both nil → true; one nil → false
	}
	return *a == *b
}

func valueOrNil(p *string) string {
	if p == nil {
		return "<nil>"
	}
	return *p
}
