package errs

import "fmt"

const (
	UserSettingNotFoundMessage = "User settings for user %q could not be found."
)

type UserSettingFoundErr struct {
	subject string
}

func NewUserSettingNotFoundErr(subject string) UserSettingFoundErr {
	return UserSettingFoundErr{subject: subject}
}

func (e UserSettingFoundErr) Error() string {
	return fmt.Sprintf(UserSettingNotFoundMessage, e.subject)
}
