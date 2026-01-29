package service

import (
	"fmt"
	"log"
	"net/smtp"

	"linuxdo-review/config"
)

// EmailService 邮件服务
type EmailService struct {
	host     string
	port     int
	user     string
	password string
	from     string
	enabled  bool
}

// NewEmailService 创建邮件服务
func NewEmailService(cfg *config.Config) *EmailService {
	enabled := cfg.SMTP.Host != "" && cfg.SMTP.User != ""
	return &EmailService{
		host:     cfg.SMTP.Host,
		port:     cfg.SMTP.Port,
		user:     cfg.SMTP.User,
		password: cfg.SMTP.Password,
		from:     cfg.SMTP.From,
		enabled:  enabled,
	}
}

// IsEnabled 检查邮件服务是否启用
func (s *EmailService) IsEnabled() bool {
	return s.enabled
}

// SendInviteCode 发送邀请码邮件
func (s *EmailService) SendInviteCode(to, username, inviteCode string) error {
	if !s.enabled {
		log.Printf("[EmailService] SMTP未配置,跳过发送邀请码邮件给 %s", to)
		return nil
	}

	subject := "🎉 恭喜！您的Linux.do邀请码申请已通过"
	body := fmt.Sprintf(`亲爱的 %s：

恭喜您！您在Linux.do Review系统中提交的邀请码申请已通过审核。

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
您的邀请码是：%s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

请妥善保管此邀请码，并在Linux.do网站上使用完成注册。

注意事项：
1. 每个邀请码只能使用一次
2. 请在有效期内使用
3. 请勿将邀请码分享给他人

如有任何问题，欢迎联系我们。

祝您使用愉快！

---
此邮件由Linux.do Review系统自动发送，请勿回复。
`, username, inviteCode)

	if err := s.send(to, subject, body); err != nil {
		log.Printf("[EmailService] 发送邀请码邮件失败: %v", err)
		return err
	}

	log.Printf("[EmailService] 邀请码邮件已发送给 %s", to)
	return nil
}

// SendRejectionNotification 发送拒绝通知邮件
func (s *EmailService) SendRejectionNotification(to, username, postTitle, reason string) error {
	if !s.enabled {
		log.Printf("[EmailService] SMTP未配置,跳过发送拒绝通知邮件给 %s", to)
		return nil
	}

	subject := "关于您的Linux.do邀请码申请"

	reasonText := "未达到审核标准"
	if reason != "" {
		reasonText = reason
	}

	body := fmt.Sprintf(`亲爱的 %s：

感谢您对Linux.do社区的关注。

很遗憾地通知您，您提交的邀请码申请「%s」未能通过审核。

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
原因：%s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

您可以：
1. 重新提交一份更详细的申请
2. 完善您的申请内容后再次申请
3. 联系管理员了解更多详情

我们期待您的再次申请！

---
此邮件由Linux.do Review系统自动发送，请勿回复。
`, username, postTitle, reasonText)

	if err := s.send(to, subject, body); err != nil {
		log.Printf("[EmailService] 发送拒绝通知邮件失败: %v", err)
		return err
	}

	log.Printf("[EmailService] 拒绝通知邮件已发送给 %s", to)
	return nil
}

// SendStatusNotification 发送状态变更通知邮件
func (s *EmailService) SendStatusNotification(to, username, postTitle, statusText, message string) error {
	if !s.enabled {
		return nil
	}

	subject := fmt.Sprintf("您的申请「%s」状态已更新", postTitle)
	body := fmt.Sprintf(`亲爱的 %s：

您的邀请码申请「%s」状态已更新为：%s

%s

如有疑问，请联系管理员。

---
此邮件由Linux.do Review系统自动发送，请勿回复。
`, username, postTitle, statusText, message)

	return s.send(to, subject, body)
}

// send 发送邮件
func (s *EmailService) send(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.user, s.password, s.host)

	// 发件人地址（用于SMTP信封）
	fromAddr := s.from
	if fromAddr == "" {
		fromAddr = s.user
	}

	// 显示的发件人（带名称）
	displayFrom := fmt.Sprintf("Linux.do Review <%s>", fromAddr)

	// 构建邮件头
	headers := make(map[string]string)
	headers["From"] = displayFrom
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=UTF-8"

	// 构建完整消息
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	return smtp.SendMail(addr, auth, fromAddr, []string{to}, []byte(message))
}

// SendNotification 发送通知邮件
func (s *EmailService) SendNotification(to, subject, body string) error {
	if !s.enabled {
		log.Printf("[EmailService] SMTP未配置,跳过发送通知邮件给 %s", to)
		return nil
	}
	return s.send(to, subject, body)
}
