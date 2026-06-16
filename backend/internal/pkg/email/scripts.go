package email

import (
	_ "embed"
	"html/template"
)

//go:embed templates/verify_code.html
var verifyCodeTemplate string

var verifyCodeTmpl = template.Must(template.New("verify_code").Parse(verifyCodeTemplate))
