package db

import (
	"context"
	"database/sql"
	"inventory_management/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomGood(t *testing.T, category Category, unit Unit) Good {

	arg := CreateGoodParams{
		Category: category.ID,
		Model:    util.RandomName(),
		Unit:     unit.ID,
		Amount:   util.RandomInt(3, 7),
		GoodDesc: util.RandomName(),
	}

	good, err := testQueries.CreateGood(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, good)
	require.Equal(t, arg.Category, good.Category)
	require.Equal(t, arg.Model, good.Model)
	require.Equal(t, arg.Unit, good.Unit)
	require.Equal(t, arg.Amount, good.Amount)
	require.Equal(t, arg.GoodDesc, good.GoodDesc)

	require.NotZero(t, good.ID)

	return good
}

func TestCreateGood(t *testing.T) {
	category := createRandomCategory(t)
	unit := createRandomUnit(t)
	createRandomGood(t, category, unit)
}

func TestGetGood(t *testing.T) {
	category := createRandomCategory(t)
	unit := createRandomUnit(t)
	good1 := createRandomGood(t, category, unit)
	good2, err := testQueries.GetGood(context.Background(), good1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, good2)

	require.Equal(t, good1.Category, good2.Category)
	require.Equal(t, good1.Model, good2.Model)
	require.Equal(t, good1.Unit, good2.Unit)
	require.Equal(t, good1.Amount, good2.Amount)
	require.Equal(t, good1.GoodDesc, good2.GoodDesc)
	require.WithinDuration(t, good1.CreatedAt, good2.CreatedAt, time.Second)
}

func TestListGoods(t *testing.T) {

	category := createRandomCategory(t)
	unit := createRandomUnit(t)

	for i := 0; i < 10; i++ {
		createRandomGood(t, category, unit)
	}

	arg := ListGoodsParams{
		Category: category.ID,
		Limit:    5,
		Offset:   5,
	}
	goods, err := testQueries.ListGoods(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, goods, 5)

	for _, good := range goods {
		require.NotEmpty(t, good)
		require.True(t, good.Category == category.ID || good.Unit == unit.ID)
	}
}

func TestUpdateGood(t *testing.T) {
	category := createRandomCategory(t)
	unit1 := createRandomUnit(t)
	unit2 := createRandomUnit(t)
	good1 := createRandomGood(t, category, unit1)

	arg := UpdateGoodParams {
		ID: good1.ID,
		Unit: unit2.ID,
		Amount: util.RandomInt(8, 11),
	}

	good2, err := testQueries.UpdateGood(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, good2)
	require.Equal(t, good1.ID, good2.ID)
	require.NotEqual(t, good1.Unit, good2.Unit)
	require.NotEqual(t, good1.Amount, good2.Amount)
}

func TestDeleteGood(t *testing.T) {
	category := createRandomCategory(t)
	unit := createRandomUnit(t)
	good1 := createRandomGood(t, category, unit)

	err1:= testQueries.DeleteGood(context.Background(), good1.ID)

	require.NoError(t, err1)

	good2, err2 := testQueries.GetGood(context.Background(), good1.ID)

	require.Error(t, err2)
	require.EqualError(t, err2, sql.ErrNoRows.Error())

	require.Empty(t, good2)
}