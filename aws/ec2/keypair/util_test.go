package keypair

import (
	"testing"
)

func TestStringValue(t *testing.T) {
	type args struct {
		properties   map[string]interface{}
		key          string
		defaultValue string
	}
	tests := []struct {
		name      string
		args      args
		wantValue string
	}{
		{
			"zero-from-zero",
			args{
				map[string]interface{}{},
				"missing",
				"",
			},
			"",
		},
		{
			"default-from-zero",
			args{
				map[string]interface{}{},
				"missing",
				"default",
			},
			"default",
		},
		{
			"default-from-missing",
			args{
				map[string]interface{}{
					"key": "value",
				},
				"missing",
				"default",
			},
			"default",
		},
		{
			"zero-from-missing",
			args{
				map[string]interface{}{
					"key": "value",
				},
				"missing",
				"default",
			},
			"default",
		},
		{
			"match-from-match",
			args{
				map[string]interface{}{
					"match": "value",
				},
				"match",
				"default",
			},
			"value",
		},
		{
			"zero-from-match-zero",
			args{
				map[string]interface{}{
					"match": "",
				},
				"match",
				"default",
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotValue := StringValue(tt.args.properties, tt.args.key, tt.args.defaultValue); gotValue != tt.wantValue {
				t.Errorf("StringValue() = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func TestBoolValue(t *testing.T) {
	type args struct {
		properties   map[string]interface{}
		key          string
		defaultValue bool
	}
	tests := []struct {
		name      string
		args      args
		wantValue bool
	}{
		{
			"false-from-zero",
			args{
				map[string]interface{}{},
				"missing",
				false,
			},
			false,
		},
		{
			"true-from-zero",
			args{
				map[string]interface{}{},
				"missing",
				true,
			},
			true,
		},
		{
			"true-from-missing",
			args{
				map[string]interface{}{
					"key": "value",
				},
				"missing",
				true,
			},
			true,
		},
		{
			"false-from-missing",
			args{
				map[string]interface{}{
					"key": "value",
				},
				"missing",
				false,
			},
			false,
		},
		{
			"true-from-match-nil",
			args{
				map[string]interface{}{
					"match": nil,
				},
				"match",
				true,
			},
			true,
		},
		{
			"false-from-match-nil",
			args{
				map[string]interface{}{
					"match": nil,
				},
				"match",
				false,
			},
			false,
		},
		{
			"false-from-match-asdf",
			args{
				map[string]interface{}{
					"match": "asdf",
				},
				"match",
				true,
			},
			false,
		},
		{
			"match-from-match-bool",
			args{
				map[string]interface{}{
					"match": true,
				},
				"match",
				false,
			},
			true,
		},
		{
			"match-from-match-string-0",
			args{
				map[string]interface{}{
					"match": "0",
				},
				"match",
				true,
			},
			false,
		},
		{
			"match-from-match-string-1",
			args{
				map[string]interface{}{
					"match": "1",
				},
				"match",
				false,
			},
			true,
		},
		{
			"match-from-match-string-false",
			args{
				map[string]interface{}{
					"match": "false",
				},
				"match",
				true,
			},
			false,
		},
		{
			"match-from-match-string-true",
			args{
				map[string]interface{}{
					"match": "true",
				},
				"match",
				false,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotValue := BoolValue(tt.args.properties, tt.args.key, tt.args.defaultValue); gotValue != tt.wantValue {
				t.Errorf("BoolValue() = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}
