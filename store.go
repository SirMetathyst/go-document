package document

import (
	"context"
	"errors"
)

var (
	ErrDocumentFound    = errors.New("document: document found")
	ErrDocumentNotFound = errors.New("document: document not found")
)

type Store[T Document] interface {
	Storer
	Creater
	Reader[T]
	Updater
	Deleter
	Lister[T]
}

type Storer interface {
	StoreDocument(ctx context.Context, b []byte, v ...Marshaler) error
	//StoreDocumentFn(ctx context.Context, b []byte, fn func(ctx PutContext) error) error
}

type Creater interface {
	CreateDocument(ctx context.Context, b []byte, v ...Marshaler) error
	//CreateDocumentFn(ctx context.Context, b []byte, fn func(ctx PutContext) error) error
}

type Reader[T Document] interface {
	ReadDocument(ctx context.Context, b []byte, v ...[]byte) ([]T, error)
	ReadDocumentFn(ctx context.Context, b []byte, factory func() (T, error), v ...[]byte) ([]T, error)
	//ReadDocumentFn(ctx context.Context, b []byte, fn func(ctx GetContext) error) error
}

type Updater interface {
	UpdateDocument(ctx context.Context, b []byte, v ...Marshaler) error
	//ReadDocumentFn(ctx context.Context, b []byte, fn func(ctx GetContext) error) error
}

type Deleter interface {
	DeleteDocument(ctx context.Context, b []byte, v ...[]byte) error
}

type Lister[T Document] interface {
	ListDocument(ctx context.Context, b []byte) ([]T, error)
	ListDocumentFn(ctx context.Context, b []byte, factory func() (T, error)) ([]T, error)
}

//type PutContext interface {
//	context.Context
//	Put(v Marshaler) error
//}

type Document interface {
	Marshaler
	Unmarshaler
}

type Marshaler interface {
	MarshalDocument() ([]byte, []byte, error)
}

type Unmarshaler interface {
	UnmarshalDocument([]byte, []byte) error
}

type Bucket []byte

type Key []byte
