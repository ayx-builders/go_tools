package normalize_structure

import (
	"github.com/tlarsendataguy/goalteryx/sdk"
	"testing"
)

func TestNormalizeStructure(t *testing.T) {
	plugin := &NormalizeStructure{}
	runner := sdk.RegisterToolTest(plugin, 1, `<Configuration><StringType>V_WString</StringType><StringLength>200</StringLength><IntType>Int64</IntType><FloatType>Double</FloatType></Configuration>`)
	collector := runner.CaptureOutgoingAnchor(`Output`)
	runner.ConnectInput(`Input`, `test.txt`)

	runner.SimulateLifecycle()

	if fieldLen := len(collector.Data); fieldLen != 11 {
		t.Fatalf(`expected 11 fields but got %v`, fieldLen)
	}
	expected := map[string]string{
		`field1`:  `Int64`,
		`field2`:  `Int64`,
		`field3`:  `Int64`,
		`field4`:  `Int64`,
		`field5`:  `Double`,
		`field6`:  `Double`,
		`field7`:  `Double`,
		`field8`:  `V_WString`,
		`field9`:  `V_WString`,
		`field10`: `V_WString`,
		`field11`: `V_WString`,
	}
	actual := map[string]string{}
	for _, field := range collector.Config.Fields() {
		actual[field.Name] = field.Type
	}
	for name, expectedType := range expected {
		if actual[name] != expectedType {
			t.Fatalf(`expected field %v to be type %v but got %v`, name, expectedType, actual[name])
		}
	}
	totalRecords := len(collector.Data[`field1`])
	if totalRecords != 1 {
		t.Fatalf(`expected 1 record but got %v`, totalRecords)
	}
}
