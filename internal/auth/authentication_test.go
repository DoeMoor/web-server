package auth

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	type args struct {
		tokenString string
		tokenSecret string
	}
	tests := []struct {
		name    string
		args    args
		want    uuid.UUID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateJWT(tt.args.tokenString, tt.args.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeJWT(t *testing.T) {
	type args struct {
		userID              uuid.UUID
		tokenSecret         string
		secretForValidation string
		expiriesIn          time.Duration
	}
	tests := []struct {
		name      string
		args      args
		want      string
		wantErr   bool
		errorText string
	}{
		{
			name: "full positive test",
			args: args{
				userID:              uuid.New(),
				tokenSecret:         "printenv",
				secretForValidation: "printenv",
				expiriesIn:          time.Minute * 5,
			},
			wantErr: false,
		},
		{
			name: "wrong secret",
			args: args{
				userID:              uuid.New(),
				tokenSecret:         "printenv",
				secretForValidation: "printenv2",
				expiriesIn:          time.Minute * 5,
			},
			wantErr: true,
		},
		{
			name: "expired token",
			args: args{
				userID:              uuid.New(),
				tokenSecret:         "printenv",
				secretForValidation: "printenv",
				expiriesIn:          time.Nanosecond * 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeJWT(tt.args.userID, tt.args.tokenSecret, tt.args.expiriesIn)
			if err != nil {
				t.Errorf("MakeJWT() error = %v", err)
				return
			}
			result, err := ValidateJWT(got, tt.args.secretForValidation)

			if err != nil {
				if tt.wantErr {
					return
				}
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result.String() != tt.args.userID.String() {
				t.Errorf("MakeJWT() = %v, want %v", got, tt.args.userID.String())
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	type args struct {
		header http.Header
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "full positive test",
			args: args{
				header: http.Header{"Authorization": []string{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ"}},
			},
			wantErr: false,
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ",
		},
		{
			name: "wrong token",
			args: args{
				header: http.Header{"Authorization": []string{"BearereyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.ey"}},
			},
			wantErr: true,
		},
		{
			name: "no token",
			args: args{
				header: http.Header{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := GetBearerToken(tt.args.header)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
