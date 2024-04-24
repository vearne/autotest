package rule

import "github.com/go-resty/resty/v2"

type VerifyRule interface {
	Name() string
	Valid(response *resty.Response) bool
}
