package instana_test

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

func TestCustomEventSpecificationResource(t *testing.T) {
	unitTest := &customerEventSpecificationUnitTest{}
	t.Run("CRUD integration test of with Entity Verification Rule", customerEventSpecificationIntegrationTestWithEntityVerificationRule().testCrud)
	t.Run("CRUD integration test of with System Rule", customerEventSpecificationIntegrationTestWithSystemRule().testCrud)
	t.Run("CRUD integration test of with Threshold Rule and Metric Name and Rollup", customerEventSpecificationIntegrationTestWithThresholdRuleAndMetricNameAndRollup().testCrud)
	t.Run("CRUD integration test of with Threshold Rule and Metric Pattern", customerEventSpecificationIntegrationTestWithThresholdRuleAndMetricPattern().testCrud)
	t.Run("schema should be valid", unitTest.schemaShouldBeValid)
	t.Run("should have schema version 0", unitTest.shouldHaveSchemaVersion0)
	t.Run("should have no state upgrader", unitTest.shouldHaveNoStateUpgraders)
	t.Run("should have correct resource name", unitTest.shouldHaveCorrectResourceName)
	t.Run("should map entity verification rule to state", unitTest.shouldMapEntityVerificationRuleToState)
	t.Run("should map system rule to state", unitTest.shouldMapSystemRuleToState)
	t.Run("should map threshold rule and metric name to state", unitTest.shouldMapThresholdRuleAndMetricNameToState)
	t.Run("should map threshold rule and metric pattern to state", unitTest.shouldMapThresholdRuleAndMetricPatternToState)
	t.Run("should fail to map rule when severity is not valid", unitTest.shouldFailToMapRuleWhenSeverityIsNotValid)
	t.Run("should fail to map rule when rule type is not valid", unitTest.shouldFailToMapRuleWhenRuleTypeIsNotValid)
	t.Run("should map state of entity verification rule to data model", unitTest.shouldMapStateOfEntityVerificationRuleToDataModel)
	t.Run("should fail to map state of entity verification rule when severity is not valid", unitTest.shouldFailToMapStateOfEntityVerificationRuleToDataModelWhenSeverityIsNotValid)
	t.Run("should map state of system rule to data model", unitTest.shouldMapStateOfSystemRuleToDataModel)
	t.Run("should fail to map state of system rule when severity is not valid", unitTest.shouldFailToMapStateOfSystemRuleToDataModelWhenSeverityIsNotValid)
	t.Run("should map state of threshold rule with metric name to data model", unitTest.shouldMapStateOfThresholdRuleWithMetricNameToDataModel)
	t.Run("should map state of threshold rule with metric pattern to data model", unitTest.shouldMapStateOfThresholdRuleWithMetricPatternToDataModel)
	t.Run("should fail to map state of threshold rule when severity is not valid", unitTest.shouldFailToMapStateOfThresholdRuleToDataModelWhenSeverityIsNotValid)
	t.Run("should fail to map state when no rule is provided", unitTest.shouldFailToMapStateWhenNoRuleIsProvided)
}

const (
	customEventSpecificationConfigResourceName            = "instana_custom_event_specification.example"
	customEventSpecificationRuleFieldPattern              = "%s.0.%s.0.%s"
	customEventSpecificationRuleMetricPatternFieldPattern = "%s.0.%s.0.%s.0.%s"

	customEventSpecificationWithThresholdRuleSeverity       = "warning"
	customEventSpecificationWithThresholdRuleID             = "custom-system-event-id"
	customEventSpecificationWithThresholdRuleEntityType     = "entity_type"
	customEventSpecificationWithThresholdRuleQuery          = "query"
	customEventSpecificationWithThresholdRuleExpirationTime = 60000
	customEventSpecificationWithThresholdRuleDescription    = "description"
	customEventSpecificationWithThresholdRuleMetricName     = "metric_name"
	customEventSpecificationWithThresholdRuleRollup         = 40000
	customEventSpecificationWithThresholdRuleWindow         = 60000
	customEventSpecificationWithThresholdRuleAggregation    = "sum"
	customEventSpecificationWithThresholdRuleConditionValue = float64(1.2)

	entityVerificationRuleEntityType = "host"
	systemRuleEntityType             = "any"
)

func customerEventSpecificationIntegrationTestWithEntityVerificationRule() *customerEventSpecificationIntegrationTest {
	resourceTemplate := `
resource "instana_custom_event_specification" "example" {
  name = "name %d"
  entity_type = "host"
  query = "query"
  enabled = true
  triggering = true
  description = "description"
  expiration_time = "60000"
  rules {
    entity_verification {
      severity = "warning"
      matching_entity_label = "matching-entity-label"
      matching_entity_type = "matching-entity-type"
      matching_operator = "startsWith"
	  offline_duration = 60000
    }
  }
}`

	httpServerResponseTemplate := `
{
	"id" : "%s",
	"name" : "name %d",
	"entityType" : "host",
	"query" : "query",
	"enabled" : true,
	"triggering" : true,
	"description" : "description",
	"expirationTime" : 60000,
	"rules" : [{ 
		"ruleType" : "entity_verification", 
		"severity" : 5, 
		"matchingEntityLabel" : "matching-entity-label", 
		"matchingEntityType" : "matching-entity-type", 
		"matchingOperator" : "startsWith", 
		"offlineDuration" : 60000 
	}]
}`

	return newCustomerEventSpecificationIntegrationTest(
		"host",
		resourceTemplate,
		customEventSpecificationConfigResourceName,
		httpServerResponseTemplate,
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldEntityVerificationRule, CustomEventSpecificationRuleFieldSeverity), customEventSpecificationWithThresholdRuleSeverity),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldEntityVerificationRule, CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityLabel), "matching-entity-label"),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldEntityVerificationRule, CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityType), "matching-entity-type"),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldEntityVerificationRule, CustomEventSpecificationEntityVerificationRuleFieldMatchingOperator), "startsWith"),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldEntityVerificationRule, CustomEventSpecificationEntityVerificationRuleFieldOfflineDuration), "60000"),
		},
	)
}

