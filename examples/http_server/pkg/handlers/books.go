// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package handlers

import (
	"errors"
	"net/http"

	"github.com/ralvescosta/gokit/httpw/middlewares"
	"github.com/ralvescosta/gokit/httpw/server"
	"github.com/ralvescosta/gokit/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/ralvescosta/gokit/examples/http_server/internals/services"
	_ "github.com/ralvescosta/gokit/examples/http_server/pkg/openapi"
)

type (
	HTTPHandlers interface {
		Install(router server.HTTPServer)
	}

	httpHandlers struct {
		logger         logging.Logger
		service        services.BookService
		authMiddleware middlewares.Authorization
		tracer         trace.Tracer
	}
)

func NewHandler(logger logging.Logger, service services.BookService, authMiddleware middlewares.Authorization) HTTPHandlers {
	return &httpHandlers{logger, service, authMiddleware, otel.Tracer("http handler")}
}

func (h *httpHandlers) Install(router server.HTTPServer) {
	router.Group("/books", []*server.Route{
		server.NewRouteBuilder().Path("/").Method(http.MethodPost).Handler(h.postHandler).Build(),
		server.NewRouteBuilder().Path("/{id}").Method(http.MethodGet).Handler(h.getHandler).Middlewares(h.authMiddleware.Handle).Build(),
		server.NewRouteBuilder().Path("/").Method(http.MethodGet).Handler(h.listHandler).Middlewares(h.authMiddleware.Handle).Build(),
		server.NewRouteBuilder().Path("/{id}").Method(http.MethodPut).Handler(h.putHandler).Middlewares(h.authMiddleware.Handle).Build(),
		server.NewRouteBuilder().Path("/{id}").Method(http.MethodDelete).Handler(h.deleteHandler).Middlewares(h.authMiddleware.Handle).Build(),
	})
}

// CreateBook
// @Summary      CreateBook
// @Description  Create a new Book
// @Tags         books
// @Accept       json
// @Produce      json
// @Param				 books	body		httpHandlers.Book	true	"Add Book"
// @Success      201  {object}  httpHandlers.Book
// @Failure      400  {object}  httpHandlers.HTTPError
// @Failure      404  {object}  httpHandlers.HTTPError
// @Failure      500  {object}  httpHandlers.HTTPError
// @Router       /books/ [post]
func (h *httpHandlers) postHandler(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("postHandler")
	h.service.RegisterBook(req.Context())

	w.WriteHeader(http.StatusCreated)
}

// GetBook
// @Summary      GetBook
// @Description  Get Book By ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      200  {object}  httpHandlers.Book
// @Failure      400  {object}  httpHandlers.HTTPError
// @Failure      404  {object}  httpHandlers.HTTPError
// @Failure      500  {object}  httpHandlers.HTTPError
// @Router       /books/{id} [get]
func (h *httpHandlers) getHandler(w http.ResponseWriter, req *http.Request) {
	_, span := h.tracer.Start(req.Context(), "getHandler")
	defer span.End()

	span.SetStatus(codes.Error, "some error")
	span.RecordError(errors.New("some error"))

	h.logger.Info("getHandler")
	h.service.GetBook(req.Context())

	w.WriteHeader(200)
}

func (h *httpHandlers) listHandler(w http.ResponseWriter, req *http.Request) {
	_, span := h.tracer.Start(req.Context(), "listHandler")
	defer span.End()

	h.logger.Info("listHandler")
	h.service.ListBook(req.Context())

	w.WriteHeader(200)
	w.Write([]byte("oi"))
}

func (h *httpHandlers) putHandler(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("putHandler")
	h.service.UpdateBook(req.Context())

	w.WriteHeader(200)
}

func (h *httpHandlers) deleteHandler(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("deleteHandler")
	h.service.DeleteBook(req.Context())

	w.WriteHeader(200)
}
