package db

import "testing"

func TestSaveUserRole(t *testing.T) {
	err := SaveUserRole(5, []int64{1, 2})
	if err != nil {
		panic(err)
	}
}
