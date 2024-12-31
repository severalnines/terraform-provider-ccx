package api

type ContentService struct {
	client HttpClient
}

func Content(client HttpClient) (*ContentService, error) {
	c := ContentService{
		client: client,
	}

	return &c, nil
}
