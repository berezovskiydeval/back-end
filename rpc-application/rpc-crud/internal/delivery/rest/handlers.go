package rest

import (
	"context"
	"net/http"
	"github.com/berezovskiydeval/rpc-crud/internal/domain"

	"github.com/gorilla/mux"
)

type ProductService interface {
	Create(ctx context.Context, product domain.Product) (int64, error)
	GetAll(ctx context.Context) ([]domain.Product, error)
	Update(ctx context.Context, id int64, product domain.ProductUpdate) (int64, error)
	Delete(ctx context.Context, id int64) (int64, error)
}

type Handler struct{
	productService ProductService
}

func NewHandler(productService ProductService) *Handler{
	return &Handler{
		productService: productService,
	}
}

func (h *Handler) InitRoutes() http.Handler{
	r := mux.NewRouter()

	r.Use(loggingMiddleware)
	products := r.PathPrefix("/products").Subrouter()
	{
		products.HandleFunc("", h.Create).Methods(http.MethodPost)
		products.HandleFunc("/{id}", h.Update).Methods(http.MethodPut)
		products.HandleFunc("/{id}", h.Delete).Methods(http.MethodDelete)
		products.HandleFunc("", h.GetAll).Methods(http.MethodGet)
	}

	return r
}