package response

type AuthenticationResult struct {
	Status    int         `json:"status"`
	PI        *string      `json:"pi"`
	PUid      string       `json:"puid"`
}