package main

import (
	"bytes"
	"context"
	"log"
	"strconv"
	"text/template"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

const (
	ProjectID = "" // your project ID
	TableName = "" // your table name

	Selector = `
  resource_type,
  start_time,
  end_time,
  cost,
  usage.unit`

	QueryTemplate = `
SELECT {{.Selector}}
FROM {{.TableName}}.{{.TargetYmd}}
WHERE
  project.id = "{{.ProjectID}}" AND
  cost > 0
ORDER BY
  cost desc
`
)

type (
	QueryVars struct {
		Selector  string
		TargetYmd string
		TableName string
		ProjectID string
	}

	CostResult struct {
		ResourceType string
		Start        string
		End          string
		Cost         float64
		Currency     string
	}
)

type BigQuery interface {
	CreateQuery(targetYmd string) (string, error)
	SelectTable(query string) ([]CostResult, error)
}

// BigQueryImpl implementes BigQuery interface.
type BigQueryImpl struct{}

func (b *BigQueryImpl) CreateQuery(targetYmd string) (string, error) {
	vars := QueryVars{
		Selector:  Selector,
		TargetYmd: targetYmd,
		TableName: TableName,
		ProjectID: ProjectID,
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
	cli, err := bigquery.NewClient(ctx, ProjectID)
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	q := cli.Query(query)
	iter, err := q.Read(ctx)
	if err != nil {
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
			log.Printf("Failed to iterate: %v", err)
			continue
		}

		res = append(res, convertStruct(values))
	}

	return res, nil
}

func convertStruct(values []bigquery.Value) CostResult {
	var res CostResult
	for i := range values {
		switch i {
		case 0:
			res.ResourceType = values[i].(string)
			break
		case 1:
			t := values[i].(time.Time)
			res.Start = t.String()
			break
		case 2:
			t := values[i].(time.Time)
			res.End = t.String()
			break
		case 3:
			res.Cost = values[i].(float64)
			break
		case 4:
			res.Currency = values[i].(string)
			break
		default:
			panic("Unexpected index: " + strconv.Itoa(i))
		}
	}
	return res
}
