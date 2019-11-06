package response

import (
	"net/http"

	"github.com/go-chi/render"
)

type Data struct {
	Data Model `json:"data"`
}

type Response struct {
	Status     *int             `json:"status"`
	StatusDesc *string          `json:"status_desc"`
	Error      string           `json:"error,omitempty"`
	Data       map[string]*Data `json:"data,omitempty"`
}

func Respond(w http.ResponseWriter, r *http.Request, status int, datum ...*Data) {
	var dataMap = make(map[string]*Data)
	for _, d := range datum {
		dataMap[d.Data.GetJSONKey()] = d
	}

	text := http.StatusText(status)
	res := Response{
		Status:     &status,
		StatusDesc: &text,
		Data:       dataMap,
	}

	render.Status(r, status)
	render.SetContentType(render.ContentTypeJSON)
	render.JSON(w, r, res)
	return
}

func RespondOk(w http.ResponseWriter, r *http.Request, datum ...*Data) {
	Respond(w, r, http.StatusOK, datum...)
	return
}

func RespondError(w http.ResponseWriter, r *http.Request, status int, err error) {
	text := http.StatusText(status)
	res := Response{
		Status:     &status,
		StatusDesc: &text,
		Error:      err.Error(),
	}

	render.Status(r, status)
	render.SetContentType(render.ContentTypeJSON)
	render.JSON(w, r, res)
	return
}
