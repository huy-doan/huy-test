package fixtures

import "github.com/huydq/test/src/domain/models"

// GetMockRole returns a mock role with the specified ID and code
func GetMockRole(id int, code string) *models.Role {
	return &models.Role{
		ID:   id,
		Code: code,
	}
}
