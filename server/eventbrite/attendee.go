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
	Id       int64     `json:"id,string"`
	Quantity int       `json:"quantity"`
	Profile  Profile   `json:"profile"`
	Barcodes []Barcode `json:"barcodes"`
}

type Profile struct {
	Gender    string `json:"gender"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       string `json:"age,omitempty"` // Eventbrite doesn't always return ints
	Email     string `json:"email"`
}

type Barcode struct {
	Status  string `json:"status"`
	Barcode string `json:"barcode"`
}
