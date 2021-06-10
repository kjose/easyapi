// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utils

import "sync"

type Context struct {
	Context sync.Map
}

func NewContext() *Context {
	return &Context{}
}

func (s *Context) With(key string, value interface{}) *Context {
	s.Set(key, value)
	return s
}

func (s *Context) Set(key string, value interface{}) {
	s.Context.Store(key, value)
}

func (s *Context) Get(key string) interface{} {
	r, ok := s.Context.Load(key)
	if !ok {
		return nil
	}
	return r
}