func customerEventSpecificationIntegrationTestWithSystemRule() *customerEventSpecificationIntegrationTest {
	resourceTemplate := `
resource "instana_custom_event_specification" "example" {
  name = "name %d"
  entity_type = "any"
  query = "query"
  enabled = true
  triggering = true
  description = "description"
  expiration_time = "60000"
  rules {
    system {
      severity = "warning"
      system_rule_id = "system_rule_id"
    }
  }
}`

	httpServerResponseTemplate := `
{
	"id" : "%s",
	"name" : "name %d",
	"entityType" : "any",
	"query" : "query",
	"enabled" : true,
	"triggering" : true,
	"description" : "description",
	"expirationTime" : 60000,
	"rules" : [{
		"ruleType" : "system", 
		"severity" : 5, 
		"systemRuleId" : "system_rule_id"
	}]
}`

	return newCustomerEventSpecificationIntegrationTest(
		"any",
		resourceTemplate,
		customEventSpecificationConfigResourceName,
		httpServerResponseTemplate,
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldSystemRule, CustomEventSpecificationRuleFieldSeverity), customEventSpecificationWithThresholdRuleSeverity),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldSystemRule, CustomEventSpecificationSystemRuleFieldSystemRuleId), "system_rule_id"),
		},
	)
}

func customerEventSpecificationIntegrationTestWithThresholdRuleAndMetricNameAndRollup() *customerEventSpecificationIntegrationTest {
	resourceTemplate := `
resource "instana_custom_event_specification" "example" {
  name = "name %d"
  entity_type = "entity_type"
  query = "query"
  enabled = true
  triggering = true
  description = "description"
  expiration_time = "60000"
  rules {
    threshold {
      severity = "warning"
      metric_name = "metric_name"
	  aggregation = "sum"
      window = "60000"
      rollup = "40000"
      condition_operator = "="
      condition_value = "1.2"
    }
  }
}`

	httpServerResponseTemplate := `
{
	"id" : "%s",
	"name" : "name %d",
	"entityType" : "entity_type",
	"query" : "query",
	"enabled" : true,
	"triggering" : true,
	"description" : "description",
	"expirationTime" : 60000,
	"rules" : [{
		"ruleType" : "threshold", 
		"severity" : 5, 
		"metricName" : "metric_name", 
		"aggregation" : "sum", 
		"window" : 60000,
		"rollup" : 40000, 
		"conditionOperator" : "=", 
		"conditionValue" : 1.2 
	}]
}`

	return newCustomerEventSpecificationIntegrationTest(
		customEventSpecificationWithThresholdRuleEntityType,
		resourceTemplate,
		customEventSpecificationConfigResourceName,
		httpServerResponseTemplate,
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldThresholdRule, CustomEventSpecificationRuleFieldSeverity), customEventSpecificationWithThresholdRuleSeverity),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldThresholdRule, CustomEventSpecificationThresholdRuleFieldMetricName), customEventSpecificationWithThresholdRuleMetricName),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldThresholdRule, CustomEventSpecificationThresholdRuleFieldAggregation), string(customEventSpecificationWithThresholdRuleAggregation)),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldThresholdRule, CustomEventSpecificationThresholdRuleFieldWindow), strconv.FormatInt(customEventSpecificationWithThresholdRuleWindow, 10)),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldThresholdRule, CustomEventSpecificationThresholdRuleFieldRollup), strconv.FormatInt(customEventSpecificationWithThresholdRuleRollup, 10)),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldThresholdRule, CustomEventSpecificationThresholdRuleFieldConditionOperator), string("=")),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldThresholdRule, CustomEventSpecificationThresholdRuleFieldConditionValue), "1.2"),
		},
	)
}

func customerEventSpecificationIntegrationTestWithThresholdRuleAndMetricPattern() *customerEventSpecificationIntegrationTest {
	resourceTemplate := `
resource "instana_custom_event_specification" "example" {
  name = "name %d"
  entity_type = "entity_type"
  query = "query"
  enabled = true
  triggering = true
  description = "description"
  expiration_time = "60000"
  rules {
    threshold {
      severity = "warning"
      metric_pattern {
      	prefix = "prefix"
		postfix = "postfix"
		placeholder = "placeholder"
		operator = "startsWith"
      }
	  aggregation = "sum"
      window = "60000"
      rollup = "40000"
      condition_operator = "="
      condition_value = "1.2"
    }
  }
}`

	httpServerResponseTemplate := `
{
	"id" : "%s",
	"name" : "name %d",
	"entityType" : "entity_type",
	"query" : "query",
	"enabled" : true,
	"triggering" : true,
	"description" : "description",
	"expirationTime" : 60000,
	"rules" : [{
		"ruleType" : "threshold", 
		"severity" : 5, 
		"metricPattern" : { 
			"prefix" : "prefix", 
			"postfix" : "postfix", 
			"placeholder" : "placeholder", 
			"operator" : "startsWith"
		}, 
		"aggregation" : "sum", 
		"window" : 60000,
		"rollup" : 40000, 
		"conditionOperator" : "=", 
		"conditionValue" : 1.2 
	}]
}`

	return newCustomerEventSpecificationIntegrationTest(
		customEventSpecificationWithThresholdRuleEntityType,
		resourceTemplate,
		customEventSpecificationConfigResourceName,
		httpServerResponseTemplate,
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldThresholdRule, CustomEventSpecificationRuleFieldSeverity), customEventSpecificationWithThresholdRuleSeverity),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleMetricPatternFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldThresholdRule, CustomEventSpecificationThresholdRuleFieldMetricPattern, CustomEventSpecificationThresholdRuleFieldMetricPatternPrefix), "prefix"),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleMetricPatternFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldThresholdRule, CustomEventSpecificationThresholdRuleFieldMetricPattern, CustomEventSpecificationThresholdRuleFieldMetricPatternPostfix), "postfix"),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleMetricPatternFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldThresholdRule, CustomEventSpecificationThresholdRuleFieldMetricPattern, CustomEventSpecificationThresholdRuleFieldMetricPatternPlaceholder), "placeholder"),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleMetricPatternFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldThresholdRule, CustomEventSpecificationThresholdRuleFieldMetricPattern, CustomEventSpecificationThresholdRuleFieldMetricPatternOperator), "startsWith"),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldThresholdRule, CustomEventSpecificationThresholdRuleFieldAggregation), string(customEventSpecificationWithThresholdRuleAggregation)),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldThresholdRule, CustomEventSpecificationThresholdRuleFieldWindow), strconv.FormatInt(customEventSpecificationWithThresholdRuleWindow, 10)),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldThresholdRule, CustomEventSpecificationThresholdRuleFieldRollup), strconv.FormatInt(customEventSpecificationWithThresholdRuleRollup, 10)),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldThresholdRule, CustomEventSpecificationThresholdRuleFieldConditionOperator), string("=")),
			resource.TestCheckResourceAttr(customEventSpecificationConfigResourceName, fmt.Sprintf(customEventSpecificationRuleFieldPattern, CustomEventSpecificationFieldRules, CustomEventSpecificationFieldThresholdRule, CustomEventSpecificationThresholdRuleFieldConditionValue), "1.2"),
		},
	)
}

