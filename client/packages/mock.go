package packages

import (
	"fmt"
	"net/url"
)

type Mock struct {
	Packages       map[string]Mock
	MockRepo       map[string]Pkg
	Name           string
	Version        string
	RawBody        []byte
	Dependencies   []Pkg
	DependencyTree []Pkg
}

func (m Mock) GetJSONKey() string {
	return "mock_package"
}

func (m Mock) GetURL() (*url.URL, error) {
	return nil, nil
}

func (m Mock) ParseBody([]byte) error {
	return nil
}

func (m Mock) GetDependencies() []Pkg {
	return m.Dependencies
}

func (m Mock) GetCacheKey() string {
	return fmt.Sprintf("%s-%s", m.Name, m.Version)
}

func (m Mock) BuildTree(packages map[string]Pkg) {
	return
}

func (m Mock) Fetch(func() ([]byte, error)) ([]byte, error) {
	return m.RawBody, nil
}
