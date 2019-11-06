package api

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/m1/go-dependency-check/api/response"
	"github.com/m1/go-dependency-check/client"
	"github.com/m1/go-dependency-check/client/packages"
)

type PackagesHandler struct {
	*API
}

func NewPackagesHandler(api *API) *PackagesHandler {
	return &PackagesHandler{api}
}

func (h *PackagesHandler) GetRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/npm/{package}", h.GetNpmPackage)
	router.Get("/npm/{package}/{version}", h.GetNpmPackageVersion)
	return router
}

func (h *PackagesHandler) GetNpmPackage(w http.ResponseWriter, r *http.Request) {
	h.handleNpmPackage(w, r)
}

func (h *PackagesHandler) GetNpmPackageVersion(w http.ResponseWriter, r *http.Request) {
	h.handleNpmPackage(w, r)
}

func (h *PackagesHandler) handleNpmPackage(w http.ResponseWriter, r *http.Request) {
	c := client.New(client.ClientConfig{
		MaxWorkers: 10,
		Cache:      h.cache,
	})

	pkg := packages.NewNpm(
		chi.URLParam(r, "package"),
		chi.URLParam(r, "version"),
	)
	err := c.GetDependencyTree(pkg)
	if err != nil {
		response.RespondError(w, r, http.StatusNotFound, err)
		return
	}

	response.RespondOk(w, r, &response.Data{Data: pkg})
}
