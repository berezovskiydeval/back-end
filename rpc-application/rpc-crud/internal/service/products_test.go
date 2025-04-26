package service_test

import (
	"context"
	"testing"

	"github.com/berezovskiydeval/rpc-crud/internal/domain"
	"github.com/berezovskiydeval/rpc-crud/internal/service"
	"github.com/berezovskiydeval/rpc-crud/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestProducts_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProductRepository(ctrl)
	productSvc := service.NewProducts(mockRepo)

	product := domain.Product{Id: 1, Name: "Test Product", Price: 100.0}
	mockRepo.EXPECT().Create(context.Background(), product).Return(int64(1), nil)

	Id, err := productSvc.Create(context.Background(), product)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), Id)
}

func TestProducts_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProductRepository(ctrl)
	productSvc := service.NewProducts(mockRepo)

	expected := []domain.Product{
		{Id: 1, Name: "Test", Price: 10},
		{Id: 2, Name: "Demo", Price: 20},
	}

	mockRepo.EXPECT().GetAll(context.Background()).Return(expected, nil)

	products, err := productSvc.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expected, products)
}

func TestProducts_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProductRepository(ctrl)
	productSvc := service.NewProducts(mockRepo)

	update := domain.ProductUpdate{Name: "Updated", Price: 150}
	mockRepo.EXPECT().Update(context.Background(), int64(1), update).Return(int64(1), nil)

	Id, err := productSvc.Update(context.Background(), 1, update)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), Id)
}

func TestProducts_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProductRepository(ctrl)
	productSvc := service.NewProducts(mockRepo)

	mockRepo.EXPECT().Delete(context.Background(), int64(1)).Return(int64(1), nil)

	Id, err := productSvc.Delete(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), Id)
}
