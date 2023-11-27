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

func TestGetGood(t *testing.T) {
	good := randomGood()

	testCases := []struct {
		name          string
		goodID        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(T *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			goodID: good.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetGood(gomock.Any(), gomock.Eq(good.ID)).Times(1).Return(good, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchGood(t, recorder.Body, good)
			},
		},
		{
			name:   "NotFound",
			goodID: good.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetGood(gomock.Any(), gomock.Eq(good.ID)).Times(1).Return(db.Good{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			goodID: good.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetGood(gomock.Any(), gomock.Eq(good.ID)).Times(1).Return(db.Good{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			goodID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetGood(gomock.Any(), gomock.Any()).Times(0)
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
			url := fmt.Sprintf("/goods/%d", tc.goodID)
			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}

}

func TestListGoods(t *testing.T) {

	n := 5
	goods := make([]db.Good, n)

	for i := 0; i < n; i++ {
		goods[i] = randomGood()
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
		// {
		// 	name: "OK",
		// 	query: Query{
		// 		pageID:   1,
		// 		pageSize: n,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		arg := db.ListGoodsParams{
		// 			Category: goods[1].Category,
		// 			Limit:    int32(n),
		// 			Offset:   0,
		// 		}
		// 		store.EXPECT().ListGoods(gomock.Any(), gomock.Eq(arg)).Times(1).Return(goods, nil)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusOK, recorder.Code)
		// 		requireBodyMatchGoods(t, recorder.Body, goods)
		// 	},
		// },
		// {
		// 	name: "InternalError",
		// 	query: Query{
		// 		pageID:   1,
		// 		pageSize: n,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		arg := db.ListGoodsParams{
		// 			Category: goods[1].Category,
		// 			Limit:    int32(n),
		// 			Offset:   0,
		// 		}
		// 		store.EXPECT().ListGoods(gomock.Any(), gomock.Eq(arg)).Times(1).Return([]db.Good{}, sql.ErrConnDone)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		// 	},
		// },
		{
			name: "InvalidPageID",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListGoods(gomock.Any(), gomock.Any()).Times(0)
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
				store.EXPECT().ListGoods(gomock.Any(), gomock.Any()).Times(0)
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

			url := "/goods"
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

func TestCreateGood(t *testing.T) {
	good := randomGood()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"category":  good.Category,
				"model":     good.Model,
				"unit":      good.Unit,
				"amount":    good.Amount,
				"good_desc": good.GoodDesc,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateGoodParams{
					Category: int64(good.Category),
					Model:    good.Model,
					Unit:     int64(good.Unit),
					Amount:   int64(good.Amount),
					GoodDesc: good.GoodDesc,
				}
				store.EXPECT().CreateGood(gomock.Any(), gomock.Eq(arg)).Times(1).Return(good, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchGoodRequest(t, recorder.Body, good)
			},
		},
		// {
		// 	name: "InternalError",
		// 	body: gin.H{
		// 		"category_name": category.CategoryName,
		// 		"section_name":  category.SectionName,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		arg := db.CreateCategoryParams{
		// 			CategoryName: category.CategoryName,
		// 			SectionName:  category.SectionName,
		// 		}
		// 		store.EXPECT().CreateCategory(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Category{}, sql.ErrConnDone)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "InvalidContext",
		// 	body: gin.H{
		// 		"category_name": "",
		// 		"section_name":  "",
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().CreateCategory(gomock.Any(), gomock.Any()).Times(0)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
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

			url := "/goods"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})

	}

}

func randomGood() db.Good {
	g_category := randomCategory()
	g_unit := randomUnit()
	return db.Good{
		ID:       util.RandomInt(1, 1000),
		Category: g_category.ID,
		Model:    util.RandomName(),
		Unit:     g_unit.ID,
		Amount:   util.RandomInt(5, 9),
		GoodDesc: util.RandomName(),
	}
}

func requireBodyMatchGood(t *testing.T, body *bytes.Buffer, good db.Good) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotGood db.Good
	err = json.Unmarshal(data, &gotGood)
	require.NoError(t, err)
	require.Equal(t, good, gotGood)
}

func requireBodyMatchGoods(t *testing.T, body *bytes.Buffer, goods []db.Good) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotGoods []db.Good
	err = json.Unmarshal(data, &gotGoods)
	require.NoError(t, err)
	require.Equal(t, goods, gotGoods)
}

func requireBodyMatchGoodRequest(t *testing.T, body *bytes.Buffer, good db.Good) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotGood db.Good
	err = json.Unmarshal(data, &gotGood)
	require.NoError(t, err)
	require.Equal(t, good.Category, gotGood.Category)
	require.Equal(t, good.Model, gotGood.Model)
	require.Equal(t, good.Unit, gotGood.Unit)
	require.Equal(t, good.Amount, gotGood.Amount)
	require.Equal(t, good.GoodDesc, gotGood.GoodDesc)
}
