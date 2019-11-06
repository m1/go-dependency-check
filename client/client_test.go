package client

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m1/go-dependency-check/client/packages"
	"github.com/m1/go-dependency-check/worker"
)

var (
	testMock = packages.Mock{
		Name:    "test",
		Version: "1.0.0",
	}

	mockRepo = map[string]packages.Mock{}
)

func NewTestClient() Client {
	mockRepo = map[string]packages.Mock{
		"test-1-1.0.0": packages.Mock{
			Name:    "test-1",
			Version: "1.0.0",
			RawBody: []byte(`{"name":"test-1","version":"1.0.0","dependencies":[{"name":"test-2","version":"1.0.0"}]}`),
			Dependencies: []packages.Pkg{
				packages.Mock{
					Name:    "test-2",
					Version: "1.0.0",
					RawBody: []byte(`{"name":"test-1","version":"1.0.0","dependencies":[{"name":"test-3","version":"1.0.0"}]}`),
				},
			},
		},
		"test-2-1.0.0": packages.Mock{
			Name:    "test-2",
			Version: "1.0.0",
			RawBody: []byte(`{"name":"test-1","version":"1.0.0","dependencies":[{"name":"test-3","version":"1.0.0"}]}`),
			Dependencies: []packages.Pkg{
				packages.Mock{
					Name:    "test-3",
					Version: "1.0.0",
					RawBody: []byte(`{"name":"test-3","version":"1.0.0"`),
				},
			},
		},
		"test-3-1.0.0": packages.Mock{
			Name:         "test-3",
			Version:      "1.0.0",
			RawBody:      []byte(`{"name":"test-3","version":"1.0.0"`),
			Dependencies: []packages.Pkg{},
		},
	}
	return Client{
		Queued:   &sync.Map{},
		Pool:     worker.NewPool(2),
		Packages: make(map[string]packages.Pkg),
	}
}

func TestClient_GetDependencyTree(t *testing.T) {
	type args struct {
		pkg packages.Pkg
	}
	tests := []struct {
		name string
		args args
		test func(c Client)
	}{
		{
			name: "valid",
			args: args{
				pkg: testMock,
			},
			test: func(c Client) {
				assert.Equal(t, 1, c.JobsCreated)
				assert.Equal(t, 1, c.JobsCompleted)

				for key := range c.Packages {
					if key != "test-1.0.0" && key != "test-1.0.1" {
						t.Errorf("GetURL() expecting proper version numbers")
					}
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewTestClient()
			c.GetDependencyTree(testMock)
			tt.test(c)
		})
	}
}
