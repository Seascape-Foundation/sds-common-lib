package database

import (
	"github.com/Seascape-Foundation/sds-common-lib/data_type/key_value"
)

// Struct interface adds the database CRUD to the data struct.
//
// The interface that it accepts is the *remote.ClientSocket from the
// "github.com/Seascape-Foundation/sds-service-lib/remote" package.
type Crud interface {
	// Update the parameters by int flag. It calls UPDATE command
	Update(interface{}, uint8) error
	// Exist in the database or not. It calls EXIST command
	Exist(interface{}) bool

	// Insert into the database. It calls INSERT command
	Insert(interface{}) error
	// Load the database from database. It calls SELECT_ROW command
	Select(interface{}) error

	// It calls SELECT_ALL without WHERE clause of query.
	//
	// Result is then put to the second argument
	SelectAll(interface{}, interface{}) error

	// AllByCondition returns structs from database to the second argument.
	// The sql query should match to the condition.
	//
	// It calls SELECT_ALL with WHERE clause
	SelectAllByCondition(interface{}, key_value.KeyValue, interface{}) error // uses SELECT_ROW
}
