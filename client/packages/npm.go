package packages

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

const (
	npmPackageUrlFmt = "https://registry.npmjs.org/%s/%s"
	NpmVersionLatest = "latest"
)

type Npm struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	DependenciesRaw map[string]string `json:"dependencies"`
	Dependencies    []Pkg             `json:"-"`
	DependencyTree  []Pkg
	TreeBuilt       bool
}

func NewNpm(name, version string) *Npm {
	n := &Npm{Name: name, Version: version}
	n.parseVersion()
	return n
}

func (n *Npm) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name         string `json:"name"`
		Version      string `json:"version"`
		Dependencies []Pkg  `json:"dependencies"`
	}{
		Name:         n.Name,
		Version:      n.Version,
		Dependencies: n.DependencyTree,
	})
}

func (n *Npm) GetJSONKey() string {
	return "npm_package"
}

func (n *Npm) parseVersion() {
	if n.Version == "*" || n.Version == NpmVersionLatest || n.Version == "" {
		n.Version = NpmVersionLatest
		return
	}

	orSplit := strings.Split(n.Version, "||")
	if len(orSplit) > 1 {
		n.Version = strings.TrimSpace(orSplit[1])
	}

	toSplit := strings.Split(n.Version, "-")
	if len(toSplit) > 1 {
		n.Version = strings.TrimSpace(toSplit[0])
	}

	spaceSplit := strings.Split(n.Version, " ")
	if len(spaceSplit) > 1 {
		n.Version = strings.TrimSpace(spaceSplit[1])
	}

	replacer := strings.NewReplacer(
		"~", "",
		"^", "",
		">", "",
		"<", "",
		"x", "0",
		"=", "",
	)
	n.Version = replacer.Replace(n.Version)

	if len(strings.Split(n.Version, ".")) < 2 {
		n.Version += ".0.0"
	}

	if len(strings.Split(n.Version, ".")) < 3 {
		n.Version += ".0"
	}

	n.Version = strings.Replace(n.Version, "*", "0", -1)
	if n.Version == "0.0.0" {
		n.Version = "0.0.1"
	}
}

func (n *Npm) GetURL() (*url.URL, error) {
	rawUrl := fmt.Sprintf(npmPackageUrlFmt, n.Name, n.Version)
	return url.Parse(rawUrl)
}

func (n *Npm) GetCacheKey() string {
	return fmt.Sprintf("npm-%s-%s", n.Name, n.Version)
}

func (n *Npm) ParseBody(body []byte) error {
	var parsedPkg Npm
	err := json.Unmarshal(body, &parsedPkg)
	if err != nil {
		return err
	}

	n.Name = parsedPkg.Name
	n.Version = parsedPkg.Version
	n.DependenciesRaw = parsedPkg.DependenciesRaw

	for lib, version := range n.DependenciesRaw {
		n.Dependencies = append(n.Dependencies, NewNpm(lib, version))
	}
	return nil
}

func (n Npm) GetDependencies() []Pkg {
	return n.Dependencies
}

func (n *Npm) BuildTree(packages map[string]Pkg) {
	if n.TreeBuilt {
		return
	}
	n.TreeBuilt = true

	for _, pkg := range n.Dependencies {
		p, ok := packages[pkg.GetCacheKey()]
		if ok {
			basePkg := p
			basePkg.BuildTree(packages)
			n.DependencyTree = append(n.DependencyTree, basePkg)
		}
	}
}

func (n *Npm) Fetch(fetchFunc func() ([]byte, error)) ([]byte, error) {
	return fetchFunc()
}
