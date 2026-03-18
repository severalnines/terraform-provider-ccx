package ccx

type ContentClient struct {
	client HTTPClient
}

func NewContentClient(client HTTPClient) (ContentService, error) {
	c := ContentClient{
		client: client,
	}

	return &c, nil
}
