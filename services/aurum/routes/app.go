package routes // Access determines if a user is allowed access to a certain application
import (
	"encoding/json"
	"github.com/finitum/aurum/pkg/aurum"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"net/http"
)

// GET /access/{app}/{user}
func (rs Routes) Access(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")
	name := chi.URLParam(r, "app")

	if name == "" || user == "" {
		_ = RenderError(w, errors.New("name empty"), InvalidRequest)
		return
	}


	resp, err := aurum.Access(r.Context(), rs.store, user, name)
	if err != nil {
		_ = RenderError(w, err, ServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(&resp)
}
