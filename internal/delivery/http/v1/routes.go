package v1

import (
	"encoding/json"
	"net/http"

	"github.com/AlekSi/pointer"
	"github.com/gorilla/mux"

	"goals_scheduler/internal/cns"
	"goals_scheduler/internal/models"
	"goals_scheduler/internal/usecase"
)

func RegisterHandlers(uc *usecase.UseCase) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/callback", makeCallbackHandler(uc))
	return router
}

func makeCallbackHandler(uc *usecase.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type callbackData struct {
			ID int64 `json:"ID"`
		}

		var cb callbackData
		err := json.NewDecoder(r.Body).Decode(&cb)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = uc.NotifierUpdate(r.Context(), models.NotifierCU{
			Status: pointer.To(cns.StatusNotifierEnded),
		}, cb.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
