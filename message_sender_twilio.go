package messagesendertwilio

import (
	"context"
	"github.com/applicaset/sms-svc"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
)

type messageSender struct {
	from       string
	accountSid string
	authToken  string
	urlStr     string
}

func (ms *messageSender) Send(ctx context.Context, phoneNumber, message string) error {
	msgData := url.Values{}
	msgData.Set("To", phoneNumber)
	msgData.Set("From", ms.from)
	msgData.Set("Body", message)

	msgDataReader := strings.NewReader(msgData.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ms.urlStr, msgDataReader)
	if err != nil {
		return errors.Wrap(err, "error on create new http request")
	}

	req.SetBasicAuth(ms.accountSid, ms.authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error on do http request")
	}

	defer func() { _ = res.Body.Close() }()

	if res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	return errors.Errorf("http response is '%d'", res.StatusCode)
}

func New(accountSid, authToken, from string) smssvc.MessageSender {
	ms := messageSender{
		from:       from,
		accountSid: accountSid,
		authToken:  authToken,
		urlStr:     "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json",
	}

	return &ms
}
