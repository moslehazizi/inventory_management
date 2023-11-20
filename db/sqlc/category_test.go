package db

import (
	"context"
	"testing"
)

func TestCreateCategory(t *testing.T) {
	arg := 	CreateCategory {
		CategoryName: "hammer",
	}

	category, err := testQueries.CreateCategory(context.Background(), arg)

}