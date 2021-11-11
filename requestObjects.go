package mbsGo

type StatementRequestObject struct {
	AccountNo     string `json:"accountNo"`
	BankId        int    `json:"bankId"`
	DestinationId string `json:"destinationId"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	Role          string `json:"role"`
	Username      string `json:"username"`
	Country       string `json:"country"`
	Phone         string `json:"phone"`
	Applicants    []struct {
		Name          string `json:"name"`
		ApplicationNo string `json:"applicationNo"`
	} `json:"applicants"`
}

//TODO
//	Add Validators

type ConfirmStatementRequest struct {
	TicketNo string `json:"ticketNo"`
	Password string `json:"password"`
}

type RequestId struct {
	RequestId int `json:"requestId"`
}
