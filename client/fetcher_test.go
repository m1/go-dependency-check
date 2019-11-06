package client

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/m1/go-dependency-check/cache"
	"github.com/m1/go-dependency-check/client/packages"
)

type MockCache struct {
	Packages map[string]packages.Mock
}

func (m MockCache) Get(key string) (*string, error) {
	for key, pkg := range m.Packages {
		if pkg.GetCacheKey() == key {
			body := string(pkg.RawBody)
			return &body, nil
		}
	}

	return nil, cache.ErrCacheKeyNotFound
}

func (m MockCache) Set(key, data string, duration time.Duration) error {
	return nil
}

func TestFetcher_fetchBody(t *testing.T) {
	type fields struct {
		Pkg       packages.Pkg
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Pkg: mockRepo["test-1-1.0.0"],
			},
			want: mockRepo["test-1-1.0.0"].RawBody,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Fetcher{
				Pkg:       tt.fields.Pkg,
			}
			got, err := f.fetchBody()
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchBody() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetcher_fetchURL(t *testing.T) {
	// TODO: Add NewTestHTTPServer here
	type fields struct {
		Pkg       packages.Pkg
		Results   chan Fetcher
		Queue     chan packages.Pkg
		Queued    *sync.Map
		PkgType   packages.PkgType
		Err       error
		Cache     cache.Cache
		CacheTime *time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Fetcher{
				Pkg:       tt.fields.Pkg,
				Results:   tt.fields.Results,
				Queue:     tt.fields.Queue,
				Queued:    tt.fields.Queued,
				PkgType:   tt.fields.PkgType,
				Err:       tt.fields.Err,
				Cache:     tt.fields.Cache,
				CacheTime: tt.fields.CacheTime,
			}
			got, err := f.fetchURL()
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}