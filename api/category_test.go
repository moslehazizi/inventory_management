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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetCategory(t *testing.T) {
	category := randomCategory()

	testCases := []struct {
		name          string
		categoryID    int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(T *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			categoryID: category.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCategory(gomock.Any(), gomock.Eq(category.ID)).Times(1).Return(category, nil)
			},
			checkResponse: func(T *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCategory(t, recorder.Body, category)

			},
		},
		{
			name:       "NotFound",
			categoryID: category.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCategory(gomock.Any(), gomock.Eq(category.ID)).Times(1).Return(db.Category{}, sql.ErrNoRows)
			},
			checkResponse: func(T *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "InternalError",
			categoryID: category.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCategory(gomock.Any(), gomock.Eq(category.ID)).Times(1).Return(db.Category{}, sql.ErrConnDone)
			},
			checkResponse: func(T *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "InvalidID",
			categoryID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCategory(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(T *testing.T, recorder *httptest.ResponseRecorder) {
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
			url := fmt.Sprintf("/categories/%d", tc.categoryID)
			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}

}

func randomCategory() db.Category {
	return db.Category{
		ID:           util.RandomInt(1, 1000),
		CategoryName: util.RandomName(),
		SectionName:  util.RandomName(),
	}
}

func requireBodyMatchCategory(t *testing.T, body *bytes.Buffer, category db.Category) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotCategory db.Category
	err = json.Unmarshal(data, &gotCategory)
	require.NoError(t, err)
	require.Equal(t, category, gotCategory)
}
