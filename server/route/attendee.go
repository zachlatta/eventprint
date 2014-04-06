package route

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/zachlatta/eventprint/server/model"
	"github.com/zachlatta/eventprint/server/websockets"
)

// GetAttendees fetches all of the attendees in the database.
func GetAttendees(db gorp.SqlExecutor, w http.ResponseWriter,
	params martini.Params, log *log.Logger) (int, string) {

	var attendees []model.Attendee
	_, err := db.Select(&attendees, "SELECT * FROM Attendee ORDER BY Id")
	if err != nil {
		log.Println("Error retrieving attendees:", err)
		return http.StatusInternalServerError, "Error retrieving attendees"
	}

	json, err := json.Marshal(attendees)
	if err != nil {
		log.Println("Error marshalling attendees to JSON:", err)
		return http.StatusInternalServerError, "Error retrieving attendees"
	}

	w.Header().Set("Content-Type", "application/json")
	return http.StatusOK, string(json)
}

// GetAttendee fetches an attendee by id and returns it serialized as JSON.
func GetAttendee(db gorp.SqlExecutor, w http.ResponseWriter,
	params martini.Params, log *log.Logger) (int, string) {

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Printf("Error converting id %d to string: %v\n", params["id"], err)
		return http.StatusBadRequest, "Id not an int"
	}

	obj, err := db.Get(model.Attendee{}, id)
	if err != nil {
		log.Printf("Error retrieving attendee with id %d: %v\n", id, err)
		return http.StatusNotFound, "Error retrieving attendee"
	}

	attendee := obj.(*model.Attendee)

	json, err := json.Marshal(attendee)
	if err != nil {
		log.Println("Error marshalling attendee to JSON:", err)
		return http.StatusInternalServerError, "Error when retrieving attendee"
	}

	w.Header().Set("Content-Type", "application/json")
	return http.StatusOK, string(json)
}

// CheckInAttendee sets the CheckedIn field of an attendee to true. It returns
// the updated attendee serialized as JSON.
func CheckInAttendee(db gorp.SqlExecutor, w http.ResponseWriter,
	params martini.Params, log *log.Logger) (int, string) {

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Printf("Error converting id %d to string: %v\n", params["id"], err)
		return http.StatusBadRequest, "Id not an int"
	}

	obj, err := db.Get(model.Attendee{}, id)
	if err != nil {
		log.Printf("Error retrieving attendee with id %d: %v\n", id, err)
		return http.StatusNotFound, ""
	}

	attendee := obj.(*model.Attendee)

	attendee.CheckedIn = true

	if _, err := db.Update(attendee); err != nil {
		log.Println("Error updating attendee:", err)
		return http.StatusInternalServerError, "Error signing in attendee."
	}

	json, err := json.Marshal(attendee)
	if err != nil {
		log.Println("Error marshalling attendee to JSON:", err)
		return http.StatusInternalServerError, "Error when retrieving attendee"
	}

	websockets.Hub.Broadcast(string(json))

	w.Header().Set("Content-Type", "application/json")
	return http.StatusOK, string(json)
}
