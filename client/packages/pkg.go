package packages

import (
	"net/url"

	"github.com/m1/go-dependency-check/api/response"
)

const (
	PkgTypeNpm PkgType = iota
	PkgTypeMaven
)

// Pkg is the package interface
type Pkg interface {
	response.Model

	// GetURL returns the url to fetch for the package info
	GetURL() (*url.URL, error)

	// ParseBody parses the body of the response for the package info
	ParseBody([]byte) error

	// GetDependencies returns the dependencies of this package
	GetDependencies() []Pkg

	// GetCacheKey returns the key used for caching this package
	GetCacheKey() string

	// BuildTree builds the depedency tree for this package
	BuildTree(packages map[string]Pkg)

	// Fetch returns the raw data of the package
	Fetch(func() ([]byte, error)) ([]byte, error)
}

// PkgType is the type of package, i.e npm/maven
type PkgType int
