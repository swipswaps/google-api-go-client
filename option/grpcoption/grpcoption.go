// Copyright 2020 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpcoption

import (
	"google.golang.org/api/internal"
	"google.golang.org/api/internal/setting/grpcsetting"
	"google.golang.org/grpc"
)

// A ClientOption is an option for a Google API client.
type ClientOption interface {
	Apply(*internal.DialSettings)
}

// WithGRPCConn returns a ClientOption that specifies the gRPC client
// connection to use as the basis of communications. This option may only be
// used with services that support gRPC as their communication transport. When
// used, the WithGRPCConn option takes precedent over all other supplied
// options.
func WithGRPCConn(conn *grpc.ClientConn) ClientOption {
	return grpcsetting.GRPCConnSetting{conn}
}

// WithGRPCDialOption returns a ClientOption that appends a new grpc.DialOption
// to an underlying gRPC dial. It does not work with WithGRPCConn.
func WithGRPCDialOption(opt grpc.DialOption) ClientOption {
	return withGRPCDialOption{opt}
}

type withGRPCDialOption struct{ opt grpc.DialOption }

func (w withGRPCDialOption) Apply(o *internal.DialSettings) {
	o.GRPCDialOpts = append(o.GRPCDialOpts, w.opt)
}

// WithGRPCConnectionPool returns a ClientOption that creates a pool of gRPC
// connections that requests will be balanced between.
//
// This is an EXPERIMENTAL API and may be changed or removed in the future.
func WithGRPCConnectionPool(size int) ClientOption {
	return withGRPCConnectionPool(size)
}

type withGRPCConnectionPool int

func (w withGRPCConnectionPool) Apply(o *internal.DialSettings) {
	o.GRPCConnPoolSize = int(w)
}