func newCustomerEventSpecificationIntegrationTest(entityType string, resourceTemplate string, resourceName string, serverResponseTemplate string, useCaseSpecificChecks []resource.TestCheckFunc) *customerEventSpecificationIntegrationTest {
	return &customerEventSpecificationIntegrationTest{
		entityType:             entityType,
		resourceTemplate:       resourceTemplate,
		resourceName:           resourceName,
		serverResponseTemplate: serverResponseTemplate,
		useCaseSpecificChecks:  useCaseSpecificChecks,
	}
}

type customerEventSpecificationIntegrationTest struct {
	entityType             string
	resourceTemplate       string
	resourceName           string
	serverResponseTemplate string
	useCaseSpecificChecks  []resource.TestCheckFunc
}

func (r *customerEventSpecificationIntegrationTest) testCrud(t *testing.T) {
	httpServer := createMockHttpServerForResource(restapi.CustomEventSpecificationResourcePath, r.serverResponseTemplate)
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: appendProviderConfig(fmt.Sprintf(r.resourceTemplate, 0), httpServer.GetPort()),
				Check:  r.createTestCheckFunctions(0),
			},
			testStepImport(r.resourceName),
			{
				Config: appendProviderConfig(fmt.Sprintf(r.resourceTemplate, 1), httpServer.GetPort()),
				Check:  r.createTestCheckFunctions(1),
			},
			testStepImport(r.resourceName),
		},
	})
}

func (r *customerEventSpecificationIntegrationTest) createTestCheckFunctions(iteration int) resource.TestCheckFunc {
	defaultCheckFunctions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(r.resourceName, "id"),
		resource.TestCheckResourceAttr(r.resourceName, CustomEventSpecificationFieldName, formatResourceName(iteration)),
		resource.TestCheckResourceAttr(r.resourceName, CustomEventSpecificationFieldEntityType, r.entityType),
		resource.TestCheckResourceAttr(r.resourceName, CustomEventSpecificationFieldQuery, customEventSpecificationWithThresholdRuleQuery),
		resource.TestCheckResourceAttr(r.resourceName, CustomEventSpecificationFieldTriggering, trueAsString),
		resource.TestCheckResourceAttr(r.resourceName, CustomEventSpecificationFieldDescription, customEventSpecificationWithThresholdRuleDescription),
		resource.TestCheckResourceAttr(r.resourceName, CustomEventSpecificationFieldExpirationTime, strconv.Itoa(customEventSpecificationWithThresholdRuleExpirationTime)),
		resource.TestCheckResourceAttr(r.resourceName, CustomEventSpecificationFieldEnabled, trueAsString),
	}
	allFunctions := append(defaultCheckFunctions, r.useCaseSpecificChecks...)
	return resource.ComposeTestCheckFunc(allFunctions...)
}

type customerEventSpecificationUnitTest struct{}

func (r *customerEventSpecificationUnitTest) schemaShouldBeValid(t *testing.T) {
	schemaData := NewCustomEventSpecificationResourceHandle().MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaData, t)
	require.Len(t, schemaData, 8)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationFieldName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationFieldEntityType)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationFieldQuery)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationFieldTriggering, false)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationFieldDescription)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(CustomEventSpecificationFieldExpirationTime)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationFieldEnabled, true)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeListOfResource(CustomEventSpecificationFieldRules)

	r.validateRuleSchema(t, schemaData[CustomEventSpecificationFieldRules].Elem.(*schema.Resource).Schema)
}

func (r *customerEventSpecificationUnitTest) validateRuleSchema(t *testing.T, ruleSchema map[string]*schema.Schema) {
	schemaAssert := testutils.NewTerraformSchemaAssert(ruleSchema, t)
	require.Len(t, ruleSchema, 3)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(CustomEventSpecificationFieldEntityVerificationRule)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(CustomEventSpecificationFieldSystemRule)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(CustomEventSpecificationFieldThresholdRule)

	r.validateEntityVerificationRuleSchema(t, ruleSchema[CustomEventSpecificationFieldEntityVerificationRule].Elem.(*schema.Resource).Schema)
	r.validateSystemRuleSchema(t, ruleSchema[CustomEventSpecificationFieldSystemRule].Elem.(*schema.Resource).Schema)
	r.validateThresholdRuleSchema(t, ruleSchema[CustomEventSpecificationFieldThresholdRule].Elem.(*schema.Resource).Schema)
}

func (r *customerEventSpecificationUnitTest) validateEntityVerificationRuleSchema(t *testing.T, ruleSchema map[string]*schema.Schema) {
	require.Len(t, ruleSchema, 5)
	schemaAssert := testutils.NewTerraformSchemaAssert(ruleSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationRuleFieldSeverity)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityType)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationEntityVerificationRuleFieldMatchingOperator)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityLabel)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeInt(CustomEventSpecificationEntityVerificationRuleFieldOfflineDuration)
}

func (r *customerEventSpecificationUnitTest) validateSystemRuleSchema(t *testing.T, ruleSchema map[string]*schema.Schema) {
	require.Len(t, ruleSchema, 2)
	schemaAssert := testutils.NewTerraformSchemaAssert(ruleSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationRuleFieldSeverity)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationSystemRuleFieldSystemRuleId)
}

func (r *customerEventSpecificationUnitTest) validateThresholdRuleSchema(t *testing.T, ruleSchema map[string]*schema.Schema) {
	require.Len(t, ruleSchema, 8)
	schemaAssert := testutils.NewTerraformSchemaAssert(ruleSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationRuleFieldSeverity)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationThresholdRuleFieldMetricName)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(CustomEventSpecificationThresholdRuleFieldMetricPattern)
	r.validateMetricPatternSchema(t, ruleSchema[CustomEventSpecificationThresholdRuleFieldMetricPattern].Elem.(*schema.Resource).Schema)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(CustomEventSpecificationThresholdRuleFieldRollup)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeInt(CustomEventSpecificationThresholdRuleFieldWindow)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationThresholdRuleFieldAggregation)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationThresholdRuleFieldConditionOperator)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeFloat(CustomEventSpecificationThresholdRuleFieldConditionValue)
}

func (r *customerEventSpecificationUnitTest) validateMetricPatternSchema(t *testing.T, metricPatternSchema map[string]*schema.Schema) {
	require.Len(t, metricPatternSchema, 4)
	schemaAssert := testutils.NewTerraformSchemaAssert(metricPatternSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationThresholdRuleFieldMetricPatternPrefix)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationThresholdRuleFieldMetricPatternPostfix)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationThresholdRuleFieldMetricPatternPlaceholder)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationThresholdRuleFieldMetricPatternOperator)
}

