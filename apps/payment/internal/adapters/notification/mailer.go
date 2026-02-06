package notification

import (
	"context"
	"fmt"
	notifyv1 "hunger4data/pb/notification"
	"payment-service/internal/adapters/db"
	"payment-service/internal/utils"
)

type Mailer interface {
	SendCheckoutURL(ctx context.Context, p *db.Payment, url string) error
}

type mailer struct {
	client notifyv1.EmailServiceClient
}

func NewMailer(client notifyv1.EmailServiceClient) Mailer {
	return &mailer{
		client: client,
	}
}

func (m *mailer) SendCheckoutURL(ctx context.Context, p *db.Payment, url string) error {

	appName := "Hunger4Data"
	userName := p.User.Username
	countryName := utils.ISO3166Alpha3[p.CountryCode]

	checkoutURL := url

	emailHTML := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Complete Your Payment</title>
  </head>
  <body
    style="
      margin: 0;
      padding: 0;
      background-color: #f9fafb;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI',
        Roboto, Helvetica, Arial, sans-serif;
      color: #111827;
    "
  >
    <table width="100%%" cellpadding="0" cellspacing="0">
      <tr>
        <td align="center" style="padding: 32px 16px">
          <table
            width="100%%"
            cellpadding="0"
            cellspacing="0"
            style="
              max-width: 480px;
              background-color: #ffffff;
              border-radius: 8px;
              padding: 32px;
            "
          >
            <tr>
              <td style="font-size: 18px; font-weight: 600; padding-bottom: 16px">
                Complete Your Donation for %s
              </td>
            </tr>

            <tr>
              <td style="font-size: 14px; line-height: 1.5; padding-bottom: 24px">
                Hi %s,
                <br /><br />
                You’re one step away from completing your donation to <strong>%s</strong>.
                Click the button below to proceed with payment.
              </td>
            </tr>

            <tr>
              <td align="center" style="padding-bottom: 24px">
                <a
                  href="%s"
                  style="
                    display: inline-block;
                    background-color: #111827;
                    color: #ffffff;
                    text-decoration: none;
                    font-size: 14px;
                    font-weight: 500;
                    padding: 12px 20px;
                    border-radius: 6px;
                  "
                >
                  Pay Now
                </a>
              </td>
            </tr>

            <tr>
              <td
                style="
                  font-size: 12px;
                  line-height: 1.5;
                  color: #6b7280;
                  padding-bottom: 16px;
                "
              >
                If you didn’t intend to make this donation, you can safely ignore this email.
              </td>
            </tr>

            <tr>
              <td
                style="
                  font-size: 12px;
                  color: #9ca3af;
                  border-top: 1px solid #e5e7eb;
                  padding-top: 16px;
                "
              >
                — The %s team
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>
`,
		appName, // Header
		userName,
		countryName,
		checkoutURL, // Stripe Checkout URL
		appName,     // Footer
	)

	email := notifyv1.SendEmailRequest{
		To:      p.User.Username,
		Subject: "Complete Your Hunger4Data Donation",
		Body:    emailHTML,
	}

	_, err := m.client.SendTransactionEmail(ctx, &email)
	return err
}
