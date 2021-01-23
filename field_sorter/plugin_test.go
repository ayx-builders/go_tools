package field_sorter_test

import (
	"github.com/ayx-builders/go_tools/field_sorter"
	"github.com/tlarsen7572/goalteryx/sdk"
	"reflect"
	"testing"
)

func TestSort(t *testing.T) {
	plugin := &field_sorter.Plugin{}
	config := `<Configuration>
  <alphabetical>false</alphabetical>
  <field0>
    <text>Field3</text>
    <isPattern>false</isPattern>
  </field0>
  <field1>
    <text>Field.*</text>
    <isPattern>true</isPattern>
  </field1>
</Configuration>`
	runner := sdk.RegisterToolTest(plugin, 0, config)
	output := runner.CaptureOutgoingAnchor(`Output`)
	runner.ConnectInput(`Input`, `testfile.txt`)
	runner.SimulateLifecycle()
	metadata := output.Input.Metadata()

	expectedFields := []string{
		"Field3",
		"Field1",
		"Field2",
		"Field4",
		"Field5",
		"Field6",
	}
	actualFields := make([]string, len(metadata.Fields()))
	for index, actualField := range metadata.Fields() {
		actualFields[index] = actualField.Name
	}
	if !reflect.DeepEqual(expectedFields, actualFields) {
		t.Fatalf("expected\n%v\nbut got\n%v", expectedFields, actualFields)
	}
}
