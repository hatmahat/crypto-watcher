package whatsapp_cloud_api

type (
	MetaMessageRequest struct {
		MessaingProduct string `json:"messaging_product"`
		To              string `json:"to"`
		Type            string `json:"type"`
		Template        `json:"template"`
	}

	Template struct {
		Name       string `json:"name"`
		Language   `json:"language"`
		Components []Component `json:"components"`
	}

	Language struct {
		Code string `json:"code"`
	}

	Component struct {
		Type       string      `json:"type"`
		Parameters []Parameter `json:"parameters"`
	}

	Parameter struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
)
