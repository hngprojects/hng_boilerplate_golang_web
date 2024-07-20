package send_invites

import "time"

func MockSendInvite(email string, orgName string, expiresAt time.Time) (string, error) {
	return MockGenerateInviteLink(), nil
}


func MockGenerateInviteLink() string {
	return "http://localhost:8080/invite/1234"
}