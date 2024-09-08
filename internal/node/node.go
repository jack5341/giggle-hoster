package node

import (
	"errors"

	"github.com/jack5341/giggle-hoster/internal/types"
	"gorm.io/gorm"
)

var (
	ErrFitNodeCouldNotBeFound = errors.New("fit not could not be found")
)

// func CreateNode(db *gorm.DB) {

// }

// func DeleteNode(db *gorm.DB) {

// }

func FindFitNode(db *gorm.DB, requestedMem int, requestedCPU int) (types.Node, error) {
	var node types.Node
	tx := db.Begin()

	err := tx.Where("free_mem >= ? AND free_cpu >= ?", requestedMem, requestedCPU).
		Order("free_mem DESC, free_cpu DESC").
		First(&node).Error

	if err != nil {
		return node, errors.Join(ErrFitNodeCouldNotBeFound, err)
	}

	return node, nil
}
