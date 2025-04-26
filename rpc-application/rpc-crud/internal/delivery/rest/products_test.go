package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/berezovskiydeval/rpc-crud/internal/delivery/rest"
	"github.com/berezovskiydeval/rpc-crud/internal/domain"
	"github.com/berezovskiydeval/rpc-crud/mocks"
	"github.com/gorilla/mux"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"bytes"
)

func TestGetAllProductsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockProductService(ctrl)

	handler := rest.NewHandler(mockService)

	expected := []domain.Product{
		{Id: 1, Name: "Phone", Price: 999},
		{Id: 2, Name: "Laptop", Price: 1499},
	}

	mockService.EXPECT().
		GetAll(gomock.Any()).
		Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var got []domain.Product
	err := json.NewDecoder(res.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestCreateProductsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockProductService(ctrl)
	handler := rest.NewHandler(mockService)

	request := domain.Product{Name: "test1_name", Price: 199}
	expectedID := int64(1)

	mockService.EXPECT().Create(gomock.Any(), request).Return(expectedID, nil)

	body, err := json.Marshal(request)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var response map[string]int64
	err = json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, response["id"])
}

func TestUpdateProductsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockProductService(ctrl)
	handler := rest.NewHandler(mockService)

	request := domain.ProductUpdate{Name: "test2_name", Price: 19}
	expected := int64(1)

	mockService.EXPECT().Update(gomock.Any(), expected, request).Return(expected, nil)

	body, err := json.Marshal(request)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	
	w := httptest.NewRecorder()

	handler.Update(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response map[string]int64
	err = json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expected, response["id"])
}

func TestDeleteProductsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockProductService(ctrl)
	handler := rest.NewHandler(mockService)

	request := int64(1)
	expectedID := int64(1)

	mockService.EXPECT().Delete(gomock.Any(), request).Return(expectedID, nil)

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)

	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	w := httptest.NewRecorder()
	handler.Delete(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response map[string]int64
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, response["id"])
}
