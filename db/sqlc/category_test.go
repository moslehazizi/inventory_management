package db

import (
	"context"
	"database/sql"
	"inventory_management/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomCategory(t *testing.T) Category {
	arg := CreateCategoryParams{
		CategoryName: util.RandomName(),
		SectionName:  util.RandomName(),
	}

	category, err := testQueries.CreateCategory(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, category)
	require.Equal(t, arg.CategoryName, category.CategoryName)
	require.Equal(t, arg.SectionName, category.SectionName)
	require.NotZero(t, category.ID)

	return category
}

func TestCreateCategory(t *testing.T) {
	createRandomCategory(t)
}

func TestGetCategory(t *testing.T) {
	category1 := createRandomCategory(t)
	category2, err := testQueries.GetCategory(context.Background(), category1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, category2)

	require.Equal(t, category1.ID, category2.ID)
	require.Equal(t, category1.CategoryName, category2.CategoryName)
	require.Equal(t, category1.SectionName, category2.SectionName)
}

func TestUpdateCategory(t *testing.T) {
	category1 := createRandomCategory(t)

	arg := UpdateCategoryParams{
		ID:           category1.ID,
		CategoryName: util.RandomName(),
		SectionName:  util.RandomName(),
	}

	category2, err := testQueries.UpdateCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, category2)

	require.NotEqual(t, category1.CategoryName, category2.CategoryName)
	require.NotEqual(t, category1.SectionName, category2.SectionName)
}

func TestDeleteCategory(t *testing.T) {
	category1 := createRandomCategory(t)

	err := testQueries.DeleteCategory(context.Background(), category1.ID)

	require.NoError(t, err)

	category2, err1 := testQueries.GetCategory(context.Background(), category1.ID)

	require.Error(t, err1)
	require.EqualError(t, err1, sql.ErrNoRows.Error())
	require.Empty(t, category2)
}

func TestListCategories(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomCategory(t)
	}
	arg := ListCategoriesParams{
		Limit:  5,
		Offset: 5,
	}
	categories, err := testQueries.ListCategories(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, categories, 5)

	for _, category := range categories {
		require.NotEmpty(t, category)
	}
}
