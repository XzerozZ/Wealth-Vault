package mail

import (
	"context"
	"fmt"
	"strconv"
	"wealth-vault/auth-service/configs"
	"wealth-vault/auth-service/internal/domain"

	"gopkg.in/gomail.v2"
)

type gomailClient struct {
	dialer *gomail.Dialer
	from   string
}

type NotificationClient interface {
	SendOTP(ctx context.Context, req domain.SendEmailRequest) error
}

func NewMailClient(cfg configs.Mail) NotificationClient {
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		return nil
	}

	d := gomail.NewDialer(cfg.Host, port, cfg.Sender, cfg.Key)

	return &gomailClient{
		dialer: d,
		from:   cfg.Sender,
	}
}

func (g *gomailClient) SendOTP(ctx context.Context, req domain.SendEmailRequest) error {
	m := gomail.NewMessage()
	m.SetHeader("From", g.from)
	m.SetHeader("To", req.ToEmail)
	m.SetHeader("Subject", fmt.Sprintf("Wealth Vault: Hello %v", req.ToEmail))
	m.SetBody("text/html", g.generateOTPHTML(req))

	return g.dialer.DialAndSend(m)
}

func (g *gomailClient) generateOTPHTML(req domain.SendEmailRequest) string {
	expiry := req.ExpiredAt
	if expiry == "" {
		expiry = "5 นาที"
	}

	return fmt.Sprintf(`
		<div style="font-family: 'Sarabun', sans-serif; max-width: 600px; margin: auto; border: 1px solid #ddd; border-radius: 8px; overflow: hidden; background-color: #ffffff;">
			<div style="background-color: #d97706; padding: 15px; text-align: center;">
				<h2 style="color: white; margin: 0; font-size: 20px;">รหัสยืนยันตัวตน (OTP)</h2>
			</div>
			
			<div style="padding: 30px 20px; text-align: center;">
				<p style="color: #555; font-size: 16px; margin-bottom: 20px;">
					โปรดใช้รหัสด้านล่างเพื่อดำเนินการยืนยันตัวตนในแอปพลิเคชัน Wealth Vault
				</p>
				
				<div style="background-color: #f3f4f6; border-radius: 8px; padding: 20px; margin: 20px 0; letter-spacing: 5px;">
					<span style="font-size: 36px; font-weight: bold; color: #333; display: block;">%s</span>
				</div>
				
				<div style="margin-top: 20px; font-size: 14px; color: #777;">
					<p style="margin: 5px 0; color: #dc2626;">*รหัสนี้มีอายุการใช้งาน %s และห้ามเปิดเผยแก่ผู้อื่น</p>
				</div>

				<hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;" />
				
				<p style="font-size: 12px; color: #999;">
					หากคุณไม่ได้ทำรายการนี้ โปรดติดต่อเจ้าหน้าที่ทันที<br/>
					Wealth Vault Application
				</p>
			</div>
		</div>
	`, req.OTP, expiry)
}
