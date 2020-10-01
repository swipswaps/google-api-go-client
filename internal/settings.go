// Copyright 2017 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package internal supports the options and transport packages.
package internal

import (
	"crypto/tls"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/internal/impersonate"
	"google.golang.org/grpc"
)

const (
	GRPCConnSettingKey = SettingKey("GRPCConnSettingKey")
)

// SettingKey is a value used to retrieve a setting from DialSettings. Must be
// unique.
type SettingKey string

// DialSettings holds information needed to establish a connection with a
// Google API service.
type DialSettings struct {
	Endpoint            string
	DefaultEndpoint     string
	DefaultMTLSEndpoint string
	Scopes              []string
	TokenSource         oauth2.TokenSource
	Credentials         *google.Credentials
	CredentialsFile     string // if set, Token Source is ignored.
	CredentialsJSON     []byte
	UserAgent           string
	APIKey              string
	Audiences           []string
	HTTPClient          *http.Client
	GRPCDialOpts        []grpc.DialOption
	GRPCConnPool        ConnPool
	GRPCConnPoolSize    int
	NoAuth              bool
	TelemetryDisabled   bool
	ClientCertSource    func(*tls.CertificateRequestInfo) (*tls.Certificate, error)
	CustomClaims        map[string]interface{}
	SkipValidation      bool
	ImpersonationConfig *impersonate.Config

	// Google API system parameters. For more information please read:
	// https://cloud.google.com/apis/docs/system-parameters
	QuotaProject  string
	RequestReason string

	m map[SettingKey]interface{}
}

// SetSetting sets a setting to a value. Should not be called directly but
// instead from a setting's Apply method.
func (ds *DialSettings) SetSetting(k SettingKey, v interface{}) {
	if ds.m == nil {
		ds.m = make(map[SettingKey]interface{})
	}
	ds.m[k] = v
}

// GetSetting gets a setting for a given key. Should not be called directly but
// instead from Get helpers that return the proper types.
func (ds *DialSettings) GetSetting(k SettingKey) (v interface{}, ok bool) {
	if ds.m == nil {
		ds.m = make(map[SettingKey]interface{})
	}
	v, ok = ds.m[k]
	return
}

// IsSet checks if a setting is set.
func (ds *DialSettings) IsSet(k SettingKey) bool {
	_, ok := ds.m[k]
	return ok
}

// Validate reports an error if ds is invalid.
func (ds *DialSettings) Validate() error {
	if ds.SkipValidation {
		return nil
	}
	hasCreds := ds.APIKey != "" || ds.TokenSource != nil || ds.CredentialsFile != "" || ds.Credentials != nil
	if ds.NoAuth && hasCreds {
		return errors.New("options.WithoutAuthentication is incompatible with any option that provides credentials")
	}
	// Credentials should not appear with other options.
	// We currently allow TokenSource and CredentialsFile to coexist.
	// TODO(jba): make TokenSource & CredentialsFile an error (breaking change).
	nCreds := 0
	if ds.Credentials != nil {
		nCreds++
	}
	if ds.CredentialsJSON != nil {
		nCreds++
	}
	if ds.CredentialsFile != "" {
		nCreds++
	}
	if ds.APIKey != "" {
		nCreds++
	}
	if ds.TokenSource != nil {
		nCreds++
	}
	if len(ds.Scopes) > 0 && len(ds.Audiences) > 0 {
		return errors.New("WithScopes is incompatible with WithAudience")
	}
	// Accept only one form of credentials, except we allow TokenSource and CredentialsFile for backwards compatibility.
	if nCreds > 1 && !(nCreds == 2 && ds.TokenSource != nil && ds.CredentialsFile != "") {
		return errors.New("multiple credential options provided")
	}
	if ds.IsSet(GRPCConnSettingKey) && ds.GRPCConnPool != nil {
		return errors.New("WithGRPCConn is incompatible with WithConnPool")
	}
	if ds.HTTPClient != nil && ds.GRPCConnPool != nil {
		return errors.New("WithHTTPClient is incompatible with WithConnPool")
	}
	if ds.HTTPClient != nil && ds.IsSet(GRPCConnSettingKey) {
		return errors.New("WithHTTPClient is incompatible with WithGRPCConn")
	}
	if ds.HTTPClient != nil && ds.GRPCDialOpts != nil {
		return errors.New("WithHTTPClient is incompatible with gRPC dial options")
	}
	if ds.HTTPClient != nil && ds.QuotaProject != "" {
		return errors.New("WithHTTPClient is incompatible with QuotaProject")
	}
	if ds.HTTPClient != nil && ds.RequestReason != "" {
		return errors.New("WithHTTPClient is incompatible with RequestReason")
	}
	if ds.HTTPClient != nil && ds.ClientCertSource != nil {
		return errors.New("WithHTTPClient is incompatible with WithClientCertSource")
	}
	if ds.ClientCertSource != nil && (ds.IsSet(GRPCConnSettingKey) || ds.GRPCConnPool != nil || ds.GRPCConnPoolSize != 0 || ds.GRPCDialOpts != nil) {
		return errors.New("WithClientCertSource is currently only supported for HTTP. gRPC settings are incompatible")
	}
	if ds.ImpersonationConfig != nil && len(ds.ImpersonationConfig.Scopes) == 0 && len(ds.Scopes) == 0 {
		return errors.New("WithImpersonatedCredentials requires scopes being provided")
	}
	return nil
}
