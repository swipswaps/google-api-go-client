// Copyright 2020 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpcsetting

import (
	"google.golang.org/api/internal"
	"google.golang.org/grpc"
)

type GRPCConnSetting struct{ Value *grpc.ClientConn }

func (s GRPCConnSetting) Apply(ds *internal.DialSettings) {
	ds.SetSetting(internal.GRPCConnSettingKey, s.Value)
}

func GetGRPCConnSetting(ds *internal.DialSettings) *grpc.ClientConn {
	v, ok := ds.GetSetting(internal.GRPCConnSettingKey)
	if !ok {
		return nil
	}
	return v.(*grpc.ClientConn)
}
