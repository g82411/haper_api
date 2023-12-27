package utils

type Claims struct {
	Email        string `json:"email"`
	Sub          string `json:"sub"`
	CogUsername  string `json:"cognito:username"`
	Name         string `json:"name"`
	IsDoneSurvey string `json:"custom:isDoneSurvey"`
}
