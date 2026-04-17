package handlers

import (
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/models"
	"gorm.io/gorm"
)

func collectDepartmentSubtreeIDs(db *gorm.DB, rootID uint) ([]uint, error) {
	if rootID == 0 {
		return []uint{}, nil
	}

	ids := []uint{rootID}
	frontier := []uint{rootID}
	seen := map[uint]struct{}{
		rootID: {},
	}

	for len(frontier) > 0 {
		var children []uint
		if err := db.Model(&models.Department{}).Where("parent_id IN ?", frontier).Pluck("id", &children).Error; err != nil {
			return nil, err
		}
		next := make([]uint, 0, len(children))
		for _, childID := range children {
			if _, ok := seen[childID]; ok {
				continue
			}
			seen[childID] = struct{}{}
			ids = append(ids, childID)
			next = append(next, childID)
		}
		frontier = next
	}

	return ids, nil
}
