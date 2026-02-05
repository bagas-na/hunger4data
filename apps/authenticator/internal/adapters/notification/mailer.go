package notification

import (
	"authenticator/internal/adapters/repo"
	"context"
	"fmt"
	notifyv1 "hunger4data/pb/notification"
)

type Mailer interface {
	SendRegistrationActivationLink(ctx context.Context, u *repo.User) error
}

type mailer struct {
	client notifyv1.EmailServiceClient
}

func NewMailer(client notifyv1.EmailServiceClient) Mailer {
	return &mailer{
		client: client,
	}
}

func (m *mailer) SendRegistrationActivationLink(ctx context.Context, u *repo.User) error {

	appName := "Hunger4Data"
	userName := u.Username
	activationURL := fmt.Sprintf("http://localhost:9000/auth/activate/%s", u.ActivationString)

	emailHTML := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Activate your account</title>
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
                Activate your %s account
              </td>
            </tr>

            <tr>
              <td style="font-size: 14px; line-height: 1.5; padding-bottom: 24px">
                Hi %s,
                <br /><br />
                Thanks for signing up for %s.
                Click the button below to activate your account and get started.
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
                  Activate account
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
                If you didn’t sign up for %s, you can safely ignore this email.
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
		appName, // App name (header)
		userName,
		appName, // App name (body)
		activationURL,
		appName, // App name (ignore note)
		appName, // App name (footer)
	)

	email := notifyv1.SendEmailRequest{
		To:      u.Username,
		Subject: "Activate Your Hunger4Data Account",
		Body:    emailHTML,
	}

	_, err := m.client.SendTransactionEmail(ctx, &email)
	return err
}
