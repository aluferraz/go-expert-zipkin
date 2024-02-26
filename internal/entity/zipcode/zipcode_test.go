package zipcode

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ZipcodeTestSuite struct {
	suite.Suite
	ctx context.Context
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ZipcodeTestSuite))
}

func (suite *ZipcodeTestSuite) SetupSuite() {
	suite.ctx = context.Background()
}
func (suite *ZipcodeTestSuite) TestNewZipcode() {
	tests := []struct {
		name    string
		zipcode string
		wantErr bool
	}{
		{"Valid zipcode", "12345678", false},
		{"Less than 8 digits", "1234567", true},
		{"More than 8 digits", "1234567890", true},
		{"Contains non-digit characters", "123A5678", true},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			_, err := NewZipcode(tt.zipcode)
			if (err != nil) != tt.wantErr {
				suite.T().Errorf("NewZipcode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func (suite *ZipcodeTestSuite) TestIsValid() {
	tests := []struct {
		name    string
		zipcode Zipcode
		wantErr bool
	}{
		{"Valid zipcode", Zipcode{"12345678"}, false},
		{"Less than 8 digits", Zipcode{"1234567"}, true},
		{"More than 8 digits", Zipcode{"1234567890"}, true},
		{"Contains non-digit characters", Zipcode{"123A5678"}, true},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := tt.zipcode.IsValid()
			if (err != nil) != tt.wantErr {
				suite.T().Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
