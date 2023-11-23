package msgraphgocore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_convertToPage(t *testing.T) {
	type args struct {
		response interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    PageResult[validTestStruct]
		wantErr bool
	}{
		{
			name: "should pass",
			args: args{
				response: &validTestStruct{
					obj: []validTestStruct{
						{
							obj: []validTestStruct{
								{
									obj: []validTestStruct{},
								},
							},
						},
					}},
			},
			want: PageResult[validTestStruct]{
				oDataNextLink:  nil,
				oDataDeltaLink: nil,
				value: []validTestStruct{
					{
						obj: []validTestStruct{
							{
								obj: []validTestStruct{},
							},
						},
					},
				},
			},
		},
		{
			name: "should return error 'saying response cannot be nil' for nil response",
			args: args{
				response: nil,
			},
			want: PageResult[validTestStruct]{
				oDataNextLink:  nil,
				oDataDeltaLink: nil,
				value:          nil,
			},
			wantErr: true,
		},
		{
			name: "should return error 'value property missing in response object' for missing 'GetValue' method",
			args: args{
				response: &invalidTestStruct{
					obj: []invalidTestStruct{
						{
							obj: []invalidTestStruct{
								{
									obj: []invalidTestStruct{},
								},
							},
						},
					}},
			},
			want: PageResult[validTestStruct]{
				oDataNextLink:  nil,
				oDataDeltaLink: nil,
			},
			wantErr: true,
		},
		{
			name: "should return error 'response does not have next link accessor' for missing 'GetOdataNextLink() *string' method",
			args: args{
				response: &invalidTestStruct{
					obj: []invalidTestStruct{
						{
							obj: []invalidTestStruct{
								{
									obj: []invalidTestStruct{},
								},
							},
						},
					}},
			},
			want: PageResult[validTestStruct]{
				oDataNextLink:  nil,
				oDataDeltaLink: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToPage[validTestStruct](tt.args.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want.oDataDeltaLink, got.oDataDeltaLink, "got %v, want %v", got.oDataNextLink, tt.want.oDataNextLink)
			assert.Equal(t, tt.want.oDataNextLink, got.oDataNextLink, "got %v, want %v", got.oDataDeltaLink, tt.want.oDataDeltaLink)
			assert.Equal(t, tt.want.value, got.value, "got %v, want %v", got.value, tt.want.value)
		})
	}
}

type validTestStruct struct {
	obj []validTestStruct
}

func (t *validTestStruct) GetValue() []validTestStruct {
	return t.obj
}

func (t *validTestStruct) GetOdataNextLink() *string {
	return nil
}

func (t *validTestStruct) GetOdataDeltaLink() *string {
	return nil
}

type invalidTestStruct struct {
	obj []invalidTestStruct
}

func (t *invalidTestStruct) GetOdataDeltaLink() *string {
	return nil
}
