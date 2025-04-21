// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package services

import "context"

type (
	BookRepository interface {
		Create(ctx context.Context)
		Get(ctx context.Context)
		List(ctx context.Context)
		Update(ctx context.Context)
		Delete(ctx context.Context)
	}
)
