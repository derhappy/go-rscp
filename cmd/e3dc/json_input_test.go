package main

import (
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/spali/go-rscp/rscp"
)

func Test_unmarshalJSONRequest(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    rscp.Message
		wantErr bool
	}{
		{`tuple with data type and value`,
			`["RSCP_AUTHENTICATION_USER","CString","testuser"]`,
			rscp.Message{Tag: rscp.RSCP_AUTHENTICATION_USER, DataType: rscp.CString, Value: "testuser"},
			false,
		},
		{"empty results in error",
			"",
			rscp.Message{},
			true,
		},
		{`string`,
			`"INFO_REQ_UTC_TIME"`,
			rscp.Message{Tag: rscp.INFO_REQ_UTC_TIME},
			false,
		},
		{`tuple`,
			`["INFO_REQ_UTC_TIME"]`,
			rscp.Message{Tag: rscp.INFO_REQ_UTC_TIME},
			false,
		},
		{`tuple with data type`,
			`["INFO_REQ_UTC_TIME","None"]`,
			rscp.Message{Tag: rscp.INFO_REQ_UTC_TIME},
			false,
		},
		{`tuple with value`,
			`["RSCP_AUTHENTICATION_USER","testuser"]`,
			rscp.Message{Tag: rscp.RSCP_AUTHENTICATION_USER, DataType: rscp.CString, Value: "testuser"},
			false,
		},
		{`tuple with data type and value`,
			`["RSCP_AUTHENTICATION_USER","CString","testuser"]`,
			rscp.Message{Tag: rscp.RSCP_AUTHENTICATION_USER, DataType: rscp.CString, Value: "testuser"},
			false,
		},
		{`tuple with nested data`,
			`["BAT_REQ_DATA",[["BAT_INDEX",0],"BAT_REQ_DEVICE_STATE"]]`,
			rscp.Message{Tag: rscp.BAT_REQ_DATA, DataType: rscp.Container, Value: []rscp.Message{{Tag: rscp.BAT_INDEX, DataType: rscp.UInt16, Value: uint16(0)}, {Tag: rscp.BAT_REQ_DEVICE_STATE}}},
			false,
		},
		{`tuple (invalid empty)`,
			`[]`,
			rscp.Message{},
			true,
		},
		{`tuple (invalid)`,
			`[x]`,
			rscp.Message{},
			true,
		},
		{`tuple (invalid tag)`,
			`[1]`,
			rscp.Message{},
			true,
		},
		{`tuple (invalid tag and datatype)`,
			`[1,1]`,
			rscp.Message{},
			true,
		},
		{`tuple (invalid tag, datatype and value)`,
			`[1,1,1]`,
			rscp.Message{},
			true,
		},
		{`tuple (invalid data type) keeps untouched`,
			`["INFO_REQ_UTC_TIME", "UChar8"]`,
			rscp.Message{Tag: rscp.INFO_REQ_UTC_TIME, DataType: rscp.UChar8},
			false,
		},
		{`tuple (invalid value) get's fixed`,
			`["INFO_REQ_MAC_ADDRESS","None",""]`,
			rscp.Message{Tag: rscp.INFO_REQ_MAC_ADDRESS},
			false,
		},
		{`object`,
			`{ "Tag": "INFO_REQ_UTC_TIME" }`,
			rscp.Message{Tag: rscp.INFO_REQ_UTC_TIME},
			false,
		},
		{`object with value`,
			`{ "Tag": "BAT_INDEX", "Value": 0 }`,
			rscp.Message{Tag: rscp.BAT_INDEX, DataType: rscp.UInt16, Value: uint16(0)},
			false,
		},
		{`object all fields`,
			`{ "Tag": "BAT_INDEX", "DataType": "UInt16", "Value": 0 }`,
			rscp.Message{Tag: rscp.BAT_INDEX, DataType: rscp.UInt16, Value: uint16(0)},
			false,
		},
		{`object with nested data`,
			`{ "Tag": "BAT_REQ_DATA", "Value": [ { "Tag": "BAT_INDEX", "Value": 0 }, { "Tag": "BAT_REQ_DEVICE_STATE" } ] }`,
			rscp.Message{Tag: rscp.BAT_REQ_DATA, DataType: rscp.Container, Value: []rscp.Message{{Tag: rscp.BAT_INDEX, DataType: rscp.UInt16, Value: uint16(0)}, {Tag: rscp.BAT_REQ_DEVICE_STATE}}},
			false,
		},
		{`object (invalid tag)`,
			`{ "Tag": 1 }`,
			rscp.Message{},
			true,
		},
		{`string (invalid tag)`,
			`"INVALID_TAG"`,
			rscp.Message{},
			true,
		},
		{`time value`,
			`["INFO_SET_TIME","1234-05-06T07:08:09.123456Z"]`,
			rscp.Message{Tag: rscp.INFO_SET_TIME, DataType: rscp.Timestamp, Value: time.Date(1234, 5, 6, 7, 8, 9, 123456000, time.UTC)},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := rscp.Message{}
			err := unmarshalJSONRequest([]byte(tt.message), &m)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalJSONRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(m, tt.want); diff != nil {
				t.Errorf("unmarshalJSONRequest() = %v, want %v\n%s", m, tt.want, diff)
			}
		})
	}
}

func Test_unmarshalJSONRequests(t *testing.T) {

	tests := []struct {
		name    string
		message string
		want    []rscp.Message
		wantErr bool
	}{
		{"empty results in error",
			"",
			nil,
			true,
		},
		{`array of tags`,
			`["INFO_REQ_MAC_ADDRESS","INFO_REQ_UTC_TIME"]`,
			[]rscp.Message{{Tag: rscp.INFO_REQ_MAC_ADDRESS}, {Tag: rscp.INFO_REQ_UTC_TIME}},
			false,
		},
		{`array of tuples`,
			`[["INFO_REQ_MAC_ADDRESS"],["INFO_REQ_UTC_TIME"]]`,
			[]rscp.Message{{Tag: rscp.INFO_REQ_MAC_ADDRESS}, {Tag: rscp.INFO_REQ_UTC_TIME}},
			false,
		},
		{`array of messages`,
			`[{ "Tag": "INFO_REQ_MAC_ADDRESS" }, { "Tag": "INFO_REQ_UTC_TIME" }]`,
			[]rscp.Message{{Tag: rscp.INFO_REQ_MAC_ADDRESS}, {Tag: rscp.INFO_REQ_UTC_TIME}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unmarshalJSONRequests([]byte(tt.message))
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalJSONRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("unmarshalJSONRequests() = %v, want %v\n%s", got, tt.want, diff)
			}
		})
	}
}