func (r *customerEventSpecificationUnitTest) shouldHaveSchemaVersion0(t *testing.T) {
	require.Equal(t, 0, NewCustomEventSpecificationResourceHandle().MetaData().SchemaVersion)
}

func (r *customerEventSpecificationUnitTest) shouldHaveNoStateUpgraders(t *testing.T) {
	resourceHandler := NewCustomEventSpecificationResourceHandle()

	require.Equal(t, 0, len(resourceHandler.StateUpgraders()))
}

func (r *customerEventSpecificationUnitTest) shouldHaveCorrectResourceName(t *testing.T) {
	name := NewCustomEventSpecificationResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_custom_event_specification")
}

func (r *customerEventSpecificationUnitTest) shouldMapEntityVerificationRuleToState(t *testing.T) {
	description := customEventSpecificationWithThresholdRuleDescription
	expirationTime := customEventSpecificationWithThresholdRuleExpirationTime
	query := customEventSpecificationWithThresholdRuleQuery
	matchingEntityLabel := "matching-entity-label"
	matchingEntityType := "matching-entity-type"
	matchingOperator := "is"
	offlineDuration := 1234

	spec := &restapi.CustomEventSpecification{
		ID:             customEventSpecificationWithThresholdRuleID,
		Name:           resourceName,
		EntityType:     entityVerificationRuleEntityType,
		Query:          &query,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Triggering:     true,
		Enabled:        true,
		Rules: []restapi.RuleSpecification{
			{
				DType:               restapi.EntityVerificationRuleType,
				Severity:            restapi.SeverityWarning.GetAPIRepresentation(),
				MatchingEntityLabel: &matchingEntityLabel,
				MatchingEntityType:  &matchingEntityType,
				MatchingOperator:    &matchingOperator,
				OfflineDuration:     &offlineDuration,
			},
		},
	}

	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	sut := NewCustomEventSpecificationResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec)

	require.Nil(t, err)
	require.Equal(t, customEventSpecificationWithThresholdRuleID, resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(CustomEventSpecificationFieldName))
	require.Equal(t, entityVerificationRuleEntityType, resourceData.Get(CustomEventSpecificationFieldEntityType))
	require.Equal(t, customEventSpecificationWithThresholdRuleQuery, resourceData.Get(CustomEventSpecificationFieldQuery))
	require.Equal(t, description, resourceData.Get(CustomEventSpecificationFieldDescription))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldTriggering).(bool))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldEnabled).(bool))

	require.IsType(t, []interface{}{}, resourceData.Get(CustomEventSpecificationFieldRules))
	require.Len(t, resourceData.Get(CustomEventSpecificationFieldRules).([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, resourceData.Get(CustomEventSpecificationFieldRules).([]interface{})[0])

	rules := resourceData.Get(CustomEventSpecificationFieldRules).([]interface{})[0].(map[string]interface{})
	require.Len(t, rules, 3)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldEntityVerificationRule])
	require.Len(t, rules[CustomEventSpecificationFieldEntityVerificationRule].([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, rules[CustomEventSpecificationFieldEntityVerificationRule].([]interface{})[0])
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldSystemRule])
	require.Len(t, rules[CustomEventSpecificationFieldSystemRule].([]interface{}), 0)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldThresholdRule])
	require.Len(t, rules[CustomEventSpecificationFieldThresholdRule].([]interface{}), 0)

	rule := rules[CustomEventSpecificationFieldEntityVerificationRule].([]interface{})[0].(map[string]interface{})
	require.Len(t, rule, 5)
	require.Equal(t, restapi.SeverityWarning.GetTerraformRepresentation(), rule[CustomEventSpecificationRuleFieldSeverity])
	require.Equal(t, matchingEntityType, rule[CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityType])
	require.Equal(t, matchingEntityLabel, rule[CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityLabel])
	require.Equal(t, matchingOperator, rule[CustomEventSpecificationEntityVerificationRuleFieldMatchingOperator])
	require.Equal(t, offlineDuration, rule[CustomEventSpecificationEntityVerificationRuleFieldOfflineDuration])
}

func (r *customerEventSpecificationUnitTest) shouldMapSystemRuleToState(t *testing.T) {
	description := customEventSpecificationWithThresholdRuleDescription
	expirationTime := customEventSpecificationWithThresholdRuleExpirationTime
	query := customEventSpecificationWithThresholdRuleQuery
	systemRuleId := "system-rule-id"

	spec := &restapi.CustomEventSpecification{
		ID:             customEventSpecificationWithThresholdRuleID,
		Name:           resourceName,
		EntityType:     systemRuleEntityType,
		Query:          &query,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Triggering:     true,
		Enabled:        true,
		Rules: []restapi.RuleSpecification{
			{
				DType:        restapi.SystemRuleType,
				Severity:     restapi.SeverityWarning.GetAPIRepresentation(),
				SystemRuleID: &systemRuleId,
			},
		},
	}

	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	sut := NewCustomEventSpecificationResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec)

	require.Nil(t, err)
	require.Equal(t, customEventSpecificationWithThresholdRuleID, resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(CustomEventSpecificationFieldName))
	require.Equal(t, systemRuleEntityType, resourceData.Get(CustomEventSpecificationFieldEntityType))
	require.Equal(t, customEventSpecificationWithThresholdRuleQuery, resourceData.Get(CustomEventSpecificationFieldQuery))
	require.Equal(t, description, resourceData.Get(CustomEventSpecificationFieldDescription))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldTriggering).(bool))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldEnabled).(bool))

	require.IsType(t, []interface{}{}, resourceData.Get(CustomEventSpecificationFieldRules))
	require.Len(t, resourceData.Get(CustomEventSpecificationFieldRules).([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, resourceData.Get(CustomEventSpecificationFieldRules).([]interface{})[0])

	rules := resourceData.Get(CustomEventSpecificationFieldRules).([]interface{})[0].(map[string]interface{})
	require.Len(t, rules, 3)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldEntityVerificationRule])
	require.Len(t, rules[CustomEventSpecificationFieldEntityVerificationRule].([]interface{}), 0)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldSystemRule])
	require.Len(t, rules[CustomEventSpecificationFieldSystemRule].([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, rules[CustomEventSpecificationFieldSystemRule].([]interface{})[0])
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldThresholdRule])
	require.Len(t, rules[CustomEventSpecificationFieldThresholdRule].([]interface{}), 0)

	rule := rules[CustomEventSpecificationFieldSystemRule].([]interface{})[0].(map[string]interface{})
	require.Len(t, rule, 2)
	require.Equal(t, restapi.SeverityWarning.GetTerraformRepresentation(), rule[CustomEventSpecificationRuleFieldSeverity])
	require.Equal(t, systemRuleId, rule[CustomEventSpecificationSystemRuleFieldSystemRuleId])
}

