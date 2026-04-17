package auth

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenInfo struct {
	ID           string  `json:"id"`
	Document     string  `json:"document"`
	Name         string  `json:"name"`
	AccessGroups []int16 `json:"accessGroups"`
}

type User struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Document     string  `json:"document"`
	AccessGroups []int16 `json:"accessGroups"`
}
