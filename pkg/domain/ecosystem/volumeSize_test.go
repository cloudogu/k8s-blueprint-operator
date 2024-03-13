package ecosystem

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

func TestGetQuantityReference(t *testing.T) {
	twoGigaByte := resource.MustParse("2G")
	type args struct {
		quantityStr string
	}
	tests := []struct {
		name    string
		args    args
		want    *resource.Quantity
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "should parse quantity and return a reference",
			args:    args{quantityStr: "2G"},
			want:    &twoGigaByte,
			wantErr: assert.NoError,
		},
		{
			name:    "should return error on invalid quantity string",
			args:    args{quantityStr: "2GG"},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetQuantityReference(tt.args.quantityStr)
			if !tt.wantErr(t, err, fmt.Sprintf("GetQuantityReference(%v)", tt.args.quantityStr)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetQuantityReference(%v)", tt.args.quantityStr)
		})
	}
}

func TestGetQuantityString(t *testing.T) {
	twoGigaByte := resource.MustParse("2G")
	type args struct {
		quantity *resource.Quantity
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should return string if reference is not nil",
			args: args{quantity: &twoGigaByte},
			want: "2G",
		},
		{
			name: "should return empty string if reference is nil",
			args: args{quantity: nil},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetQuantityString(tt.args.quantity), "GetQuantityString(%v)", tt.args.quantity)
		})
	}
}
