package mbsGo

import (
	"encoding/json"
	"fmt"
	"github.com/0sax/err2"
	"net/http"
	"strings"
)

const (
	BaseUrl                        = "https://mybankstatement.net/TP/api"
	RequestStatementEndpoint       = "/RequestStatement"
	GetFeedbackByRequestIdEndpoint = "/GetFeedbackByRequestID"
	ConfirmStatementEndpoint       = "/ConfirmStatement"
	ReConfirmStatementEndpoint     = "/ReconfirmStatement"
	ListBanksEndpoint              = "/SelectActiveRequestBanks"
	GetStatementJSONEndpoint       = "/GetStatementJSONObject"
)

type Client struct {
	clientId     string
	clientSecret string
	baseUrl      string
}

func NewClient(baseUrl, clientId, clientSecret string) *Client {
	return &Client{
		clientId:     clientId,
		clientSecret: clientSecret,
		baseUrl:      baseUrl,
	}
}

func NewStatementRequestObject(bankId int, accountNo, destinationId, startDate, endDate, role, username, country, phone, applicantName string) *StatementRequestObject {
	return &StatementRequestObject{
		AccountNo:     accountNo,
		BankId:        bankId,
		DestinationId: destinationId,
		StartDate:     startDate,
		EndDate:       endDate,
		Role:          role,
		Username:      username,
		Country:       country,
		Phone:         phone,
		Applicants: []struct {
			Name          string `json:"name"`
			ApplicationNo string `json:"applicationNo"`
		}{{applicantName, ""}},
	}
}

func (cl *Client) RequestStatement(req *StatementRequestObject) (requestId int, err error) {

	req.DestinationId = cl.clientId

	var resp RequestStatementResponse

	err = cl.standardRequest(http.MethodPost, RequestStatementEndpoint, req, &resp)
	if err != nil {
		return
	}

	if resp.isSuccessful() {
		return resp.getResultId(), nil
	}

	err = err2.NewClientErr(nil, fmt.Sprint(resp.getErrors()), 400)
	return
}

func (cl *Client) GetFeedbackByRequestID(reqID int) (status, feedback string, err error) {

	reqIDWrap := &RequestId{
		reqID,
	}

	var resp RequestIDFeedbackResponse

	err = cl.standardRequest(http.MethodPost, GetFeedbackByRequestIdEndpoint, reqIDWrap, resp)
	if err != nil {
		return
	}

	if resp.isSuccessful() {
		status = resp.Result.Status
		feedback = resp.Result.Feedback
		return
	}

	status = resp.Message
	feedback = resp.Message
	err = err2.NewClientErr(nil, resp.Message, 400)
	return
}

func (cl *Client) ConfirmStatement(ticketNumber, password string) (message string, err error) {

	req := &ConfirmStatementRequest{
		TicketNo: ticketNumber,
		Password: password,
	}
	var resp RequestStatementResponse

	err = cl.standardRequest(http.MethodPost, ConfirmStatementEndpoint, req, &resp)
	if err != nil {
		return
	}

	if resp.isSuccessful() {
		message = resp.Message
	}

	err = err2.NewClientErr(nil, resp.Message, 400)
	return
}

func (cl *Client) ReConfirmStatement(reqID int) (message string, err error) {

	reqIDWrap := &RequestId{
		reqID,
	}

	var resp RequestStatementResponse

	err = cl.standardRequest(http.MethodPost, ReConfirmStatementEndpoint, reqIDWrap, &resp)
	if err != nil {
		return
	}

	if resp.isSuccessful() {
		message = resp.Message
	}

	err = err2.NewClientErr(nil, resp.Message, 400)
	return
}

func (cl *Client) GetBankList() (list []Bank, err error) {

	var resp BankListResponse

	err = cl.standardRequest(http.MethodPost, ListBanksEndpoint, nil, &resp)
	if err != nil {
		return
	}

	if resp.isSuccessful() {
		list = resp.Result
		return
	}

	err = err2.NewClientErr(nil, resp.Message, 400)
	return
}

func (cl *Client) GetStatementJSON(ticketNumber, password string) (st JSONStatement, err error) {

	req := &ConfirmStatementRequest{
		TicketNo: ticketNumber,
		Password: password,
	}
	var resp RequestStatementResponse

	err = cl.standardRequest(http.MethodPost, GetStatementJSONEndpoint, req, &resp)
	if err != nil {
		return
	}

	if resp.isSuccessful() {
		if resp.resultIsString() {
			// todo unmarshal the string into a statement object
			statementString := strings.Trim(resp.getResultString(), `\`)
			err = json.Unmarshal([]byte(statementString), &st)
			return
		}
	}

	err = err2.NewClientErr(nil, resp.Message, 400)
	return
}