func (r *customerEventSpecificationUnitTest) shouldMapThresholdRuleAndMetricNameToState(t *testing.T) {
	description := customEventSpecificationWithThresholdRuleDescription
	expirationTime := customEventSpecificationWithThresholdRuleExpirationTime
	query := customEventSpecificationWithThresholdRuleQuery

	window := customEventSpecificationWithThresholdRuleWindow
	rollup := customEventSpecificationWithThresholdRuleRollup
	aggregation := customEventSpecificationWithThresholdRuleAggregation
	conditionValue := customEventSpecificationWithThresholdRuleConditionValue
	metricName := customEventSpecificationWithThresholdRuleMetricName
	conditionOperator := "="

	spec := &restapi.CustomEventSpecification{
		ID:             customEventSpecificationWithThresholdRuleID,
		Name:           resourceName,
		EntityType:     customEventSpecificationWithThresholdRuleEntityType,
		Query:          &query,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Triggering:     true,
		Enabled:        true,
		Rules: []restapi.RuleSpecification{
			{
				DType:             restapi.ThresholdRuleType,
				Severity:          restapi.SeverityWarning.GetAPIRepresentation(),
				MetricName:        &metricName,
				Window:            &window,
				Rollup:            &rollup,
				Aggregation:       &aggregation,
				ConditionOperator: &conditionOperator,
				ConditionValue:    &conditionValue,
			},
		},
	}

	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	sut := NewCustomEventSpecificationResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec)

	require.Nil(t, err)
	require.Equal(t, customEventSpecificationWithThresholdRuleID, resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(CustomEventSpecificationFieldName))
	require.Equal(t, customEventSpecificationWithThresholdRuleEntityType, resourceData.Get(CustomEventSpecificationFieldEntityType))
	require.Equal(t, customEventSpecificationWithThresholdRuleQuery, resourceData.Get(CustomEventSpecificationFieldQuery))
	require.Equal(t, description, resourceData.Get(CustomEventSpecificationFieldDescription))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldTriggering).(bool))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldEnabled).(bool))

	require.IsType(t, []interface{}{}, resourceData.Get(CustomEventSpecificationFieldRules))
	require.Len(t, resourceData.Get(CustomEventSpecificationFieldRules).([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, resourceData.Get(CustomEventSpecificationFieldRules).([]interface{})[0])

	rules := resourceData.Get(CustomEventSpecificationFieldRules).([]interface{})[0].(map[string]interface{})
	require.Len(t, rules, 3)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldEntityVerificationRule])
	require.Len(t, rules[CustomEventSpecificationFieldEntityVerificationRule].([]interface{}), 0)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldSystemRule])
	require.Len(t, rules[CustomEventSpecificationFieldSystemRule].([]interface{}), 0)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldThresholdRule])
	require.Len(t, rules[CustomEventSpecificationFieldThresholdRule].([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, rules[CustomEventSpecificationFieldThresholdRule].([]interface{})[0])

	rule := rules[CustomEventSpecificationFieldThresholdRule].([]interface{})[0].(map[string]interface{})
	require.Len(t, rule, 8)
	require.Equal(t, metricName, rule[CustomEventSpecificationThresholdRuleFieldMetricName])
	require.Equal(t, []interface{}{}, rule[CustomEventSpecificationThresholdRuleFieldMetricPattern])
	require.Equal(t, window, rule[CustomEventSpecificationThresholdRuleFieldWindow])
	require.Equal(t, rollup, rule[CustomEventSpecificationThresholdRuleFieldRollup])
	require.Equal(t, string(aggregation), rule[CustomEventSpecificationThresholdRuleFieldAggregation])
	require.Equal(t, conditionOperator, rule[CustomEventSpecificationThresholdRuleFieldConditionOperator])
	require.Equal(t, conditionValue, rule[CustomEventSpecificationThresholdRuleFieldConditionValue])
	require.Equal(t, restapi.SeverityWarning.GetTerraformRepresentation(), rule[CustomEventSpecificationRuleFieldSeverity])
}

func (r *customerEventSpecificationUnitTest) shouldMapThresholdRuleAndMetricPatternToState(t *testing.T) {
	description := customEventSpecificationWithThresholdRuleDescription
	expirationTime := customEventSpecificationWithThresholdRuleExpirationTime
	query := customEventSpecificationWithThresholdRuleQuery

	window := customEventSpecificationWithThresholdRuleWindow
	rollup := customEventSpecificationWithThresholdRuleRollup
	aggregation := customEventSpecificationWithThresholdRuleAggregation
	conditionValue := customEventSpecificationWithThresholdRuleConditionValue
	conditionOperator := "="
	prefix := "prefix"
	postfix := "postfix"
	placeholder := "placeholder"
	operator := "startsWith"

	spec := &restapi.CustomEventSpecification{
		ID:             customEventSpecificationWithThresholdRuleID,
		Name:           resourceName,
		EntityType:     customEventSpecificationWithThresholdRuleEntityType,
		Query:          &query,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Triggering:     true,
		Enabled:        true,
		Rules: []restapi.RuleSpecification{
			{
				DType:    restapi.ThresholdRuleType,
				Severity: restapi.SeverityWarning.GetAPIRepresentation(),
				MetricPattern: &restapi.MetricPattern{
					Prefix:      prefix,
					Postfix:     &postfix,
					Placeholder: &placeholder,
					Operator:    operator,
				},
				Window:            &window,
				Rollup:            &rollup,
				Aggregation:       &aggregation,
				ConditionOperator: &conditionOperator,
				ConditionValue:    &conditionValue,
			},
		},
	}

	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	sut := NewCustomEventSpecificationResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec)

	require.Nil(t, err)
	require.Equal(t, customEventSpecificationWithThresholdRuleID, resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(CustomEventSpecificationFieldName))
	require.Equal(t, customEventSpecificationWithThresholdRuleEntityType, resourceData.Get(CustomEventSpecificationFieldEntityType))
	require.Equal(t, customEventSpecificationWithThresholdRuleQuery, resourceData.Get(CustomEventSpecificationFieldQuery))
	require.Equal(t, description, resourceData.Get(CustomEventSpecificationFieldDescription))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldTriggering).(bool))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldEnabled).(bool))

	require.IsType(t, []interface{}{}, resourceData.Get(CustomEventSpecificationFieldRules))
	require.Len(t, resourceData.Get(CustomEventSpecificationFieldRules).([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, resourceData.Get(CustomEventSpecificationFieldRules).([]interface{})[0])

	rules := resourceData.Get(CustomEventSpecificationFieldRules).([]interface{})[0].(map[string]interface{})
	require.Len(t, rules, 3)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldEntityVerificationRule])
	require.Len(t, rules[CustomEventSpecificationFieldEntityVerificationRule].([]interface{}), 0)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldSystemRule])
	require.Len(t, rules[CustomEventSpecificationFieldSystemRule].([]interface{}), 0)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldThresholdRule])
	require.Len(t, rules[CustomEventSpecificationFieldThresholdRule].([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, rules[CustomEventSpecificationFieldThresholdRule].([]interface{})[0])

	rule := rules[CustomEventSpecificationFieldThresholdRule].([]interface{})[0].(map[string]interface{})
	require.Len(t, rule, 8)
	require.Equal(t, "", rule[CustomEventSpecificationThresholdRuleFieldMetricName])
	require.Equal(t, window, rule[CustomEventSpecificationThresholdRuleFieldWindow])
	require.Equal(t, rollup, rule[CustomEventSpecificationThresholdRuleFieldRollup])
	require.Equal(t, string(aggregation), rule[CustomEventSpecificationThresholdRuleFieldAggregation])
	require.Equal(t, conditionOperator, rule[CustomEventSpecificationThresholdRuleFieldConditionOperator])
	require.Equal(t, conditionValue, rule[CustomEventSpecificationThresholdRuleFieldConditionValue])
	require.Equal(t, restapi.SeverityWarning.GetTerraformRepresentation(), rule[CustomEventSpecificationRuleFieldSeverity])

	require.IsType(t, []interface{}{}, rule[CustomEventSpecificationThresholdRuleFieldMetricPattern])
	require.Len(t, rule[CustomEventSpecificationThresholdRuleFieldMetricPattern].([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, rule[CustomEventSpecificationThresholdRuleFieldMetricPattern].([]interface{})[0])

	metricPatternData := rule[CustomEventSpecificationThresholdRuleFieldMetricPattern].([]interface{})[0].(map[string]interface{})

	require.Len(t, metricPatternData, 4)
	require.Equal(t, prefix, metricPatternData[CustomEventSpecificationThresholdRuleFieldMetricPatternPrefix])
	require.Equal(t, postfix, metricPatternData[CustomEventSpecificationThresholdRuleFieldMetricPatternPostfix])
	require.Equal(t, placeholder, metricPatternData[CustomEventSpecificationThresholdRuleFieldMetricPatternPlaceholder])
	require.Equal(t, string(operator), metricPatternData[CustomEventSpecificationThresholdRuleFieldMetricPatternOperator])
}

func (r *customerEventSpecificationUnitTest) shouldFailToMapRuleWhenSeverityIsNotValid(t *testing.T) {
	spec := &restapi.CustomEventSpecification{
		Rules: []restapi.RuleSpecification{
			{
				DType:    restapi.ThresholdRuleType,
				Severity: 123,
			},
		},
	}

	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	sut := NewCustomEventSpecificationResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "is not a valid severity")
}

func (r *customerEventSpecificationUnitTest) shouldFailToMapRuleWhenRuleTypeIsNotValid(t *testing.T) {
	spec := &restapi.CustomEventSpecification{
		Rules: []restapi.RuleSpecification{
			{
				DType:    "invalid",
				Severity: restapi.SeverityWarning.GetAPIRepresentation(),
			},
		},
	}

	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	sut := NewCustomEventSpecificationResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "unsupported rule type invalid")
}

func (r *customerEventSpecificationUnitTest) shouldMapStateOfEntityVerificationRuleToDataModel(t *testing.T) {
	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	resourceHandle := NewCustomEventSpecificationResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	matchingEntityLabel := "matching-entity-label"
	matchingEntityType := "matching-entity-type"
	matchingOperator := "is"
	offlineDuration := 1234

	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldName, resourceName)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldEntityType, customEventSpecificationWithThresholdRuleEntityType)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldQuery, customEventSpecificationWithThresholdRuleQuery)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldTriggering, true)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldDescription, customEventSpecificationWithThresholdRuleDescription)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldExpirationTime, customEventSpecificationWithThresholdRuleExpirationTime)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldEnabled, true)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldRules, []interface{}{
		map[string]interface{}{
			CustomEventSpecificationFieldEntityVerificationRule: []interface{}{
				map[string]interface{}{
					CustomEventSpecificationRuleFieldSeverity:                              restapi.SeverityWarning.GetTerraformRepresentation(),
					CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityLabel: matchingEntityLabel,
					CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityType:  matchingEntityType,
					CustomEventSpecificationEntityVerificationRuleFieldMatchingOperator:    matchingOperator,
					CustomEventSpecificationEntityVerificationRuleFieldOfflineDuration:     offlineDuration,
				}},
			CustomEventSpecificationFieldSystemRule:    []interface{}{},
			CustomEventSpecificationFieldThresholdRule: []interface{}{},
		},
	})

	customEventSpec, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	require.Equal(t, customEventSpecificationWithThresholdRuleID, customEventSpec.GetIDForResourcePath())
	require.Equal(t, resourceName, customEventSpec.Name)
	require.Equal(t, customEventSpecificationWithThresholdRuleEntityType, customEventSpec.EntityType)
	require.Equal(t, customEventSpecificationWithThresholdRuleQuery, *customEventSpec.Query)
	require.Equal(t, customEventSpecificationWithThresholdRuleDescription, *customEventSpec.Description)
	require.Equal(t, customEventSpecificationWithThresholdRuleExpirationTime, *customEventSpec.ExpirationTime)
	require.True(t, customEventSpec.Triggering)
	require.True(t, customEventSpec.Enabled)

	require.Equal(t, 1, len(customEventSpec.Rules))
	require.Equal(t, restapi.EntityVerificationRuleType, customEventSpec.Rules[0].DType)
	require.Equal(t, matchingEntityLabel, *customEventSpec.Rules[0].MatchingEntityLabel)
	require.Equal(t, matchingEntityType, *customEventSpec.Rules[0].MatchingEntityType)
	require.Equal(t, matchingOperator, *customEventSpec.Rules[0].MatchingOperator)
	require.Equal(t, offlineDuration, *customEventSpec.Rules[0].OfflineDuration)
}

