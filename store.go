package document

import (
	"context"
	"errors"
)

var (
	ErrDocumentFound    = errors.New("document: document found")
	ErrDocumentNotFound = errors.New("document: document not found")
)

type ReadWriter[T Document] interface {
	Reader[T]
	Writer
}

type Writer interface {
	Storer
	Creater
	Updater
	Deleter
}

type Reader[T Unmarshaler] interface {
	Fetcher[T]
	Lister[T]
}

type Storer interface {
	StoreDocument(ctx context.Context, b []byte, v ...any) error
}

type Creater interface {
	CreateDocument(ctx context.Context, b []byte, v ...any) error
}

type Fetcher[T Unmarshaler] interface {
	FetchDocument(ctx context.Context, b []byte, v ...[]byte) ([]T, error)
	FetchDocumentFn(ctx context.Context, b []byte, factory func() (T, error), v ...[]byte) ([]T, error)
}

type Updater interface {
	UpdateDocument(ctx context.Context, b []byte, v ...any) error
}

type Deleter interface {
	DeleteDocument(ctx context.Context, b []byte, v ...[]byte) error
}

type Lister[T Unmarshaler] interface {
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
