package client

import (
	"context"
	"github.com/quanxiang-cloud/process/pkg/code"
	"github.com/quanxiang-cloud/process/pkg/config"
	"github.com/quanxiang-cloud/process/pkg/misc/client"
	"github.com/quanxiang-cloud/process/pkg/misc/error2"
	"net/http"
	"strings"
)

// Condition service
type Condition interface {
	GetConditionResult(ctx context.Context, conditionStr string, params map[string]interface{}) (bool, error)
}

// NewCondition new
func NewCondition(conf *config.Configs) (Condition, error) {
	c := &condition{
		getConditionResult: conf.APIHost.FlowHost + "/api/v1/flow/formula/calculation",
		client:             client.New(conf.InternalNet),
	}
	return c, nil
}

type condition struct {
	getConditionResult string
	client             http.Client
}

// GetConditionResult get
func (c *condition) GetConditionResult(ctx context.Context, conditionStr string, params map[string]interface{}) (bool, error) {
	var resp map[string]bool

	req := map[string]interface{}{
		"expression": strings.Replace(conditionStr, "\"", "'", -1),
		"parameter":  params,
	}

	err := client.POST(ctx, &c.client, c.getConditionResult, req, &resp)
	if err != nil {
		return false, error2.NewError(code.ConditionParamError)
	}
	if r, ok := resp["result"]; ok {
		return r, nil
	}
	return false, error2.NewError(code.NoResult)
}
