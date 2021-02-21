package db

import "testing"

func TestRemoveUserRole(t *testing.T) {
	err := RemoveUserRole(5, []int64{1, 2})
	if err != nil {
		panic(err)
	}
}
