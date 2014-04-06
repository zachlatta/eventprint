package route

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/zachlatta/eventprint/server/eventbrite"
	"github.com/zachlatta/eventprint/server/model"
)

// Sync pulls the latest attendees from Eventbrite and stores them in the
// database. Returns an array of new attendees.
func Sync(db gorp.SqlExecutor, w http.ResponseWriter) string {
	eventbriteAttendees, err := eventbrite.GetAttendees()
	if err != nil {
		log.Fatal("Error getting Eventbrite attendees:", err)
	}

	var attendees []model.Attendee
	for _, evtAttendee := range eventbriteAttendees {
		checkedIn := false

		if evtAttendee.Barcodes[0].Status == "used" {
			checkedIn = true
		}

		attendee := model.Attendee{
			EventbriteId: evtAttendee.Id,
			FirstName:    evtAttendee.Profile.FirstName,
			LastName:     evtAttendee.Profile.LastName,
			Gender:       evtAttendee.Profile.Gender,
			Email:        evtAttendee.Profile.Email,
			CheckedIn:    checkedIn,
			Barcode:      evtAttendee.Barcodes[0].Barcode,
		}

		// Error will be thrown if attendee already exists in DB.
		if err := db.Insert(&attendee); err == nil {
			// Only add new attendees to attendees slice
			attendees = append(attendees, attendee)
		}
	}

	json, err := json.Marshal(attendees)
	if err != nil {
		log.Fatal(err)
	}

	// If the attendees array is empty, the json will be null.
	if string(json) == "null" {
		json = []byte("[]")
	}

	w.Header().Set("Content-Type", "application/json")
	return string(json)
}
