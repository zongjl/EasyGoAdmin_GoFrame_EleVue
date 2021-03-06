// ============================================================================
// This is auto-generated by gf cli tool only once. Fill this file as you wish.
// ============================================================================

package dao

import (
	"easygoadmin/app/dao/internal"
)

// exampleDao is the manager for logic model data accessing
// and custom defined data operations functions management. You can define
// methods on it to extend its functionality as you wish.
type exampleDao struct {
	internal.ExampleDao
}

var (
	// Example is globally public accessible object for table sys_example operations.
	Example = exampleDao{
		internal.Example,
	}
)

// Fill with you ideas below.