package api

import (
	"context"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

type deployWizardResponse struct {
	Instance struct {
		InstanceSizes map[string][]ccx.InstanceSize `json:"instance_sizes"`
	} `json:"instance"`
}

func (svc *ContentService) InstanceSizes(ctx context.Context) (map[string][]ccx.InstanceSize, error) {
	var rs deployWizardResponse

	err := svc.client.Get(ctx, "/api/content/api/v1/deploy-wizard", &rs)
	if err != nil {
		return nil, err
	}

	return rs.Instance.InstanceSizes, nil
}
