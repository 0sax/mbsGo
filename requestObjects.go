package mbsGo

import (
	"github.com/0sax/err2"
	"github.com/thoas/go-funk"
)

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

func (sro *StatementRequestObject) Validate() (err error) {
	var errs err2.ErrList
	//TODO
	//	Add Validators

	if !funk.ContainsString(statementRequestRoles(), sro.Role) {
		errs.AddF("role '%v' is not a valid role", sro.Role)
	}

	if errs.HasErrs() {
		err = errs.ErrsAsError()
	}
	return
}

type ConfirmStatementRequest struct {
	TicketNo string `json:"ticketNo"`
	Password string `json:"password"`
}

type RequestId struct {
	RequestId int `json:"requestId"`
}
