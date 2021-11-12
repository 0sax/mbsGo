package mbsGo

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"os"
	"reflect"
	"testing"
)

var cc *Client
var ee error

func TestMain(m *testing.M) {
	// Write code here to run before tests
	ee = godotenv.Load("vars.env")
	if ee != nil {
		log.Fatalf("authentication error: %v", ee)
	}

	cc = NewClient(BaseUrl, os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))

	// Run tests
	exitVal := m.Run()

	// Write code here to run after tests

	// Exit with exit value from tests
	os.Exit(exitVal)
}

func TestClient_RequestStatement(t *testing.T) {

	type args struct {
		req *StatementRequestObject
	}
	tests := []struct {
		name    string
		cl      *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test 1",
			cc,
			args{
				&StatementRequestObject{
					AccountNo:     "0136537576",
					BankId:        1,
					DestinationId: cc.clientId,
					StartDate:     "01-Jul-2021",
					EndDate:       "10-Nov-2021",
					Role:          "Applicant",
					Username:      os.Getenv("USERNAME"),
					Country:       "NG",
					Phone:         "0115542254",
					Applicants: []struct {
						Name          string `json:"name"`
						ApplicationNo string `json:"applicationNo"`
					}{{
						"Kendrick Lamar",
						""},
					},
				}},
			false,
		},
		{"bad request",
			cc,
			args{
				&StatementRequestObject{
					AccountNo:     "",
					BankId:        0,
					DestinationId: "",
					StartDate:     "",
					EndDate:       "",
					Role:          "",
					Username:      "",
					Country:       "",
					Phone:         "",
					Applicants:    nil,
				}},
			true,
		},
		{"bad client",
			&Client{
				clientId:     "",
				clientSecret: "",
				baseUrl:      BaseUrl,
			},
			args{
				&StatementRequestObject{
					AccountNo:     "",
					BankId:        0,
					DestinationId: "",
					StartDate:     "",
					EndDate:       "",
					Role:          "",
					Username:      "",
					Country:       "",
					Phone:         "",
					Applicants:    nil,
				}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := tt.cl.RequestStatement(tt.args.req)
			if err != nil {
				t.Logf("\n\nerror: %v", err)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("\n\nRequestStatement() error = %v, wantErr %v", err, tt.wantErr)
			}
			//if gotRequestId != tt.wantRequestId {
			//	t.Errorf("RequestStatement() gotRequestId = %v, want %v", gotRequestId, tt.wantRequestId)
			//}
		})
	}
}

func TestClient_GetFeedbackByRequestID(t *testing.T) {

	type args struct {
		reqID int
	}
	tests := []struct {
		name       string
		cl         *Client
		args       args
		wantStatus string
		wantErr    bool
	}{
		{"Not Found",
			cc,
			args{212},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus, _, err := tt.cl.GetFeedbackByRequestID(tt.args.reqID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFeedbackByRequestID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("GetFeedbackByRequestID() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestClient_ConfirmStatement(t *testing.T) {
	tests := []struct {
		name         string
		cl           *Client
		ticketNumber string
		password     string
		wantMessage  string
		wantErr      bool
	}{
		{"case 1",
			cc,
			"4359975-13",
			"8678",
			"Successfully confirmed request",
			false,
		},
		{"case 2",
			cc,
			os.Getenv("VALID_TICKET_NUMBER"),
			os.Getenv("VALID_PASSWORD"),
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotMessage, err := tt.cl.ConfirmStatement(tt.ticketNumber, tt.password)
			if err != nil {
				t.Logf("errmsg %v\n", err)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfirmStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMessage != tt.wantMessage {
				t.Errorf("ConfirmStatement() gotMessage = %v, want %v", gotMessage, tt.wantMessage)
			}
		})
	}
}

func TestClient_ReConfirmStatement(t *testing.T) {
	tests := []struct {
		name        string
		cl          *Client
		reqID       int
		wantMessage string
		wantErr     bool
	}{
		{"case 1",
			cc,
			2222222222234,
			"",
			true},
		{"case 2",
			cc,
			11111,
			"Request successfully confirmed",
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotMessage, err := tt.cl.ReConfirmStatement(tt.reqID)
			if err != nil {
				t.Logf("errmsg %v\n", err)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfirmStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMessage != tt.wantMessage {
				t.Errorf("ConfirmStatement() gotMessage = %v, want %v", gotMessage, tt.wantMessage)
			}
		})
	}
}

func TestClient_GetBankList(t *testing.T) {

	banksAsString := ` [
        {
            "id": 6,
            "name": "Access Bank",
            "sortCode": "044"
        },
        {
            "id": 32,
            "name": "Eco Bank",
            "sortCode": "050"
        },
        {
            "id": 5,
            "name": "FCMB",
            "sortCode": "214"
        },
        {
            "id": 15,
            "name": "Fidelity Bank",
            "sortCode": "070"
        },
        {
            "id": 3,
            "name": "First Bank",
            "sortCode": "011"
        },
        {
            "id": 13,
            "name": "GT Bank",
            "sortCode": "058"
        },
        {
            "id": 7,
            "name": "Heritage Bank",
            "sortCode": "030"
        },
        {
            "id": 4,
            "name": "Keystone Bank",
            "sortCode": "082"
        },
        {
            "id": 2,
            "name": "Polaris Bank Limited",
            "sortCode": "076"
        },
        {
            "id": 37,
            "name": "Providus Bank",
            "sortCode": "101"
        },
        {
            "id": 10,
            "name": "Stanbic IBTC Bank",
            "sortCode": "221"
        },
        {
            "id": 1,
            "name": "Sterling Bank",
            "sortCode": "232"
        },
        {
            "id": 14,
            "name": "UBA ",
            "sortCode": "033"
        },
        {
            "id": 11,
            "name": "Union Bank",
            "sortCode": "032"
        },
        {
            "id": 9,
            "name": "Unity Bank",
            "sortCode": "215"
        },
        {
            "id": 12,
            "name": "Wema Bank",
            "sortCode": "035"
        },
        {
            "id": 17,
            "name": "Zenith Bank",
            "sortCode": "057"
        }
    ]`
	var bankList []Bank
	json.Unmarshal([]byte(banksAsString), &bankList)
	tests := []struct {
		name     string
		cl       *Client
		wantList []Bank
		wantErr  bool
	}{
		{"case1",
			cc,
			bankList,
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotList, err := tt.cl.GetBankList()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBankList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotList, tt.wantList) {
				t.Errorf("GetBankList() gotList = %v, want %v", gotList, tt.wantList)
			}
		})
	}
}

func TestClient_GetStatementPDF(t *testing.T) {

	tests := []struct {
		name         string
		cl           *Client
		ticketNumber string
		password     string
		wantErr      bool
	}{
		{"case 1",
			cc,
			os.Getenv("VALID_TICKET_NUMBER"),
			os.Getenv("VALID_PASSWORD"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotSt, err := tt.cl.GetStatementPDF(tt.ticketNumber, tt.password)
			if err != nil {
				t.Logf("error message: %v\n", err)
			}
			t.Logf("base64 String: %v\n", gotSt)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStatementPDF() error = %v, wantErr %v", err, tt.wantErr)
			}
			//if gotSt != tt.wantSt {
			//	t.Errorf("GetStatementPDF() gotSt = %v, want %v", gotSt, tt.wantSt)
			//}
		})
	}
}
