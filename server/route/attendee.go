package route

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/zachlatta/eventprint/server/model"
)

// GetAttendee fetches an attendee by id and returns it serialized as JSON.
func GetAttendee(db gorp.SqlExecutor, w http.ResponseWriter,
	params martini.Params, log *log.Logger) (int, string) {

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return http.StatusBadRequest, "Id not an int"
	}

	obj, err := db.Get(model.Attendee{}, id)
	if err != nil {
		return http.StatusNotFound, ""
	}

	attendee := obj.(*model.Attendee)

	json, err := json.Marshal(attendee)
	if err != nil {
		return http.StatusInternalServerError, "Error when retrieving attendee"
	}

	w.Header().Set("Content-Type", "application/json")
	return http.StatusOK, string(json)
}

func SignInAttendee(db gorp.SqlExecutor) {
}
