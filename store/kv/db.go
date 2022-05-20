package kv

import (
	"context"
	"encoding/json"
	"github.com/SirMetathyst/go-document"
	"github.com/SirMetathyst/go-kv"
	"reflect"
)

var _ document.Store[*Document] = &DB[*Document]{}

type DB[T document.Document] struct{ kv.Store }

func (s *DB[T]) StoreDocument(ctx context.Context, b []byte, v ...document.Marshaler) error {

	if len(b) == 0 || len(v) == 0 {
		return nil
	}

	return s.StoreKVFn(ctx, b, func(ctx kv.PutContext) error {
		return putFor(ctx, v)
	})
}

func (s *DB[T]) CreateDocument(ctx context.Context, b []byte, v ...document.Marshaler) error {

	if len(b) == 0 || len(v) == 0 {
		return nil
	}

	return interceptError(s.CreateKVFn(ctx, b, func(ctx kv.PutContext) error {
		return putFor(ctx, v)
	}))
}

func (s *DB[T]) ReadDocument(ctx context.Context, b []byte, v ...[]byte) (list []T, err error) {

	if len(b) == 0 || len(v) == 0 {
		return nil, nil
	}

	return list, interceptError(s.ReadKVFn(ctx, b, func(ctx kv.GetContext) error {

		i := 0
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				value, err := ctx.Get(v[i], false)
				if err != nil {
					return err
				}
				var resultType T
				result := newInstance(resultType).(T)
				if err = result.UnmarshalDocument(v[i], value); err != nil {
					return err
				}
				list = append(list, result)
				i++
				if i >= len(v) {
					return nil
				}
			}
		}
	}))
}

func (s *DB[T]) ReadDocumentFn(ctx context.Context, b []byte, factory func() (T, error), v ...[]byte) (list []T, err error) {

	if len(b) == 0 || len(v) == 0 {
		return nil, nil
	}

	return list, interceptError(s.ReadKVFn(ctx, b, func(ctx kv.GetContext) error {
		i := 0
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				value, err := ctx.Get(v[i], false)
				if err != nil {
					return err
				}
				result, err := factory()
				if err != nil {
					return err
				}
				if err = result.UnmarshalDocument(v[i], value); err != nil {
					return err
				}
				list = append(list, result)
				i++
				if i >= len(v) {
					return nil
				}
			}
		}
	}))
}

func (s *DB[T]) UpdateDocument(ctx context.Context, b []byte, v ...document.Marshaler) error {

	if len(b) == 0 || len(v) == 0 {
		return nil
	}

	return interceptError(s.UpdateKVFn(ctx, b, func(ctx kv.PutContext) error {
		return putFor(ctx, v)
	}))
}

func (s *DB[T]) DeleteDocument(ctx context.Context, b []byte, v ...[]byte) error {

	if len(b) == 0 || len(v) == 0 {
		return nil
	}

	return s.DeleteKVFn(ctx, b, func(ctx kv.DeleteContext) error {
		i := 0
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				err := ctx.Delete(v[i])
				if err != nil {
					return err
				}
				i++
				if i >= len(v) {
					return nil
				}
			}
		}
	})
}

func (s *DB[T]) ListDocument(ctx context.Context, b []byte) (list []T, err error) {

	if len(b) == 0 {
		return nil, nil
	}

	return list, s.ListKVFn(ctx, b, func(k []byte, v []byte) error {
		i := 0
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				var resultType T
				result := newInstance(resultType).(T)
				if err = result.UnmarshalDocument(k, v); err != nil {
					return err
				}
				list = append(list, result)
				i++
				if i >= len(v) {
					return nil
				}
			}
		}
	})
}

func (s *DB[T]) ListDocumentFn(ctx context.Context, b []byte, factory func() (T, error)) (list []T, err error) {

	if len(b) == 0 {
		return nil, nil
	}

	return list, s.ListKVFn(ctx, b, func(k []byte, v []byte) error {
		i := 0
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				result, err := factory()
				if err != nil {
					return err
				}
				if err = result.UnmarshalDocument(k, v); err != nil {
					return err
				}
				list = append(list, result)
				i++
				if i >= len(v) {
					return nil
				}
			}
		}
	})
}

func putFor(ctx kv.PutContext, v []document.Marshaler) error {

	if len(v) == 0 {
		return nil
	}

	i := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			key, value, err := v[i].MarshalDocument()
			if err != nil {
				return err
			}
			if err := ctx.Put(key, value); err != nil {
				return err
			}
			i++
			if i >= len(v) {
				return nil
			}
		}
	}
}

func interceptError(err error) error {

	if err == kv.ErrKeyFound {
		return document.ErrDocumentFound
	} else if err == kv.ErrKeyNotFound {
		return document.ErrDocumentNotFound
	}

	return err
}

func newInstance[T any](v T) any {
	if typ := reflect.TypeOf(v); typ.Kind() == reflect.Ptr {
		elem := typ.Elem()
		return reflect.New(elem).Interface()
	}
	return new(T)
}

type Document struct {
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
}

func (s *Document) MarshalDocument() ([]byte, []byte, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, nil, err
	}
	return s.Key, v, nil
}

func (s *Document) UnmarshalDocument(k []byte, v []byte) error {
	var doc Document
	err := json.Unmarshal(v, &doc)
	if err != nil {
		return err
	}
	s.Key = k
	s.Value = doc.Value
	return nil
}

func New[T document.Document](store kv.Store) (*DB[T], error) {
	return &DB[T]{store}, nil
}
