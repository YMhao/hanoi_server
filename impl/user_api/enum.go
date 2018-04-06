package user_api

type UserNameType string

const (
	USER_NAME_TYPE_EMAIL        UserNameType = "EMAIL"
	USER_NAME_TYPE_MOBILE_PHONE UserNameType = "MOBILE_PHTONE"
	USER_NAME_TYPE_UNKNOWN      UserNameType = "UNKNOWN"
)
