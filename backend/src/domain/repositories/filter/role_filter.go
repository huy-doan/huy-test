package filter

import "fmt"

// RoleFilter represents filtering options for roles
type RoleFilter struct {
	BaseFilter
	Name *string
}

// NewRoleFilter creates a new RoleFilter with valid sort fields
func NewRoleFilter() *RoleFilter {
	filter := &RoleFilter{}
	filter.ValidSortFields = map[string]bool{
		"name": true,
	}
	return filter
}

// ApplyFilters applies all filter conditions based on the filter fields
func (f *RoleFilter) ApplyFilters() {
	if f.Name != nil {
		f.AddCondition("name", Like, fmt.Sprintf("%%%s%%", *f.Name))
	}
}
