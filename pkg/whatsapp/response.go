package whatsapp

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
)
