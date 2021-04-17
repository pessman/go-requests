package response

import (
	"reflect"
	"testing"
)

func TestNewBytesResponse(t *testing.T) {
	type args struct {
		b  []byte
		sc int
	}
	tests := []struct {
		name string
		args args
		want *BytesResponse
	}{
		{
			name: "new response",
			args: args{
				b:  []byte(""),
				sc: 200,
			},
			want: &BytesResponse{
				Body:       []byte(""),
				StatusCode: 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBytesResponse(tt.args.b, tt.args.sc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBytesResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
