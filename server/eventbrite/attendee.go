package eventbrite

type EventbriteResponse struct {
	Pagination Pagination `json:"pagination"`
	Attendees  []Attendee `json:"attendees"`
}

type Pagination struct {
	ObjectCount int `json:"object_count"`
	PageNumber  int `json:"page_number"`
	PageSize    int `json:"page_size"`
	PageCount   int `json:"page_count"`
}

type Attendee struct {
	Quantity int       `json:"quantity"`
	TicketId int       `json:"ticket_id,string"`
	Profile  Profile   `json:"profile"`
	Barcodes []Barcode `json:"barcodes"`
}

type Profile struct {
	Gender    string `json:"gender"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       string `json:"age"`
	Email     string `json:"email"`
}

type Barcode struct {
	Status  string `json:"status"`
	Barcode string `json:"barcode"`
}
