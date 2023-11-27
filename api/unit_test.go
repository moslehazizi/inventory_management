package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "inventory_management/db/mock"
	db "inventory_management/db/sqlc"
	"inventory_management/util"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateUnit(t *testing.T) {
	unit := randomUnit()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"unit_name":  unit.UnitName,
				"unit_value": unit.UnitValue,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUnitParams{
					UnitName:  unit.UnitName,
					UnitValue: unit.UnitValue,
				}
				store.EXPECT().CreateUnit(gomock.Any(), gomock.Eq(arg)).Times(1).Return(unit, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUnitRequest(t, recorder.Body, unit)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"unit_name":  unit.UnitName,
				"unit_value": unit.UnitValue,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUnit(gomock.Any(), gomock.Any()).Times(1).Return(db.Unit{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidContext",
			body: gin.H{
				"unit_name":  "",
				"unit_value": "",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUnit(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/units"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})

	}

}

func TestListUnit(t *testing.T) {

	n := 5
	units := make([]db.Unit, n)

	for i := 0; i < n; i++ {
		units[i] = randomUnit()
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListUnitsParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.EXPECT().ListUnits(gomock.Any(), gomock.Eq(arg)).Times(1).Return(units, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUnits(t, recorder.Body, units)

			},
		},
		{
			name: "InternalError",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListUnits(gomock.Any(), gomock.Any()).Times(1).Return([]db.Unit{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListUnits(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				pageID:   1,
				pageSize: 10000,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListUnits(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/units"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func TestUpdateUnit(t *testing.T) {

	unit := randomUnit()
	unitUpdate := randomUnit()

	testCases := []struct {
		name          string
		UnitID        int64
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			UnitID: unit.ID,
			body: gin.H{
				"unit_name":  unitUpdate.UnitName,
				"unit_value": unitUpdate.UnitValue,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUnitParams{
					ID:        unit.ID,
					UnitName:  unitUpdate.UnitName,
					UnitValue: unitUpdate.UnitValue,
				}
				store.EXPECT().UpdateUnit(gomock.Any(), gomock.Eq(arg)).Times(1).Return(unitUpdate, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUnitRequest(t, recorder.Body, unitUpdate)
			},
		},
		{
			name:   "InternalError",
			UnitID: unit.ID,
			body: gin.H{
				"unit_name":  unitUpdate.UnitName,
				"unit_value": unitUpdate.UnitValue,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUnitParams{
					ID:        unit.ID,
					UnitName:  unitUpdate.UnitName,
					UnitValue: unitUpdate.UnitValue,
				}
				store.EXPECT().UpdateUnit(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Unit{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "InvalidContext",
			UnitID: unit.ID,
			body: gin.H{
				"category_name": "",
				"section_name":  "",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateCategory(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:       "InvalidID",
			UnitID: 0,
			body: gin.H{
				"category_name": "",
				"section_name":  "",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateCategory(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:       "NotFound",
			UnitID: unit.ID,
			body: gin.H{
				"unit_name": unitUpdate.UnitName,
				"unit_value":  unitUpdate.UnitValue,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUnitParams{
					ID:          unit.ID,
					UnitName:    unitUpdate.UnitName,
					UnitValue: unitUpdate.UnitValue,
				}
				store.EXPECT().UpdateUnit(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Unit{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/units/%d", tc.UnitID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})

	}

}

func TestDeleteUnit(t *testing.T) {
	unit := randomUnit()

	testCases := []struct {
		name          string
		unitID    int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(T *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			unitID: unit.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteUnit(gomock.Any(), gomock.Eq(unit.ID)).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:       "NotFound",
			unitID: unit.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteUnit(gomock.Any(), gomock.Eq(unit.ID)).Times(1).Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "InternalError",
			unitID: unit.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteUnit(gomock.Any(), gomock.Eq(unit.ID)).Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "InvalidID",
			unitID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteUnit(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/units/%d", tc.unitID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)

			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}

}

func randomUnit() db.Unit {
	return db.Unit{
		ID:        util.RandomInt(1, 1000),
		UnitName:  util.RandomName(),
		UnitValue: util.RandomInt(1, 8),
	}
}

// func requireBodyMatchUnit(t *testing.T, body *bytes.Buffer, unit db.Unit) {
// 	data, err := ioutil.ReadAll(body)
// 	require.NoError(t, err)

// 	var gotUnit db.Unit
// 	err = json.Unmarshal(data, &gotUnit)
// 	require.NoError(t, err)
// 	require.Equal(t, unit, gotUnit)
// }

func requireBodyMatchUnits(t *testing.T, body *bytes.Buffer, units []db.Unit) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUnits []db.Unit
	err = json.Unmarshal(data, &gotUnits)
	require.NoError(t, err)
	require.Equal(t, units, gotUnits)
}

func requireBodyMatchUnitRequest(t *testing.T, body *bytes.Buffer, unit db.Unit) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUnit db.Unit
	err = json.Unmarshal(data, &gotUnit)
	require.NoError(t, err)
	require.Equal(t, unit.UnitName, gotUnit.UnitName)
	require.Equal(t, unit.UnitValue, gotUnit.UnitValue)
}
