package model

type Attendee struct {
	Id           int64  `json:"id"`
	EventbriteId int64  `json:"eventbriteId"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Title        string `json:"title"`
	Gender       string `json:"gender"`
	Email        string `json:"email"`
	CheckedIn    bool   `json:"checkedIn"`
	Barcode      string `json:"barcode"`
}
