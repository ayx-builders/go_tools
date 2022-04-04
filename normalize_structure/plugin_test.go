package normalize_structure

import (
	"github.com/tlarsen7572/goalteryx/sdk"
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
}
