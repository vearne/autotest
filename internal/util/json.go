package util

import (
	"bytes"
	"encoding/json"
	"strings"
)

// PrettyJSON 如果 s 是合法 JSON，则返回缩进美化后的内容；否则原样返回。
// 使用 json.Indent 而非反序列化后重新序列化，可以保留字段的原始顺序。
func PrettyJSON(s string) string {
	if strings.TrimSpace(s) == "" {
		return s
	}

	var out bytes.Buffer
	if err := json.Indent(&out, []byte(s), "", "  "); err != nil {
		return s
	}
	return out.String()
}
