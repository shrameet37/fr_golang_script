package models

type SupportCommandResponse struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}

type CommandSuccessResponse struct {
	Type    string                `json:"type"`
	Message CommandSuccessMessage `json:"message"`
}

type CommandSuccessMessage struct {
	SuccessMessage string `json:"successMessage"`
}

type CommandErrorResponse struct {
	Type    string              `json:"type"`
	Message CommandErrorMessage `json:"message"`
}

type CommandErrorMessage struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}
