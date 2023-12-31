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
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
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
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "InternalError",
			categoryID: category.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCategory(gomock.Any(), gomock.Eq(category.ID)).Times(1).Return(db.Category{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "InvalidID",
			categoryID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCategory(gomock.Any(), gomock.Any()).Times(0)
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
			url := fmt.Sprintf("/categories/%d", tc.categoryID)
			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}

}

func TestCreateCategory(t *testing.T) {
	category := randomCategory()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"category_name": category.CategoryName,
				"section_name":  category.SectionName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateCategoryParams{
					CategoryName: category.CategoryName,
					SectionName:  category.SectionName,
				}
				store.EXPECT().CreateCategory(gomock.Any(), gomock.Eq(arg)).Times(1).Return(category, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCategoryRequest(t, recorder.Body, category)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"category_name": category.CategoryName,
				"section_name":  category.SectionName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateCategoryParams{
					CategoryName: category.CategoryName,
					SectionName:  category.SectionName,
				}
				store.EXPECT().CreateCategory(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Category{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidContext",
			body: gin.H{
				"category_name": "",
				"section_name":  "",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateCategory(gomock.Any(), gomock.Any()).Times(0)
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

			url := "/categories"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})

	}

}

func TestListCategory(t *testing.T) {

	n := 6
	c_categories := make([]db.Category, n)

	for i := 0; i < n; i++ {
		c_categories[i] = randomCategory()
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
				arg := db.ListCategoriesParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.EXPECT().ListCategories(gomock.Any(), gomock.Eq(arg)).Times(1).Return(c_categories, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCategories(t, recorder.Body, c_categories)

			},
		},
		{
			name: "InternalError",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListCategoriesParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.EXPECT().ListCategories(gomock.Any(), gomock.Eq(arg)).Times(1).Return([]db.Category{}, sql.ErrConnDone)
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
				store.EXPECT().ListCategories(gomock.Any(), gomock.Any()).Times(0)
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
				store.EXPECT().ListCategories(gomock.Any(), gomock.Any()).Times(0)
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

			url := "/categories"
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

func TestUpdateCategory(t *testing.T) {

	category := randomCategory()
	categoryUpdate := randomCategory()

	testCases := []struct {
		name          string
		CategoryID    int64
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			CategoryID: category.ID,
			body: gin.H{
				"category_name": categoryUpdate.CategoryName,
				"section_name":  categoryUpdate.SectionName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateCategoryParams{
					ID:           category.ID,
					CategoryName: categoryUpdate.CategoryName,
					SectionName:  categoryUpdate.SectionName,
				}
				store.EXPECT().UpdateCategory(gomock.Any(), gomock.Eq(arg)).Times(1).Return(categoryUpdate, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCategoryRequest(t, recorder.Body, categoryUpdate)
			},
		},
		{
			name: "InternalError",
			CategoryID: category.ID,
			body: gin.H{
				"category_name": categoryUpdate.CategoryName,
				"section_name":  categoryUpdate.SectionName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateCategoryParams{
					ID:           category.ID,
					CategoryName: categoryUpdate.CategoryName,
					SectionName:  categoryUpdate.SectionName,
				}
				store.EXPECT().UpdateCategory(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Category{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidContext",
			CategoryID: category.ID,
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
			CategoryID: 0,
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
			CategoryID: category.ID,
			body: gin.H{
				"category_name": categoryUpdate.CategoryName,
				"section_name":  categoryUpdate.SectionName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateCategoryParams{
					ID:           category.ID,
					CategoryName: categoryUpdate.CategoryName,
					SectionName:  categoryUpdate.SectionName,
				}
				store.EXPECT().UpdateCategory(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Category{}, sql.ErrNoRows)
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

			url := fmt.Sprintf("/categories/%d", tc.CategoryID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})

	}

}

func TestDeleteCategory(t *testing.T) {
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
				store.EXPECT().DeleteCategory(gomock.Any(), gomock.Eq(category.ID)).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:       "NotFound",
			categoryID: category.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteCategory(gomock.Any(), gomock.Eq(category.ID)).Times(1).Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "InternalError",
			categoryID: category.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteCategory(gomock.Any(), gomock.Eq(category.ID)).Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "InvalidID",
			categoryID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteCategory(gomock.Any(), gomock.Any()).Times(0)
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
			url := fmt.Sprintf("/categories/%d", tc.categoryID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)

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

func requireBodyMatchCategories(t *testing.T, body *bytes.Buffer, c_categories []db.Category) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotCategories []db.Category
	err = json.Unmarshal(data, &gotCategories)
	require.NoError(t, err)
	require.Equal(t, c_categories, gotCategories)
}

func requireBodyMatchCategoryRequest(t *testing.T, body *bytes.Buffer, category db.Category) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotCategory db.Category
	err = json.Unmarshal(data, &gotCategory)
	require.NoError(t, err)
	require.Equal(t, category.CategoryName, gotCategory.CategoryName)
	require.Equal(t, category.SectionName, gotCategory.SectionName)
}
