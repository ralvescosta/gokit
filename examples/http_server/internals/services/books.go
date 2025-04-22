// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package services

import "context"

type (
	BookService interface {
		RegisterBook(ctx context.Context)
		GetBook(ctx context.Context)
		ListBook(ctx context.Context)
		UpdateBook(ctx context.Context)
		DeleteBook(ctx context.Context)
	}

	bookService struct {
		repository BookRepository
	}
)

func NewBookService(repository BookRepository) BookService {
	return &bookService{repository}
}

func (s *bookService) RegisterBook(ctx context.Context) {
	s.repository.Create(ctx)
}

func (s *bookService) GetBook(ctx context.Context) {
	s.repository.Get(ctx)
}

func (s *bookService) ListBook(ctx context.Context) {
	s.repository.List(ctx)
}

func (s *bookService) UpdateBook(ctx context.Context) {
	s.repository.Update(ctx)
}

func (s *bookService) DeleteBook(ctx context.Context) {
	s.repository.Delete(ctx)
}
