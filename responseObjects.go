package mbsGo

type RequestStatementResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"` // int or string array
}

func (r *RequestStatementResponse) isSuccessful() bool {
	return r.Status == "00"
}

func (r *RequestStatementResponse) resultIsString() bool {
	_, ok := r.Result.(string)
	return ok
}

func (r *RequestStatementResponse) getResultString() string {
	if r.resultIsString() {
		return r.Result.(string)
	}
	return ""
}

func (r *RequestStatementResponse) getResultId() int {
	if r.isSuccessful() {
		return r.Result.(int)
	}
	return 0
}

func (r *RequestStatementResponse) getErrors() (s []string) {
	r1 := r.Result.([]interface{})

	for _, e := range r1 {
		s = append(s, e.(string))
	}

	return
}

type RequestIDFeedbackResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  struct {
		Status   string `json:"status"`
		Feedback string `json:"feedback"`
	} `json:"result"`
}

func (r *RequestIDFeedbackResponse) isSuccessful() bool {
	return r.Status == "00"
}

type BankListResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []Bank `json:"result"`
}

func (r *BankListResponse) isSuccessful() bool {
	return r.Status == "00"
}

type Bank struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	SortCode string `json:"sortCode"`
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