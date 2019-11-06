package packages

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

var (
	testJson         = `{"name":"test","version":"1.0.0","dependencies":[{"name":"test","version":"1.0.1","dependencies":null}]}`
	testInputJson    = `{"name":"test","version":"1.0.0","dependencies":{"test":"1.0.1"}}`
	testBadInputJson = `{"name":"test","version":"1.0.0","dependencies":"test"}`
	testNpm          = &Npm{
		Name:    "test",
		Version: "1.0.0",
		Dependencies: []Pkg{
			&Npm{
				Name:    "test",
				Version: "1.0.1",
			},
		},
		DependenciesRaw: map[string]string{
			"test": "1.0.1",
		},
	}
)

func TestNewNpm(t *testing.T) {
	type args struct {
		name    string
		version string
	}
	tests := []struct {
		name string
		args args
		want *Npm
	}{
		{
			name: "valid",
			args: args{"test", "1"},
			want: &Npm{Name: "test", Version: "1.0.0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNpm(tt.args.name, tt.args.version); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNpm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNpm_BuildTree(t *testing.T) {
	type fields struct {
		Name         string
		Version      string
		Dependencies []Pkg
	}
	type args struct {
		packages map[string]Pkg
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Pkg
	}{
		{
			name: "valid",
			args: args{
				packages: map[string]Pkg{
					"npm-test-1.0.1": &Npm{
						Name:    "test",
						Version: "1.0.0",
					},
				},
			},
			fields: fields{
				Name:         "test",
				Version:      "1.0.0",
				Dependencies: testNpm.Dependencies,
			},
			want: testNpm.Dependencies,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Npm{
				Name:         tt.fields.Name,
				Version:      tt.fields.Version,
				Dependencies: tt.fields.Dependencies,
			}
			n.BuildTree(tt.args.packages)
			got := n.Dependencies
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCacheKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNpm_GetCacheKey(t *testing.T) {
	type fields struct {
		Name    string
		Version string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "valid",
			fields: fields{"test", "1.0.0"},
			want:   "npm-test-1.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Npm{
				Name:    tt.fields.Name,
				Version: tt.fields.Version,
			}
			if got := n.GetCacheKey(); got != tt.want {
				t.Errorf("GetCacheKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNpm_GetURL(t *testing.T) {
	type fields struct {
		Name    string
		Version string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr error
	}{
		{
			name:   "valid",
			fields: fields{"test", "1.0.0"},
			want:   "https://registry.npmjs.org/test/1.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Npm{
				Name:    tt.fields.Name,
				Version: tt.fields.Version,
			}
			got, err := n.GetURL()
			if (err != nil) != (tt.wantErr != nil) || err != tt.wantErr {
				t.Errorf("GetURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.String(), tt.want) {
				t.Errorf("GetURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNpm_MarshalJSON(t *testing.T) {
	type fields struct {
		Name           string
		Version        string
		DependencyTree []Pkg
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
				Name:    "test",
				Version: "1.0.0",
				DependencyTree: []Pkg{
					&Npm{
						Name:    "test",
						Version: "1.0.1",
					},
				},
			},
			want: []byte(testJson),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Npm{
				Name:           tt.fields.Name,
				Version:        tt.fields.Version,
				DependencyTree: tt.fields.DependencyTree,
			}
			got, err := n.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNpm_ParseBody(t *testing.T) {
	type fields struct {
	}
	type args struct {
		body []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		want    *Npm
	}{
		{
			name: "valid",
			args: args{body: []byte(testInputJson)},
			want: testNpm,
		},
		{
			name:    "invalid",
			args:    args{body: []byte(testBadInputJson)},
			wantErr: errors.New("json"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Npm{}
			err := n.ParseBody(tt.args.body)
			if (err != nil) != (tt.wantErr != nil) || err != tt.wantErr {
				if err != nil && !strings.Contains(err.Error(), tt.wantErr.Error()) {
					t.Errorf("ParseBody() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !reflect.DeepEqual(n, tt.want) {
				t.Errorf("ParseBody() got = %v, want %v", n, tt.want)
			}
		})
	}
}

func TestNpm_parseVersion(t *testing.T) {
	type fields struct {
		Version string
	}
	var tests = []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "*",
			fields: fields{"*"},
			want:   NpmVersionLatest,
		},
		{
			name:   "only major",
			fields: fields{"1"},
			want:   "1.0.0",
		},
		{
			name:   "or",
			fields: fields{"1 || 2"},
			want:   "2.0.0",
		},
		{
			name:   "tilde",
			fields: fields{Version: "~1.0.0"},
			want:   "1.0.0",
		},
		{
			name:   "more than",
			fields: fields{">=1.0.0"},
			want:   "1.0.0",
		},
		{
			name:   "named version",
			fields: fields{"1.0.0-rce"},

			want: "1.0.0",
		},
		{
			name:   "major and minor",
			fields: fields{"1.0"},
			want:   "1.0.0",
		},
		{
			name:   "wildcard",
			fields: fields{"1.0.*"},
			want:   "1.0.0",
		},
		{
			name:   "invalid",
			fields: fields{"0.0.0"},
			want:   "0.0.1",
		},
		{
			name:   "space",
			fields: fields{">2 <=3"},
			want:   "3.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Npm{
				Version: tt.fields.Version,
			}
			n.parseVersion()
			got := n.Version

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}
