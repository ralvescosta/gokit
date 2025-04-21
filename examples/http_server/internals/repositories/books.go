// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package repositories

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type (
	bookRepository struct {
		trace trace.Tracer
	}
)

func NewBookRepository() *bookRepository {
	return &bookRepository{
		otel.Tracer("BookRepository"),
	}
}

func (r *bookRepository) Create(ctx context.Context) {
	_, span := r.trace.Start(ctx, "CreateBook")
	time.Sleep(time.Microsecond * 100)
	defer span.End()
}

func (r *bookRepository) Get(ctx context.Context) {
	_, span := r.trace.Start(ctx, "GetBook")
	time.Sleep(time.Microsecond * 40)
	defer span.End()
}

func (r *bookRepository) List(ctx context.Context) {
	_, span := r.trace.Start(ctx, "ListBook")
	time.Sleep(time.Microsecond * 60)
	defer span.End()
}

func (r *bookRepository) Update(ctx context.Context) {
	_, span := r.trace.Start(ctx, "UpdateBook")
	time.Sleep(time.Microsecond * 40)
	defer span.End()
}

func (r *bookRepository) Delete(ctx context.Context) {
	_, span := r.trace.Start(ctx, "DeleteBook")
	time.Sleep(time.Microsecond * 50)
	defer span.End()
}
