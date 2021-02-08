package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnResourcesFromInstanaAPI(t *testing.T) {
	api := NewInstanaAPI("api-token", "endpoint")

	t.Run("Should return CustomEventSpecification instance", func(t *testing.T) {
		resource := api.CustomEventSpecifications()

		require.NotNil(t, resource)
	})

	t.Run("Should return BuiltinEventSpecifications instance", func(t *testing.T) {
		resource := api.BuiltinEventSpecifications()

		require.NotNil(t, resource)
	})
	t.Run("Should return UserRole instance", func(t *testing.T) {
		resource := api.UserRoles()

		require.NotNil(t, resource)
	})
	t.Run("Should return ApplicationConfig instance", func(t *testing.T) {
		resource := api.ApplicationConfigs()

		require.NotNil(t, resource)
	})
	t.Run("Should return AlertingChannel instance", func(t *testing.T) {
		resource := api.AlertingChannels()

		require.NotNil(t, resource)
	})
	t.Run("Should return AlertingConfiguration instance", func(t *testing.T) {
		resource := api.AlertingConfigurations()

		require.NotNil(t, resource)
	})
	t.Run("Should return SliConfig instance", func(t *testing.T) {
		resource := api.SliConfigs()

		require.NotNil(t, resource)
	})
	t.Run("Should return WebsiteMonitoringConfig instance", func(t *testing.T) {
		resource := api.WebsiteMonitoringConfig()

		require.NotNil(t, resource)
	})
}
