package convert

import (
	modelScreen "github.com/huydq/test/internal/domain/model/screen"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	"github.com/huydq/test/internal/infrastructure/persistence/screen/dto"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
	"gorm.io/gorm"
)

// ToScreenDTO converts a Screen domain model to a Screen
func ToScreenDTO(screen *modelScreen.Screen) *dto.Screen {
	if screen == nil {
		return nil
	}

	result := &dto.Screen{
		ID:         screen.ID,
		Name:       screen.Name,
		ScreenCode: screen.ScreenCode,
		ScreenPath: screen.ScreenPath,
		BaseColumnTimestamp: persistence.BaseColumnTimestamp{
			CreatedAt: screen.CreatedAt,
			UpdatedAt: screen.UpdatedAt,
		},
	}

	// Handle the conversion from *time.Time to gorm.DeletedAt
	if screen.DeletedAt != nil {
		result.DeletedAt = gorm.DeletedAt{
			Time:  *screen.DeletedAt,
			Valid: true,
		}
	} else {
		result.DeletedAt = gorm.DeletedAt{
			Valid: false,
		}
	}

	return result
}

// ToScreenModel converts a Screen to a Screen domain model
func ToScreenModel(dtoObj *dto.Screen) *modelScreen.Screen {
	if dtoObj == nil {
		return nil
	}

	result := &modelScreen.Screen{
		ID:         dtoObj.ID,
		Name:       dtoObj.Name,
		ScreenCode: dtoObj.ScreenCode,
		ScreenPath: dtoObj.ScreenPath,
		BaseColumnTimestamp: util.BaseColumnTimestamp{
			CreatedAt: dtoObj.CreatedAt,
			UpdatedAt: dtoObj.UpdatedAt,
		},
	}

	// Handle the conversion from gorm.DeletedAt to *time.Time
	if dtoObj.DeletedAt.Valid {
		deletedAt := dtoObj.DeletedAt.Time
		result.DeletedAt = &deletedAt
	} else {
		result.DeletedAt = nil
	}

	return result
}

// ToScreenDTOs converts a list of Screen domain models to a list of ScreenDTOs
func ToScreenDTOs(screens []*modelScreen.Screen) []*dto.Screen {
	if screens == nil {
		return nil
	}

	result := make([]*dto.Screen, len(screens))
	for i, screen := range screens {
		result[i] = ToScreenDTO(screen)
	}
	return result
}

// ToScreenModels converts a list of ScreenDTOs to a list of Screen domain models
func ToScreenModels(dtos []*dto.Screen) []*modelScreen.Screen {
	if dtos == nil {
		return nil
	}

	result := make([]*modelScreen.Screen, len(dtos))
	for i, dto := range dtos {
		result[i] = ToScreenModel(dto)
	}
	return result
}
