package mbsGo

import (
	"errors"
	"fmt"
	"github.com/0sax/err2"
	"net/http"
)

const (
	BaseUrl                        = "https://mybankstatement.net/TP/api"
	RequestStatementEndpoint       = "/RequestStatement"
	GetFeedbackByRequestIdEndpoint = "/GetFeedbackByRequestID"
	ConfirmStatementEndpoint       = "/ConfirmStatement"
	ReConfirmStatementEndpoint     = "/ReconfirmStatement"
	ListBanksEndpoint              = "/SelectActiveRequestBanks"
	GetStatementJSONEndpoint       = "/GetStatementJSONObject"
	GetStatementPDFEndpoint        = "/GetPDFStatement"

	Applicant = "Applicant"
	Sponsor   = "Sponsor"
	Guarantor = "Guarantor"
)

func statementRequestRoles() []string {
	return []string{Applicant, Sponsor, Guarantor}
}

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

	err = req.Validate()
	if err != nil {
		err = err2.NewClientErr(err, err.Error(), 400)
		return
	}

	req.DestinationId = cl.clientId

	var resp Response

	err = cl.standardRequest(http.MethodPost, RequestStatementEndpoint, req, &resp)
	if err != nil {
		return
	}

	if resp.isSuccessful() {
		requestId, err = resp.requestId()
		return
	}

	ee, err := resp.errors()
	if err != nil {
		return
	}

	m := fmt.Sprintf("%v, %v", resp.Message, ee)

	err = err2.NewClientErr(errors.New(m), m, 400)
	return
}

func (cl *Client) GetFeedbackByRequestID(reqID int) (status, feedback string, err error) {

	reqIDWrap := &RequestId{
		reqID,
	}

	var resp Response

	err = cl.standardRequest(http.MethodPost, GetFeedbackByRequestIdEndpoint, reqIDWrap, &resp)
	if err != nil {
		return
	}

	if resp.isSuccessful() {
		fb, ee := resp.feedBack()
		if ee != nil {
			err = ee
			return
		}
		status = fb.Status
		feedback = fb.Feedback
		return
	}

	ee, err := resp.errors()
	if err != nil {
		return
	}

	m := fmt.Sprintf("%v, %v", resp.Message, ee)

	err = err2.NewClientErr(errors.New(m), m, 400)
	//fmt.Printf("hoohoh : %v", m)
	return
}

func (cl *Client) ConfirmStatement(ticketNumber, password string) (message string, err error) {

	req := &ConfirmStatementRequest{
		TicketNo: ticketNumber,
		Password: password,
	}
	var resp Response

	err = cl.standardRequest(http.MethodPost, ConfirmStatementEndpoint, req, &resp)
	if err != nil {
		return
	}

	if resp.isSuccessful() {
		message = resp.Message
		return
	}

	ee, err := resp.errors()
	if err != nil {
		return
	}
	m := fmt.Sprintf("%v, %v", resp.Message, ee)

	err = err2.NewClientErr(errors.New(m), m, 400)
	return
}

func (cl *Client) ReConfirmStatement(reqID int) (message string, err error) {

	reqIDWrap := &RequestId{
		reqID,
	}

	var resp Response

	err = cl.standardRequest(http.MethodPost, ReConfirmStatementEndpoint, reqIDWrap, &resp)
	if err != nil {
		return
	}

	if resp.isSuccessful() {
		message = resp.Message
		return
	}
	m := fmt.Sprintf("%v", resp.Message)

	err = err2.NewClientErr(errors.New(m), m, 400)
	return
}

func (cl *Client) GetBankList() (list []Bank, err error) {

	var resp Response

	err = cl.standardRequest(http.MethodPost, ListBanksEndpoint, nil, &resp)
	if err != nil {
		return
	}

	if resp.isSuccessful() {
		list, err = resp.bankList()
		return
	}
	m := fmt.Sprintf("%v", resp.Message)
	err = err2.NewClientErr(errors.New(m), m, 400)
	return
}

func (cl *Client) GetStatementJSON(ticketNumber, password string) (st *JSONStatement, err error) {

	req := &ConfirmStatementRequest{
		TicketNo: ticketNumber,
		Password: password,
	}
	var resp Response

	err = cl.standardRequest(http.MethodPost, GetStatementJSONEndpoint, req, &resp)
	if err != nil {
		return
	}

	if resp.isSuccessful() {
		st, err = resp.jsonStatementObject()
		return
	}
	ee, err := resp.errors()
	if err != nil {
		return
	}

	m := fmt.Sprintf("%v, %v", resp.Message, ee)

	err = err2.NewClientErr(errors.New(m), m, 400)
	return
}

func (cl *Client) GetStatementPDF(ticketNumber, password string) (base64PDF string, err error) {

	req := &ConfirmStatementRequest{
		TicketNo: ticketNumber,
		Password: password,
	}
	var resp Response

	err = cl.standardRequest(http.MethodPost, GetStatementPDFEndpoint, req, &resp)
	if err != nil {
		return
	}

	if resp.isSuccessful() {
		base64PDF, err = resp.pdfStatementString()
		return
	}
	ee, err := resp.errors()
	if err != nil {
		return
	}

	m := fmt.Sprintf("%v, %v", resp.Message, ee)

	err = err2.NewClientErr(errors.New(m), m, 400)
	return
}