func (r *customerEventSpecificationUnitTest) shouldFailToMapStateOfEntityVerificationRuleToDataModelWhenSeverityIsNotValid(t *testing.T) {
	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	resourceHandle := NewCustomEventSpecificationResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldRules, []interface{}{
		map[string]interface{}{
			CustomEventSpecificationFieldEntityVerificationRule: []interface{}{
				map[string]interface{}{
					CustomEventSpecificationRuleFieldSeverity: "invalid",
				}},
			CustomEventSpecificationFieldSystemRule:    []interface{}{},
			CustomEventSpecificationFieldThresholdRule: []interface{}{},
		},
	})

	_, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Error(t, err)
	require.ErrorContains(t, err, "invalid is not a valid severity")
}

func (r *customerEventSpecificationUnitTest) shouldMapStateOfSystemRuleToDataModel(t *testing.T) {
	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	resourceHandle := NewCustomEventSpecificationResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	systemRuleId := "system-rule-id"

	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldName, resourceName)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldEntityType, customEventSpecificationWithThresholdRuleEntityType)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldQuery, customEventSpecificationWithThresholdRuleQuery)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldTriggering, true)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldDescription, customEventSpecificationWithThresholdRuleDescription)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldExpirationTime, customEventSpecificationWithThresholdRuleExpirationTime)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldEnabled, true)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldRules, []interface{}{
		map[string]interface{}{
			CustomEventSpecificationFieldEntityVerificationRule: []interface{}{},
			CustomEventSpecificationFieldSystemRule: []interface{}{
				map[string]interface{}{
					CustomEventSpecificationRuleFieldSeverity:           restapi.SeverityWarning.GetTerraformRepresentation(),
					CustomEventSpecificationSystemRuleFieldSystemRuleId: systemRuleId,
				}},
			CustomEventSpecificationFieldThresholdRule: []interface{}{},
		},
	})

	customEventSpec, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	require.Equal(t, customEventSpecificationWithThresholdRuleID, customEventSpec.GetIDForResourcePath())
	require.Equal(t, resourceName, customEventSpec.Name)
	require.Equal(t, customEventSpecificationWithThresholdRuleEntityType, customEventSpec.EntityType)
	require.Equal(t, customEventSpecificationWithThresholdRuleQuery, *customEventSpec.Query)
	require.Equal(t, customEventSpecificationWithThresholdRuleDescription, *customEventSpec.Description)
	require.Equal(t, customEventSpecificationWithThresholdRuleExpirationTime, *customEventSpec.ExpirationTime)
	require.True(t, customEventSpec.Triggering)
	require.True(t, customEventSpec.Enabled)

	require.Equal(t, 1, len(customEventSpec.Rules))
	require.Equal(t, restapi.SystemRuleType, customEventSpec.Rules[0].DType)
	require.Equal(t, systemRuleId, *customEventSpec.Rules[0].SystemRuleID)
}

