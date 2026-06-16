package email

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"

	"github.com/luyb177/life-tracker/backend/common/mail"
)

// EmailSender 邮件发送接口
type EmailSender interface {
	// SendVerifyCode 发送验证码邮件
	SendVerifyCode(ctx context.Context, to, code string, expireMinutes int) error
}

type DefaultEmailSender struct {
	mailer *mail.Mailer
}

func NewEmailSender(m *mail.Mailer) EmailSender {
	return &DefaultEmailSender{
		mailer: m,
	}
}

func (s *DefaultEmailSender) SendVerifyCode(ctx context.Context, to, code string, expireMinutes int) error {
	subject := "【Life Tracker】邮箱验证码"
	data := map[string]interface{}{
		"Code":          code,
		"ExpireMinutes": expireMinutes,
	}

	return s.renderAndSend(ctx, subject, verifyCodeTmpl, data, []string{to})
}

func (s *DefaultEmailSender) renderAndSend(ctx context.Context, subject string, tmpl *template.Template, data interface{}, to []string) error {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute email template failed: %w", err)
	}
	return s.mailer.Send(subject, buf.String(), to)
}
