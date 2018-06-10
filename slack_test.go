package main

import (
	"reflect"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	nSlack "github.com/nlopes/slack"
)

func TestDivideGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []CostResult
		out  map[string][]CostResult
	}{
		{
			"empty",
			[]CostResult{},
			map[string][]CostResult{},
		},
		{
			"ok",
			[]CostResult{
				{
					ResourceType: "hoge",
				},
				{
					ResourceType: "1",
				},
				{
					ResourceType: "hoge",
				},
			},
			map[string][]CostResult{
				"hoge": {
					{
						ResourceType: "hoge",
					},
					{
						ResourceType: "hoge",
					},
				},
				"1": {
					{
						ResourceType: "1",
					},
				},
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			out := divideType(v.in)
			if diff := pretty.Compare(out, v.out); diff != "" {
				t.Errorf("convertStruct diff: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestCalcGroupCost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		inService string
		inCosts   []CostResult
		out       Type
	}{
		{
			"all display costs",
			"group1",
			[]CostResult{
				{
					ResourceType: "hoge",
					Cost:         0.011,
					Currency:     "USD",
				},
				{
					ResourceType: "fuga",
					Cost:         0.0101,
					Currency:     "USD",
				},
				{
					ResourceType: "hage",
					Cost:         0.01001,
					Currency:     "USD",
				},
			},
			Type{
				Cost: 0.03111,
				Attachment: nSlack.Attachment{
					Title: "[group1]",
					Fields: []nSlack.AttachmentField{
						{
							Title: "hoge",
							Value: "0.011 (USD)",
							Short: true,
						},
						{
							Title: "fuga",
							Value: "0.010 (USD)",
							Short: true,
						},
						{
							Title: "hage",
							Value: "0.010 (USD)",
							Short: true,
						},
					},
				},
			},
		},
		{
			"some costs limit display",
			"group2",
			[]CostResult{
				{
					ResourceType: "hoge",
					Cost:         0.011,
					Currency:     "USD",
				},
				{
					ResourceType: "fuga",
					Cost:         0.01,
					Currency:     "USD",
				},
				{
					ResourceType: "hage",
					Cost:         0.0101,
					Currency:     "USD",
				},
				{
					ResourceType: "fugu",
					Cost:         0.009,
					Currency:     "USD",
				},
			},
			Type{
				Cost: 0.0211,
				Attachment: nSlack.Attachment{
					Title: "[group2]",
					Fields: []nSlack.AttachmentField{
						{
							Title: "hoge",
							Value: "0.011 (USD)",
							Short: true,
						},
						{
							Title: "hage",
							Value: "0.010 (USD)",
							Short: true,
						},
					},
				},
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			out := calcTypeCost(v.inService, v.inCosts)
			if !reflect.DeepEqual(out, v.out) {
				t.Errorf("get: %v\n, want: %v\n", out, v.out)
			}
		})
	}
}

func TestSelectColorIndex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   int
		out  int
	}{
		{
			"zero",
			0,
			1,
		},
		{
			"one",
			1,
			2,
		},
		{
			"two",
			2,
			3,
		},
		{
			"three",
			3,
			0,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			out := selectColorIndex(v.in)
			if out != v.out {
				t.Errorf("input %v\n, get: %#v\n, want: %#v\n", v.in, out, v.out)
			}
		})
	}
}
