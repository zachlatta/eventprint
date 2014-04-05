package eventbrite

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/zachlatta/eventprint/server/helper"

	"code.google.com/p/goauth2/oauth"
)

var transport *oauth.Transport

func init() {
	transport = &oauth.Transport{}
	transport.Token = &oauth.Token{AccessToken: helper.GetConfig("ACCESS_TOKEN")}
}

func GetAttendees() ([]Attendee, error) {
	return getAttendees([]Attendee{}, 1)
}

func getAttendees(attendees []Attendee, pageNumber int) ([]Attendee, error) {
	resp, err := getAttendeesPage(pageNumber)
	if err != nil {
		return nil, err
	}

	attendees = append(attendees, resp.Attendees...)

	log.Println(pageNumber)
	if resp.Pagination.PageNumber < resp.Pagination.PageCount {
		return getAttendees(attendees, pageNumber+1)
	}

	return attendees, nil
}

func getAttendeesPage(pageNumber int) (*EventbriteResponse, error) {
	// TODO: Use config file
	r, err := transport.Client().Get(
		fmt.Sprintf("https://www.eventbriteapi.com/v3/events/%s/attendees/?page=%d", helper.GetConfig("EVENT_ID"), pageNumber))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var resp *EventbriteResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Fatal(err)
	}

	return resp, nil
}
