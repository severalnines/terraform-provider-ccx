package resources

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
	"github.com/stretchr/testify/mock"
)

func TestParameterGroup_Create(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		m, p := mockProvider(t)

		pgCreated := ccx.ParameterGroup{
			ID:              "parameter-group-id",
			Name:            "asteroid",
			DatabaseVendor:  "mariadb",
			DatabaseVersion: "10.11",
			DatabaseType:    "galera",
			DbParameters: map[string]string{
				"max_connections": "100",
				"sql_mode":        "STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION",
			},
		}

		m.content.EXPECT().DBVendors(mock.Anything).Return([]ccx.DBVendorInfo{
			{
				Name:           "mariadb",
				Code:           "mariadb",
				DefaultVersion: "10.11",
				Versions:       []string{"10.11"},
				Types:          []ccx.DBVendorInfoType{{Name: "galera", Code: "galera"}},
				NumNodes:       []int{1, 2, 3, 5},
			},
		}, nil)

		m.parameterGroup.EXPECT().Create(mock.Anything, ccx.ParameterGroup{
			Name:            "asteroid",
			DatabaseVendor:  "mariadb",
			DatabaseVersion: "10.11",
			DatabaseType:    "galera",
			DbParameters: map[string]string{
				"max_connections": "100",
				"sql_mode":        "STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION",
			},
		}).Return(&pgCreated, nil)

		m.parameterGroup.EXPECT().Read(mock.Anything, "parameter-group-id").Return(&pgCreated, nil)

		m.parameterGroup.EXPECT().Delete(mock.Anything, "parameter-group-id").Return(nil)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			PreCheck: func() {
			},
			ProviderFactories: map[string]func() (*schema.Provider, error){
				"ccx": func() (*schema.Provider, error) {
					return p, nil
				},
			},
			Steps: []resource.TestStep{
				{
					Config: `
resource "ccx_parameter_group" "asteroid" {
    name = "asteroid"
    database_vendor = "mariadb"
    database_version = "10.11"
    database_type = "galera"

    parameters = {
      max_connections = 100
      sql_mode = "STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION"
    }
}
`,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("ccx_parameter_group.asteroid", "id", "parameter-group-id"),
						resource.TestCheckResourceAttr("ccx_parameter_group.asteroid", "database_vendor", "mariadb"),
						resource.TestCheckResourceAttr("ccx_parameter_group.asteroid", "database_version", "10.11"),
						resource.TestCheckResourceAttr("ccx_parameter_group.asteroid", "database_type", "galera"),
						resource.TestCheckResourceAttr("ccx_parameter_group.asteroid", "parameters.max_connections", "100"),
						resource.TestCheckResourceAttr("ccx_parameter_group.asteroid", "parameters.sql_mode", "STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION"),
					),
				},
			},
		})

		m.AssertExpectations(t)
	})
}
