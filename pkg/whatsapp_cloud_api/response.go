package whatsapp_cloud_api

type (
	MetaMessageResponse struct {
		MessagingProduct string     `json:"messaging_product"`
		Contracts        []Contract `json:"contracts"`
		Messages         []Message  `json:"messages"`
	}

	Contract struct {
		Input string `json:"input"`
		WaId  string `json:"wa_id"`
	}

	Message struct {
		Id            string `json:"id"`
		MessageStatus string `json:"message_status"`
	}

	MetaErrorResponse struct {
		Error ErrorDetail `json:"error"`
	}

	ErrorDetail struct {
		Message      string `json:"message"`
		Type         string `json:"type"`
		Code         int    `json:"code"`
		ErrorData    `json:"error_data"`
		ErrorSubcode int
		FbtraceID    string `json:"fbtrace_id"`
	}

	ErrorData struct {
		MessagingProduct string `json:"messaging_product"`
		Details          string `json:"details"`
	}
)
