package mail

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"wealth-vault/user-service/configs"
	"wealth-vault/user-service/internal/domain"

	"gopkg.in/gomail.v2"
)

type gomailClient struct {
	dialer *gomail.Dialer
	from   string
}

type NotificationClient interface {
	SendShareInvitation(ctx context.Context, req domain.SendEmailRequest) error
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

func (g *gomailClient) SendShareInvitation(ctx context.Context, req domain.SendEmailRequest) error {
	m := gomail.NewMessage()
	m.SetHeader("From", g.from)
	m.SetHeader("To", req.ToEmail)
	m.SetHeader("Subject", fmt.Sprintf("Wealth Vault: Shared %s with you", req.AssetName))
	m.SetBody("text/html", g.generateHTML(req))

	return g.dialer.DialAndSend(m)
}

func (g *gomailClient) generateHTML(req domain.SendEmailRequest) string {
	var detailsBuilder strings.Builder
	if len(req.ItemDetail) > 0 {
		detailsBuilder.WriteString(`<div style="background-color: #f9fafb; padding: 15px; border-radius: 5px; margin: 15px 0;">`)
		for key, value := range req.ItemDetail {
			line := fmt.Sprintf(`<p style="margin: 5px 0; color: #555;"><strong>%s:</strong> %s</p>`, key, value)
			detailsBuilder.WriteString(line)
		}
		detailsBuilder.WriteString(`</div>`)
	}

	return fmt.Sprintf(`
		<div style="font-family: 'Sarabun', sans-serif; max-width: 600px; margin: auto; border: 1px solid #ddd; border-radius: 8px; overflow: hidden;">
			<div style="padding: 20px;">
				<h2 style="color: #333; margin-top: 0;">คุณได้รับคำเชิญให้ดู: %s</h2>
				<p style="color: #777;">ประเภททรัพย์สิน: <span style="background: #eee; padding: 2px 6px; border-radius: 4px;">%s</span></p>
				
				%s
				
				<br/>
				<a  style="display: block; text-align: center; background: #d97706; color: white; padding: 12px; text-decoration: none; border-radius: 5px; font-weight: bold;">
					ดูรายละเอียดทรัพย์สิน
				</a>
				<p style="font-size: 12px; color: #999; text-align: center; margin-top: 20px;">
					Wealth Vault Application
				</p>
			</div>
		</div>
	`, req.AssetName, req.AssetType, detailsBuilder.String())
}
