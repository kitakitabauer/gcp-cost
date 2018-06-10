package main

import (
	"os"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/kylelemons/godebug/pretty"
)

func init() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "../key/credentials.json")
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
			"20180601",
			"\nSELECT " + Selector + "\nFROM " + TableName + ".20180601\nWHERE\n  ProjectID = \"" + ProjectID + "\" AND\n  Cost > 0\nORDER BY\n  Cost desc\n",
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

func TestBigQueryImpl_SelectTable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   BigQuery
		out  error
	}{
		{
			"ok",
			&BigQueryImpl{},
			nil,
		},
	}

	targetYmd := time.Now().AddDate(0, 0, -2).Format("20060102")
	query := "SELECT " + Selector + " FROM " + TableName + "." + targetYmd + " WHERE project.id = \"" + ProjectID + "\" AND cost > 0 ORDER BY cost desc"

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			_, err := v.in.SelectTable(query)

			if err != v.out {
				t.Errorf("input %v\n, get: %#v\n, want: %#v\n", v.in, err, v.out)
			}
		})
	}
}

func TestConvertStruct(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name string
		in   []bigquery.Value
		out  CostResult
	}{
		{
			"ok",
			[]bigquery.Value{
				"hoge",
				now,
				now,
				123.06508,
				"USD",
			},
			CostResult{
				ResourceType: "hoge",
				Start:        now.String(),
				End:          now.String(),
				Cost:         123.06508,
				Currency:     "USD",
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
