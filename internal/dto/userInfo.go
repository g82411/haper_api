package dto

type UserInfo struct {
	Email            string `json:"email"`
	Sub              string `json:"sub"`
	Name             string `json:"name"`
	IsDoneSurvey     string `json:"custom:isDoneSurvey"`
	InternalUserName string `json:"cognito:username"`
}
