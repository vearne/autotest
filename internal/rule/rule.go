package rule

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/vearne/autotest/internal/model"
)

type VerifyRule interface {
	Name() string
	Verify(response *resty.Response) bool
}

type VerifyRuleGrpc interface {
	Name() string
	Verify(response *model.GrpcResp) bool
}

func convStr(v any) string {
	return fmt.Sprintf("%v", v)
}
