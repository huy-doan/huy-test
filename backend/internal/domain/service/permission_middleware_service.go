package service

import (
	"context"
	"sync"
	"time"

	"github.com/huydq/test/internal/pkg/logger"

	object "github.com/huydq/test/internal/domain/object/permission"
	repositoryPermission "github.com/huydq/test/internal/domain/repository/permission"
)

type PermissionMiddlewareService interface {
	HasPermission(ctx context.Context, roleID int, permissions ...object.PermissionCode) (bool, error)

	GetUserPermissions(ctx context.Context, roleID int) ([]string, error)
}

type permissionCache struct {
	roleID      int
	permissions []string
	expiry      time.Time
}

const permissionCacheExpiration = 1 * time.Minute // Cache expiration 1 minute
const permissionCacheCleanupInterval = 5 * time.Minute // Cache cleanup interval 5 minutes
type permissionMiddlewareServiceImpl struct {
	repository repositoryPermission.PermissionRepository
	logger     logger.Logger
	cache      map[int]*permissionCache
	cacheMutex sync.RWMutex
}

func NewPermissionMiddlewareService(repository repositoryPermission.PermissionRepository, logger logger.Logger) PermissionMiddlewareService {
	service := &permissionMiddlewareServiceImpl{
		repository: repository,
		logger:     logger,
		cache:      make(map[int]*permissionCache),
	}

	go service.startCacheCleaner()

	return service
}

func (s *permissionMiddlewareServiceImpl) HasPermission(ctx context.Context, roleID int, permissions ...object.PermissionCode) (bool, error) {
	permissionStrings := make([]string, len(permissions))
	for i, p := range permissions {
		permissionStrings[i] = string(p)
	}

	userPermissions, err := s.GetUserPermissions(ctx, roleID)
	if err != nil {
		s.logger.Error("Error getting user permissions", map[string]any{
			"roleID": roleID,
			"error":  err.Error(),
		})
		return false, err
	}

	for _, userPerm := range userPermissions {
		for _, requiredPerm := range permissionStrings {
			if userPerm == requiredPerm {
				return true, nil
			}
		}
	}

	return false, nil
}

func (s *permissionMiddlewareServiceImpl) GetUserPermissions(ctx context.Context, roleID int) ([]string, error) {
	s.cacheMutex.RLock()
	cacheItem, found := s.cache[roleID]
	s.cacheMutex.RUnlock()

	now := time.Now()

	if found && now.Before(cacheItem.expiry) {
		return cacheItem.permissions, nil
	}

	permissions, err := s.repository.GetPermissionCodesByRoleID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	s.cacheMutex.Lock()
	s.cache[roleID] = &permissionCache{
		roleID:      roleID,
		permissions: permissions,
		expiry:      now.Add(permissionCacheExpiration),
	}
	s.cacheMutex.Unlock()

	return permissions, nil
}

func (s *permissionMiddlewareServiceImpl) startCacheCleaner() {
	ticker := time.NewTicker(permissionCacheCleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		s.cacheMutex.Lock()
		now := time.Now()
		for roleID, item := range s.cache {
			if now.After(item.expiry) {
				delete(s.cache, roleID)
			}
		}
		s.cacheMutex.Unlock()
	}
}
