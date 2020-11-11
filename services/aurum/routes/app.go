package routes// Access determines if a user is allowed access to a certain application
import (
	"encoding/json"
	"github.com/finitum/aurum/pkg/aurum"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

// GET /access/{app}/{user}
func (rs Routes) Access(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")
	app := chi.URLParam(r, "app")

	appid, err := uuid.Parse(app)
	if err != nil {
		_ = RenderError(w, err, InvalidRequest)
		return
	} else if appid == uuid.Nil {
		_ = RenderError(w, errors.New("appid zero"), InvalidRequest)
		return
	}

	resp, err := aurum.Access(r.Context(), rs.store, user, appid)
	if err != nil {
		_ = RenderError(w, err, ServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(&resp)
}