func (r *customerEventSpecificationUnitTest) shouldFailToMapStateOfSystemRuleToDataModelWhenSeverityIsNotValid(t *testing.T) {
	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	resourceHandle := NewCustomEventSpecificationResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldRules, []interface{}{
		map[string]interface{}{
			CustomEventSpecificationFieldEntityVerificationRule: []interface{}{},
			CustomEventSpecificationFieldSystemRule: []interface{}{
				map[string]interface{}{
					CustomEventSpecificationRuleFieldSeverity: "invalid",
				}},
			CustomEventSpecificationFieldThresholdRule: []interface{}{},
		},
	})

	_, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Error(t, err)
	require.ErrorContains(t, err, "invalid is not a valid severity")
}

func (r *customerEventSpecificationUnitTest) shouldMapStateOfThresholdRuleWithMetricNameToDataModel(t *testing.T) {
	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	resourceHandle := NewCustomEventSpecificationResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldName, resourceName)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldEntityType, customEventSpecificationWithThresholdRuleEntityType)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldQuery, customEventSpecificationWithThresholdRuleQuery)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldTriggering, true)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldDescription, customEventSpecificationWithThresholdRuleDescription)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldExpirationTime, customEventSpecificationWithThresholdRuleExpirationTime)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldEnabled, true)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldRules, []interface{}{
		map[string]interface{}{
			CustomEventSpecificationFieldEntityVerificationRule: []interface{}{},
			CustomEventSpecificationFieldSystemRule:             []interface{}{},
			CustomEventSpecificationFieldThresholdRule: []interface{}{
				map[string]interface{}{
					CustomEventSpecificationRuleFieldSeverity:                   restapi.SeverityWarning.GetTerraformRepresentation(),
					CustomEventSpecificationThresholdRuleFieldMetricName:        customEventSpecificationWithThresholdRuleMetricName,
					CustomEventSpecificationThresholdRuleFieldMetricPattern:     []interface{}{},
					CustomEventSpecificationThresholdRuleFieldRollup:            customEventSpecificationWithThresholdRuleRollup,
					CustomEventSpecificationThresholdRuleFieldWindow:            customEventSpecificationWithThresholdRuleWindow,
					CustomEventSpecificationThresholdRuleFieldAggregation:       customEventSpecificationWithThresholdRuleAggregation,
					CustomEventSpecificationThresholdRuleFieldConditionOperator: "=",
					CustomEventSpecificationThresholdRuleFieldConditionValue:    customEventSpecificationWithThresholdRuleConditionValue,
				},
			},
		},
	})

	customEventSpec, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	require.Equal(t, customEventSpecificationWithThresholdRuleID, customEventSpec.GetIDForResourcePath())
	require.Equal(t, resourceName, customEventSpec.Name)
	require.Equal(t, customEventSpecificationWithThresholdRuleEntityType, customEventSpec.EntityType)
	require.Equal(t, customEventSpecificationWithThresholdRuleQuery, *customEventSpec.Query)
	require.Equal(t, customEventSpecificationWithThresholdRuleDescription, *customEventSpec.Description)
	require.Equal(t, customEventSpecificationWithThresholdRuleExpirationTime, *customEventSpec.ExpirationTime)
	require.True(t, customEventSpec.Triggering)
	require.True(t, customEventSpec.Enabled)

	require.Equal(t, 1, len(customEventSpec.Rules))
	require.Equal(t, restapi.ThresholdRuleType, customEventSpec.Rules[0].DType)
	require.Equal(t, customEventSpecificationWithThresholdRuleMetricName, *customEventSpec.Rules[0].MetricName)
	require.Nil(t, customEventSpec.Rules[0].MetricPattern)
	require.Equal(t, customEventSpecificationWithThresholdRuleWindow, *customEventSpec.Rules[0].Window)
	require.Equal(t, customEventSpecificationWithThresholdRuleRollup, *customEventSpec.Rules[0].Rollup)
	require.Equal(t, customEventSpecificationWithThresholdRuleAggregation, *customEventSpec.Rules[0].Aggregation)
	require.Equal(t, "=", *customEventSpec.Rules[0].ConditionOperator)
	require.Equal(t, customEventSpecificationWithThresholdRuleConditionValue, *customEventSpec.Rules[0].ConditionValue)
	require.Equal(t, restapi.SeverityWarning.GetAPIRepresentation(), customEventSpec.Rules[0].Severity)
}

