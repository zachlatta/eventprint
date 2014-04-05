package model

type Attendee struct {
	Id           int64  `json:"id"`
	EventbriteId int64  `json:"eventbriteId"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Gender       string `json:"gender"`
	Age          string `json:"age"`
	Email        string `json:"email"`
	CheckedIn    bool   `json:"checkedIn"`
	Barcode      string `json:"barcode"`
}
