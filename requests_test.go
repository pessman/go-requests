package requests

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pessman/go-requests/client"
	"github.com/pessman/go-requests/client/mock_client"
	"github.com/pessman/go-requests/response"
)

func TestNewBytesResponse(t *testing.T) {
	type args struct {
		b  []byte
		sc int
	}
	tests := []struct {
		name string
		args args
		want *response.BytesResponse
	}{
		{
			name: "get a new bytes response object",
			args: args{
				b:  []byte("test"),
				sc: 100,
			},
			want: &response.BytesResponse{
				Body:       []byte("test"),
				StatusCode: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := response.NewBytesResponse(tt.args.b, tt.args.sc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBytesResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_request_Response(t *testing.T) {
	type fields struct {
		URL       string
		Method    string
		Body      io.Reader
		Client    client.Client
		Modifiers []Modifier
	}
	type clientArgs struct {
		wantClient bool
		callCount  int
		f          func(*http.Request) (*http.Response, error)
	}
	tests := []struct {
		name       string
		fields     fields
		clientArgs clientArgs
		want       *http.Response
		wantErr    bool
	}{
		{
			name: "new http request error",
			fields: fields{
				URL: "://google.com",
			},
			clientArgs: clientArgs{
				wantClient: true,
				callCount:  0,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "use http default client",
			fields: fields{
				URL: "://google.com",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "return info from client Do func",
			clientArgs: clientArgs{
				wantClient: true,
				callCount:  1,
				f: func(r *http.Request) (*http.Response, error) {
					return nil, errors.New("random error")
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mockClient *mock_client.MockClient
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.clientArgs.wantClient {
				mockClient = mock_client.NewMockClient(ctrl)
				mockClient.EXPECT().Do(gomock.Any()).Times(tt.clientArgs.callCount)
				tt.fields.Client = mockClient
			}
			r := &request{
				URL:       tt.fields.URL,
				Method:    tt.fields.Method,
				Body:      tt.fields.Body,
				Client:    tt.fields.Client,
				Modifiers: tt.fields.Modifiers,
			}
			got, err := r.Response()
			if (err != nil) != tt.wantErr {
				t.Errorf("request.Response() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("request.Response() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_request_newHTTPRequest(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "https://google.com", nil)
	type fields struct {
		URL       string
		Method    string
		Body      io.Reader
		Client    client.Client
		Modifiers []Modifier
	}
	tests := []struct {
		name    string
		fields  fields
		want    *http.Request
		wantErr bool
	}{
		{
			name: "new http request error",
			fields: fields{
				URL: "://google.com",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "http request returned",
			fields: fields{
				URL:    "https://google.com",
				Method: http.MethodGet,
				Body:   nil,
				Client: nil,
				Modifiers: []Modifier{
					func(req *http.Request) {},
				},
			},
			want: r,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &request{
				URL:       tt.fields.URL,
				Method:    tt.fields.Method,
				Body:      tt.fields.Body,
				Client:    tt.fields.Client,
				Modifiers: tt.fields.Modifiers,
			}
			got, err := r.newHTTPRequest()
			if (err != nil) != tt.wantErr {
				t.Errorf("request.newHTTPRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("request.newHTTPRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRequest(t *testing.T) {
	type args struct {
		url    string
		method string
		body   io.Reader
		client client.Client
	}
	tests := []struct {
		name string
		args args
		want Request
	}{
		{
			name: "create new request",
			args: args{
				url:    "https://google.com",
				method: http.MethodGet,
				body:   nil,
				client: nil,
			},
			want: &request{
				URL:    "https://google.com",
				Method: http.MethodGet,
				Body:   nil,
				Client: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRequest(tt.args.url, tt.args.method, tt.args.body, tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_request_BytesResponse(t *testing.T) {
	type fields struct {
		URL       string
		Method    string
		Body      io.Reader
		Client    client.Client
		Modifiers []Modifier
	}
	type clientArgs struct {
		f func(*http.Request) (*http.Response, error)
	}
	tests := []struct {
		name       string
		fields     fields
		clientArgs clientArgs
		want       *response.BytesResponse
		wantErr    bool
	}{
		{
			name: "response returns error",
			fields: fields{
				URL:    "https://google.com",
				Method: http.MethodGet,
			},
			clientArgs: clientArgs{
				f: func(*http.Request) (*http.Response, error) {
					return nil, errors.New("response error")
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "return BytesResponse",
			fields: fields{
				URL:    "https://google.com",
				Method: http.MethodGet,
			},
			clientArgs: clientArgs{
				f: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						Body:       ioutil.NopCloser(bytes.NewBuffer([]byte("{}"))),
						StatusCode: 201,
					}, nil
				},
			},
			want:    response.NewBytesResponse([]byte("{}"), 201),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mock_client.NewMockClient(ctrl)
			mockClient.EXPECT().Do(gomock.Any()).DoAndReturn(tt.clientArgs.f)
			r := &request{
				URL:       tt.fields.URL,
				Method:    tt.fields.Method,
				Body:      tt.fields.Body,
				Client:    mockClient,
				Modifiers: tt.fields.Modifiers,
			}
			got, err := r.BytesResponse()
			if (err != nil) != tt.wantErr {
				t.Errorf("request.BytesResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("request.BytesResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