func (r *customerEventSpecificationUnitTest) shouldMapStateOfThresholdRuleWithMetricPatternToDataModel(t *testing.T) {
	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	resourceHandle := NewCustomEventSpecificationResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	prefix := "prefix"
	postfix := "postfix"
	placeholder := "placeholder"
	operator := "startsWith"

	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldName, resourceName)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldEntityType, customEventSpecificationWithThresholdRuleEntityType)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldQuery, customEventSpecificationWithThresholdRuleQuery)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldTriggering, true)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldDescription, customEventSpecificationWithThresholdRuleDescription)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldExpirationTime, customEventSpecificationWithThresholdRuleExpirationTime)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldEnabled, true)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldRules, []interface{}{
		map[string]interface{}{
			CustomEventSpecificationFieldEntityVerificationRule: []interface{}{},
			CustomEventSpecificationFieldSystemRule:             []interface{}{},
			CustomEventSpecificationFieldThresholdRule: []interface{}{
				map[string]interface{}{
					CustomEventSpecificationRuleFieldSeverity: restapi.SeverityWarning.GetTerraformRepresentation(),
					CustomEventSpecificationThresholdRuleFieldMetricPattern: []interface{}{
						map[string]interface{}{
							CustomEventSpecificationThresholdRuleFieldMetricPatternPrefix:      prefix,
							CustomEventSpecificationThresholdRuleFieldMetricPatternPostfix:     postfix,
							CustomEventSpecificationThresholdRuleFieldMetricPatternPlaceholder: placeholder,
							CustomEventSpecificationThresholdRuleFieldMetricPatternOperator:    operator,
						},
					},
					CustomEventSpecificationThresholdRuleFieldRollup:            customEventSpecificationWithThresholdRuleRollup,
					CustomEventSpecificationThresholdRuleFieldWindow:            customEventSpecificationWithThresholdRuleWindow,
					CustomEventSpecificationThresholdRuleFieldAggregation:       customEventSpecificationWithThresholdRuleAggregation,
					CustomEventSpecificationThresholdRuleFieldConditionOperator: "=",
					CustomEventSpecificationThresholdRuleFieldConditionValue:    customEventSpecificationWithThresholdRuleConditionValue,
				},
			},
		},
	})

	customEventSpec, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	require.Equal(t, customEventSpecificationWithThresholdRuleID, customEventSpec.GetIDForResourcePath())
	require.Equal(t, resourceName, customEventSpec.Name)
	require.Equal(t, customEventSpecificationWithThresholdRuleEntityType, customEventSpec.EntityType)
	require.Equal(t, customEventSpecificationWithThresholdRuleQuery, *customEventSpec.Query)
	require.Equal(t, customEventSpecificationWithThresholdRuleDescription, *customEventSpec.Description)
	require.Equal(t, customEventSpecificationWithThresholdRuleExpirationTime, *customEventSpec.ExpirationTime)
	require.True(t, customEventSpec.Triggering)
	require.True(t, customEventSpec.Enabled)

	require.Equal(t, 1, len(customEventSpec.Rules))
	require.Equal(t, restapi.ThresholdRuleType, customEventSpec.Rules[0].DType)
	require.Nil(t, customEventSpec.Rules[0].MetricName)
	require.NotNil(t, customEventSpec.Rules[0].MetricPattern)
	require.Equal(t, prefix, customEventSpec.Rules[0].MetricPattern.Prefix)
	require.Equal(t, postfix, *customEventSpec.Rules[0].MetricPattern.Postfix)
	require.Equal(t, placeholder, *customEventSpec.Rules[0].MetricPattern.Placeholder)
	require.Equal(t, operator, customEventSpec.Rules[0].MetricPattern.Operator)
	require.Equal(t, customEventSpecificationWithThresholdRuleWindow, *customEventSpec.Rules[0].Window)
	require.Equal(t, customEventSpecificationWithThresholdRuleRollup, *customEventSpec.Rules[0].Rollup)
	require.Equal(t, customEventSpecificationWithThresholdRuleAggregation, *customEventSpec.Rules[0].Aggregation)
	require.Equal(t, "=", *customEventSpec.Rules[0].ConditionOperator)
	require.Equal(t, customEventSpecificationWithThresholdRuleConditionValue, *customEventSpec.Rules[0].ConditionValue)
	require.Equal(t, restapi.SeverityWarning.GetAPIRepresentation(), customEventSpec.Rules[0].Severity)
}

func (r *customerEventSpecificationUnitTest) shouldFailToMapStateOfThresholdRuleToDataModelWhenSeverityIsNotValid(t *testing.T) {
	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	resourceHandle := NewCustomEventSpecificationResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldRules, []interface{}{
		map[string]interface{}{
			CustomEventSpecificationFieldEntityVerificationRule: []interface{}{},
			CustomEventSpecificationFieldSystemRule:             []interface{}{},
			CustomEventSpecificationFieldThresholdRule: []interface{}{
				map[string]interface{}{
					CustomEventSpecificationRuleFieldSeverity: "invalid",
				},
			},
		},
	})

	_, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Error(t, err)
	require.ErrorContains(t, err, "invalid is not a valid severity")
}

func (r *customerEventSpecificationUnitTest) shouldFailToMapStateWhenNoRuleIsProvided(t *testing.T) {
	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	resourceHandle := NewCustomEventSpecificationResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldRules, []interface{}{
		map[string]interface{}{
			CustomEventSpecificationFieldEntityVerificationRule: []interface{}{},
			CustomEventSpecificationFieldSystemRule:             []interface{}{},
			CustomEventSpecificationFieldThresholdRule:          []interface{}{},
		},
	})

	_, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Error(t, err)
	require.ErrorContains(t, err, "no supported rule defined")
}
