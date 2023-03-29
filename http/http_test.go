package http

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuildClientInvalidTlsVersion(t *testing.T) {
	_, err := BuildClient(60, 9)
	assert.Equal(t, "unsupported tls version 9", err.Error())
}

func TestBuildClient(t *testing.T) {
	type testCase struct {
		name     string
		inputs   []interface{}
		expected *HttpClientDecorator
	}
	testCases := []testCase{
		{
			name:     "tls version 1.0",
			inputs:   []interface{}{uint32(25), 10},
			expected: &HttpClientDecorator{Timeout: 25 * time.Second, TlsVersion: tls.VersionTLS10, client: http.Client{}},
		},
		{
			name:     "tls version 1.1",
			inputs:   []interface{}{uint32(50), 11},
			expected: &HttpClientDecorator{Timeout: 50 * time.Second, TlsVersion: tls.VersionTLS11, client: http.Client{}},
		},
		{
			name:     "tls version 1.2",
			inputs:   []interface{}{uint32(90), 12},
			expected: &HttpClientDecorator{Timeout: 90 * time.Second, TlsVersion: tls.VersionTLS12, client: http.Client{}},
		},
		{
			name:     "tls version 1.3",
			inputs:   []interface{}{uint32(14), 13},
			expected: &HttpClientDecorator{Timeout: 14 * time.Second, TlsVersion: tls.VersionTLS13, client: http.Client{}},
		},
	}

	for _, tc := range testCases {
		actual, err := BuildClient(tc.inputs[0].(uint32), tc.inputs[1].(int))
		assert.Nil(t, err, err)
		assert.Equal(t, tc.expected.Timeout, actual.Timeout)
		assert.Equal(t, tc.expected.TlsVersion, actual.TlsVersion)
		assert.NotNil(t, actual.client)
	}
}

func TestBuildClientTls10(t *testing.T) {
	actual, err := BuildClientTls10()
	assert.Nil(t, err, err)
	assert.Equal(t, 30*time.Second, actual.Timeout)
	assert.Equal(t, uint16(tls.VersionTLS10), actual.TlsVersion)
}

func TestBuildClientTls11(t *testing.T) {
	actual, err := BuildClientTls11()
	assert.Nil(t, err, err)
	assert.Equal(t, 30*time.Second, actual.Timeout)
	assert.Equal(t, uint16(tls.VersionTLS11), actual.TlsVersion)
}

func TestBuildClientTls12(t *testing.T) {
	actual, err := BuildClientTls12()
	assert.Nil(t, err, err)
	assert.Equal(t, 30*time.Second, actual.Timeout)
	assert.Equal(t, uint16(tls.VersionTLS12), actual.TlsVersion)
}

func TestBuildClientTls13(t *testing.T) {
	actual, err := BuildClientTls13()
	assert.Nil(t, err, err)
	assert.Equal(t, 30*time.Second, actual.Timeout)
	assert.Equal(t, uint16(tls.VersionTLS13), actual.TlsVersion)
}
