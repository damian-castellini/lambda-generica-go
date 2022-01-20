package aws

import "github.com/aws/aws-sdk-go/aws/session"

func NewSession() *session.Session {
	sessOpts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
	return session.Must(session.NewSessionWithOptions(sessOpts))
}
