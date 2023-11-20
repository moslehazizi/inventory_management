package db

import (
	"context"
	"inventory_management/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomUnit(t *testing.T) Unit {
	arg := CreateUnitParams {
		UnitName: util.RandomName(),
		UnitValue: util.RandomInt(1, 3),
	}
	unit, err := testQueries.CreateUnit(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, unit)

	require.Equal(t, arg.UnitName, unit.UnitName)
	require.Equal(t, arg.UnitValue, unit.UnitValue)
	require.NotZero(t, unit.ID)

	return unit
}

func TestCreateUnit(t *testing.T) {
	createRandomUnit(t)
}

func TestListUnits(t *testing.T) {

	for i:=0;i<10;i++{
		createRandomUnit(t)
	}
	arg := ListUnitsParams {
		Limit: 5,
		Offset: 5,
	}
	units, err := testQueries.ListUnits(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, units, 5)

	for _, unit :=range units {
		require.NotEmpty(t, unit)
	}
}

func TestUpdateUnit(t *testing.T) {
	unit1 := createRandomUnit(t)

	arg := UpdateUnitParams {
		ID: unit1.ID,
		UnitName: util.RandomName(),
		UnitValue: util.RandomInt(4, 6),
	}
	unit2, err := testQueries.UpdateUnit(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, unit2)
	require.NotEqual(t, unit1.UnitName, unit2.UnitName)
	require.NotEqual(t, unit1.UnitValue, unit2.UnitValue)
}

func TestDeleteUnit(t *testing.T) {
	unit1 := createRandomUnit(t) 

	err2 := testQueries.DeleteUnit(context.Background(), unit1.ID)

	require.NoError(t, err2)

	// Must be edited 
}