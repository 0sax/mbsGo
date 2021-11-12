package mbsGo

import (
	"encoding/json"
	"fmt"
	"strings"
)

// All the unnecessary looking methods for the response objects
// are here because someone at Wallz & Queen has never
// used a statically typed language. Anyway, I'm still alive sha
// glory be to the motherboard.
//
// Please, I don't know who needs to hear this, but a field in an
// api response should NEVER return more than one data type. If you
// have been doing this, please say no to cultism today.

type Response struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

func (r *Response) isSuccessful() bool {
	return r.Status == "00"
}
func (r *Response) errors() (errs []string, err error) {
	if r.Result != nil {
		err = json.Unmarshal([]byte(r.Result), &errs)
		if err != nil {
			err = fmt.Errorf("response.errors() unmarshal failed because: %v", err.Error())
		}
	}
	return
}
func (r *Response) requestId() (rId int, err error) {
	err = json.Unmarshal([]byte(r.Result), &rId)
	return
}
func (r *Response) bankList() (bl []Bank, err error) {
	err = json.Unmarshal([]byte(r.Result), &bl)
	return
}
func (r *Response) feedBack() (fb *Feedback, err error) {
	err = json.Unmarshal([]byte(r.Result), &fb)
	return
}
func (r *Response) jsonStatementObject() (js *JSONStatement, err error) {
	var jsStr string
	err = json.Unmarshal([]byte(r.Result), &jsStr)
	if err != nil {
		return
	}

	statementString := strings.Trim(jsStr, `\`)
	err = json.Unmarshal([]byte(statementString), &js)
	return
}
func (r *Response) pdfStatementString() (ps string, err error) {
	err = json.Unmarshal([]byte(r.Result), &ps)
	if err != nil {
		return
	}
	return
}

type Bank struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	SortCode string `json:"sortCode"`
}

type Feedback struct {
	Status   string `json:"status"`
	Feedback string `json:"feedback"`
}

type JSONStatement struct {
	Status          string `json:"status"`
	Name            string `json:"Name"`
	Nuban           string `json:"Nuban"`
	AccountCategory string `json:"AccountCategory"`
	AccountType     string `json:"AccountType"`
	TicketNo        string `json:"TicketNo"`
	AvailableBal    string `json:"AvailableBal"`
	BookBal         string `json:"BookBal"`
	TotalCredit     string `json:"TotalCredit"`
	TotalDebit      string `json:"TotalDebit"`
	Tenor           string `json:"Tenor"`
	Period          string `json:"Period"`
	Currency        string `json:"Currency"`
	Address         string `json:"Address"`
	Applicants      string `json:"Applicants"`
	Signatories     []struct {
		Name string `json:"Name"`
		BVN  string `json:"BVN"`
	} `json:"Signatories"`
	Details []struct {
		PTransactionDate string `json:"PTransactionDate"`
		PValueDate       string `json:"PValueDate"`
		PNarration       string `json:"PNarration"`
		PCredit          string `json:"PCredit"`
		PDebit           string `json:"PDebit"`
		PBalance         string `json:"PBalance"`
	} `json:"Details"`
}
