// Copyright 2020 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpcsetting

import (
	"crypto/tls"
	"net/http"
	"testing"

	"google.golang.org/api/internal"
	"google.golang.org/grpc"
)

func TestDialSettings_Validate(t *testing.T) {
	dummyGetClientCertificate := func(info *tls.CertificateRequestInfo) (*tls.Certificate, error) { return nil, nil }
	type option interface {
		Apply(*internal.DialSettings)
	}
	tests := []struct {
		name  string
		ds    internal.DialSettings
		opts  []option
		valid bool
	}{
		{
			name:  "GRPCConn -- valid",
			ds:    internal.DialSettings{},
			opts:  []option{GRPCConnSetting{&grpc.ClientConn{}}},
			valid: true,
		},
		{
			name:  "GRPCConn -- valid",
			ds:    internal.DialSettings{HTTPClient: &http.Client{}},
			opts:  []option{GRPCConnSetting{&grpc.ClientConn{}}},
			valid: false,
		},
		{
			name:  "GRPCConn -- valid",
			ds:    internal.DialSettings{ClientCertSource: dummyGetClientCertificate},
			opts:  []option{GRPCConnSetting{&grpc.ClientConn{}}},
			valid: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, opt := range tc.opts {
				opt.Apply(&tc.ds)
			}
			err := tc.ds.Validate()
			if tc.valid && err != nil {
				t.Errorf("got %v, want nil", err)
			}
			if !tc.valid && err == nil {
				t.Errorf("got nil, want an error")
			}
		})
	}

}
