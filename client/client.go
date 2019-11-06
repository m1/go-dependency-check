package client

import (
	"errors"
	"sync"
	"time"

	"github.com/m1/go-dependency-check/cache"
	"github.com/m1/go-dependency-check/client/packages"
	"github.com/m1/go-dependency-check/worker"
)

var (
	// ErrFetchingRootPackage is the error if the first given package cannot be
	// found
	ErrFetchingRootPackage = errors.New("couldn't find root package")
)

type Client struct {
	ClientConfig

	// Queued is a list of already cralwed/to crawl packages
	Queued        *sync.Map

	// Pool is a pool of workers for fetching packages
	Pool          *worker.Pool

	// Packages a map of packages fetched
	Packages      map[string]packages.Pkg

	JobsCreated   int
	JobsCompleted int
	queue         chan packages.Pkg
	results       chan Fetcher
}

type ClientConfig struct {
	MaxWorkers             int
	Cache                  cache.Cache
	CachePkgExpirationTime *time.Duration
}

// New returns a new client
func New(config ClientConfig) *Client {
	if config.MaxWorkers == 0 {
		config.MaxWorkers = 10
	}
	if config.Cache != nil && config.CachePkgExpirationTime == nil {
		hour := time.Hour
		config.CachePkgExpirationTime = &hour
	}
	return &Client{
		ClientConfig: config,
		Packages:     make(map[string]packages.Pkg),
		Queued:       &sync.Map{},
		Pool:         worker.NewPool(config.MaxWorkers),
	}
}

// GetDependencyTree returns a dependency tree for a given package
func (c *Client) GetDependencyTree(pkg packages.Pkg) error {
	c.queue = make(chan packages.Pkg)
	defer close(c.queue)

	c.results = make(chan Fetcher)
	defer close(c.results)

	c.Pool.Start()

	c.addJob(pkg)
	err := c.waitForResults()
	if err != nil {
		return err
	}

	pkg.BuildTree(c.Packages)
	return nil
}

func (c *Client) waitForResults() error {
	for {
		select {
		case pkg := <-c.queue:
			c.addJob(pkg)
		case f := <-c.results:
			if f.Err != nil && c.JobsCompleted == 0 {
				return f.Err
			}
			c.Packages[f.Pkg.GetCacheKey()] = f.Pkg

			c.JobsCompleted++
			if c.JobsCreated == c.JobsCompleted {
				c.Pool.Close()
				return nil
			}
		}
	}
}

func (c *Client) addJob(pkg packages.Pkg) {
	f := Fetcher{
		Pkg:       pkg,
		Results:   c.results,
		Queue:     c.queue,
		Queued:    c.Queued,
		Cache:     c.ClientConfig.Cache,
		CacheTime: c.ClientConfig.CachePkgExpirationTime,
	}
	c.Pool.AddJob(f)
	c.JobsCreated++
}
