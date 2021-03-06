// ============================================================================
// This is auto-generated by gf cli tool only once. Fill this file as you wish.
// ============================================================================

package dao

import (
	"easygoadmin/app/dao/internal"
)

// PositionDao is the manager for logic model data accessing
// and custom defined data operations functions management. You can define
// methods on it to extend its functionality as you wish.
type PositionDao struct {
	internal.PositionDao
}

var (
	// Position is globally public accessible object for table sys_position operations.
	Position = PositionDao{
		internal.Position,
	}
)

// Fill with you ideas below.