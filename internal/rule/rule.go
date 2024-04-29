package rule

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

type VerifyRule interface {
	Name() string
	Verify(response *resty.Response) bool
}

func convStr(v any) string {
	return fmt.Sprintf("%v", v)
}
