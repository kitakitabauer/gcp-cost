package main

import (
	"bytes"
	"context"
	"strconv"
	"text/template"

	"cloud.google.com/go/bigquery"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
)

const (
	ProjectCost = "gcp-cost"
	Selector    = `
	format_timestamp('%Y/%m/%d', usage_start_time, 'Asia/Tokyo') as day,
	service.description as service,
	round(
		(sum(cost) +
			sum(IFNULL((SELECT SUM(c.amount) FROM UNNEST(credits) c), 0))
		) * 100) / 100 as cost,
`
	QueryTemplate = `
SELECT{{.Selector}}
FROM billing.gcp_billing_export_v1_00CAC5_AADE6E_BC76F8
WHERE
	format_timestamp('%Y/%m/%d', usage_start_time, 'Asia/Tokyo') = "{{.TargetYmd}}"
	and project.name = 'liverpool'
GROUP BY
	day, service
ORDER BY
	cost desc
`
)

type (
	QueryVars struct {
		Selector  string
		TargetYmd string
		ProjectID string
	}

	CostResult struct {
		Service string
		Cost    float64
	}
)

type bigQuery interface {
	CreateQuery(targetYmd string) (string, error)
	SelectTable(query string) ([]CostResult, error)
}

// BigQueryImpl implementes BigQuery interface.
type BigQueryImpl struct {
	log *zap.Logger
}

func (b *BigQueryImpl) CreateQuery(targetYmd string) (string, error) {
	vars := QueryVars{
		Selector:  Selector,
		TargetYmd: targetYmd,
	}

	var msg bytes.Buffer
	err := template.Must(template.New("msg").Parse(QueryTemplate)).Execute(&msg, vars)
	if err != nil {
		err = errors.Wrap(err, "failed to parse")
		return msg.String(), err
	}

	return msg.String(), nil
}

func (b *BigQueryImpl) SelectTable(query string) ([]CostResult, error) {
	ctx := context.Background()
	cli, err := bigquery.NewClient(ctx, ProjectCost)
	if err != nil {
		b.log.Error("bigquery.new.client.error", zap.Error(err))
		return nil, err
	}
	defer cli.Close()

	q := cli.Query(query)
	iter, err := q.Read(ctx)
	if err != nil {
		b.log.Error("bigquery.read.error", zap.Error(err))
		return nil, err
	}

	var res []CostResult
	for {
		var values []bigquery.Value

		err := iter.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			b.log.Error("bigquery.iterate.error", zap.Error(err))
			continue
		}

		b.log.Info("SelectTable", zap.Any("value", values))

		res = append(res, convertStruct(values))
	}

	return res, nil
}

func convertStruct(values []bigquery.Value) CostResult {
	var res CostResult
	for i := range values {
		switch i {
		case 0:
			continue
		case 1:
			res.Service = values[i].(string)
		case 2:
			res.Cost = values[i].(float64)
		default:
			panic("Unexpected index: " + strconv.Itoa(i))
		}
	}
	return res
}