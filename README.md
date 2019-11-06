# go-dependency-check

To run: `docker-compose up --build`

Can run as server or can just import `client` to use as a library.

To query: `localhost:8181/packages/npm/{package}/{version}`

Can just leave `{version}` blank - will default to `latest`

Example:
```
➜  curl localhost:8181/packages/npm/es6-symbol/0.1.1 | jq
{
  "status": 200,
  "status_desc": "OK",
  "data": {
    "npm_package": {
      "data": {
        "name": "es6-symbol",
        "version": "0.1.1",
        "dependencies": [
          {
            "name": "d",
            "version": "0.1.1",
            "dependencies": [
              {
                "name": "es5-ext",
                "version": "0.10.2",
                "dependencies": null
              }
            ]
          },
          {
            "name": "es5-ext",
            "version": "0.10.4",
            "dependencies": [
              {
                "name": "es6-iterator",
                "version": "0.1.1",
                "dependencies": [
                  {
                    "name": "es5-ext",
                    "version": "0.10.2",
                    "dependencies": null
                  },
                  {
                    "name": "es6-symbol",
                    "version": "0.1.0",
                    "dependencies": [
                      {
                        "name": "d",
                        "version": "0.1.1",
                        "dependencies": [
                          {
                            "name": "es5-ext",
                            "version": "0.10.2",
                            "dependencies": null
                          }
                        ]
                      },
                      {
                        "name": "es5-ext",
                        "version": "0.10.2",
                        "dependencies": null
                      }
                    ]
                  },
                  {
                    "name": "d",
                    "version": "0.1.1",
                    "dependencies": [
                      {
                        "name": "es5-ext",
                        "version": "0.10.2",
                        "dependencies": null
                      }
                    ]
                  }
                ]
              },
              {
                "name": "es6-symbol",
                "version": "0.1.0",
                "dependencies": [
                  {
                    "name": "d",
                    "version": "0.1.1",
                    "dependencies": [
                      {
                        "name": "es5-ext",
                        "version": "0.10.2",
                        "dependencies": null
                      }
                    ]
                  },
                  {
                    "name": "es5-ext",
                    "version": "0.10.2",
                    "dependencies": null
                  }
                ]
              }
            ]
          }
        ]
      }
    }
  }
}
```

## Assumptions made/Future plans

- The packages are interfaces so adding future package managers will be easy.
- Testing needs expanding (only around 60% coverage atm) - will be looked at expanding in future.
- Will be moved out of docker-compose.
- Different caching methods will be added due to cache being an interface.
