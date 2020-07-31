package chat

type Client struct {
	botId       string
	token       string
	tokenType   string
	apiEndpoint string
}

type Profile struct {
	Email          string `json:"email"`
	Nickname       string `json:"nickname"`
	AccountId      string `json:"one_id"`
	ProfilePicture string `json:"profilepicture"`
}

type Friend struct {
	OneEmail    string `json:"one_email"`
	UserId      string `json:"user_id"`
	AccountId   string `json:"one_id"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
}

type Choice struct {
	Label string `json:"label"`
	Type  string `json:"type"`
	Url   string `json:"url"`
	Size  string `json:"size"`
}
type Elements struct {
	Image   string   `json:"image"`
	Title   string   `json:"title"`
	Detail  string   `json:"detail"`
	Choices []Choice `json:"choice"`
}

type QuickReply struct {
	Label   string      `json:"label"`
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}
