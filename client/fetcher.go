package client

import (
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/m1/go-dependency-check/cache"
	"github.com/m1/go-dependency-check/client/packages"
)

// Fetcher is the type that actually does the fetching of the given package
type Fetcher struct {
	// Pkg is the given package to f etch
	Pkg       packages.Pkg

	// Results is the fetcher channel
	Results   chan Fetcher

	// Queue is the queue of packages to fetch
	Queue     chan packages.Pkg

	// Queued is the map of packages already queued/visited
	Queued    *sync.Map

	// PkgType is the type of package, i.e npm/maven
	PkgType   packages.PkgType

	// Err is not nil if the fetcher occurs an error whilst fetching the package
	Err       error

	Cache     cache.Cache
	CacheTime *time.Duration
}

// Run starts the fetching of the given package
func (f Fetcher) Run() {
	body, err := f.fetchBody()
	if err != nil {
		f.Err = err
		f.Results <- f
		return
	}

	err = f.Pkg.ParseBody(body)
	if err != nil {
		f.Err = err
		f.Results <- f
		return
	}

	for _, pkg := range f.Pkg.GetDependencies() {
		// has already been queued for fetching
		queued, ok := f.Queued.Load(pkg.GetCacheKey())
		if !ok || !queued.(bool) {
			f.Queued.Store(pkg.GetCacheKey(), true)
			f.Queue <- pkg
		}
	}

	f.Results <- f
}

func (f Fetcher) fetchBody() ([]byte, error) {
	if f.Cache != nil {
		cached, err := f.Cache.Get(f.Pkg.GetCacheKey())

		if err != nil {
			body, err := f.newRequest()
			if err != nil {
				return nil, err
			}

			err = f.Cache.Set(f.Pkg.GetCacheKey(), string(body), *f.CacheTime)
			if err != nil {
				return nil, err
			}

			return body, nil
		} else {
			return []byte(*cached), nil
		}
	}

	return f.newRequest()
}

func (f Fetcher) newRequest() ([]byte, error) {
	return f.Pkg.Fetch(f.fetchURL)
}

func (f Fetcher) fetchURL() ([]byte, error) {
	url, err := f.Pkg.GetURL()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 400 {
		return nil, errors.New("couldn't fetch package")
	}

	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
