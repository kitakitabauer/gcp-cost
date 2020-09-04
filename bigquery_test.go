package main

import (
	"os"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/kylelemons/godebug/pretty"
)

func init() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "../key/gcp-cost-xxxx.json")
}

func TestBigQueryImpl_CreateQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			"ok",
			"2020/04/10",
			"\nSELECT" + Selector + "\nFROM billing.gcp_billing_export_v1_00CAC5_AADE6E_BC76F8\nWHERE\n\tformat_timestamp('%Y/%m/%d', usage_start_time, 'Asia/Tokyo') = \"2020/04/10\"\n\tand project.name = 'liverpool'\nGROUP BY\n\tday, service\nORDER BY\n\tcost desc\n",
		},
	}

	b := &BigQueryImpl{}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			out, err := b.CreateQuery(v.in)

			if err != nil {
				t.Errorf("%v", err)
			}
			if out != v.out {
				t.Errorf("input %v\n, get: %#v\n, want: %#v\n", v.in, out, v.out)
			}
		})
	}
}

func TestConvertStruct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []bigquery.Value
		out  CostResult
	}{
		{
			"ok",
			[]bigquery.Value{
				"2020/04/10",
				"Big Query",
				123.06508,
			},
			CostResult{
				Service: "Big Query",
				Cost:    123.06508,
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			res := convertStruct(v.in)
			if diff := pretty.Compare(res, v.out); diff != "" {
				t.Errorf("convertStruct diff: (-got +want)\n%s", diff)
			}
		})
	}
}
