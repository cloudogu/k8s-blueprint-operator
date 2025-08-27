package doguregistry

import (
	"context"
	"fmt"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	gocache "github.com/patrickmn/go-cache"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// TODO: We should move this implementation in the dogu-descriptor-lib, so others can profit from it as well.

type Cache struct {
	repository remoteDoguDescriptorRepository
	cache      *gocache.Cache
}

func NewCache(repository remoteDoguDescriptorRepository, cache *gocache.Cache) *Cache {
	return &Cache{
		repository: repository,
		cache:      cache,
	}
}

func (c *Cache) GetLatest(ctx context.Context, name cescommons.QualifiedName) (*core.Dogu, error) {
	// we cannot cache latest for sure
	dogu, err := c.repository.GetLatest(ctx, name)
	if err != nil {
		return nil, err
	}
	version, err := core.ParseVersion(dogu.Version)
	if err != nil {
		return nil, fmt.Errorf("cannot populate cache as latest dogu version cannot be parsed from response: %w", err)
	}
	c.SetCached(
		cescommons.QualifiedVersion{
			Name:    name,
			Version: version,
		},
		dogu,
	)
	return dogu, nil
}

func (c *Cache) Get(ctx context.Context, version cescommons.QualifiedVersion) (*core.Dogu, error) {
	logger := log.FromContext(ctx).WithName("DoguDescriptorCache")
	dogu, found := c.GetCached(version)
	if found {
		logger.V(2).Info("dogu descriptor cache hit", "dogu", version)
		return dogu, nil
	}

	// cache missed
	dogu, err := c.repository.Get(ctx, version)
	if err != nil {
		return nil, err
	}
	c.SetCached(version, dogu)
	return dogu, nil
}

func cacheKeyFromReference(qualifiedDoguVersion cescommons.QualifiedVersion) string {
	return fmt.Sprintf("%s:%s", qualifiedDoguVersion.Name.String(), qualifiedDoguVersion.Version.String())
}

func (c *Cache) GetCached(qualifiedDoguVersion cescommons.QualifiedVersion) (*core.Dogu, bool) {
	cached, found := c.cache.Get(cacheKeyFromReference(qualifiedDoguVersion))
	if found {
		return cached.(*core.Dogu), found
	}
	return nil, found
}

func (c *Cache) SetCached(qualifiedDoguVersion cescommons.QualifiedVersion, dogu *core.Dogu) {
	c.cache.SetDefault(cacheKeyFromReference(qualifiedDoguVersion), dogu)
}
