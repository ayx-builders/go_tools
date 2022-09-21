package normalize_structure

import (
	"encoding/xml"
	"fmt"
	"github.com/tlarsendataguy/goalteryx/sdk"
	"strings"
)

type NormalizeStructureConfig struct {
	StringType   string
	StringLength int
	IntType      string
	FloatType    string
}

var stringTypes = []string{`string`, `wstring`, `v_string`, `v_wstring`}
var intTypes = []string{`byte`, `int16`, `int32`, `int64`}
var floatTypes = []string{`float`, `double`, `fixeddecimal`}

const source = `Normalize Structure`

type NormalizeStructure struct {
	hasError       bool
	outputOpened   bool
	provider       sdk.Provider
	config         NormalizeStructureConfig
	output         sdk.OutputAnchor
	outInfo        *sdk.OutgoingRecordInfo
	editor         *sdk.EditingRecordInfo
	addStringField func(name, source string, length int, options ...sdk.AddFieldOptionSetter) string
	addIntField    func(name, source string, options ...sdk.AddFieldOptionSetter) string
	addFloatField  func(name, source string, options ...sdk.AddFieldOptionSetter) string
	sourceStrings  map[string]sdk.StringGetter
	sourceInts     map[string]sdk.IntGetter
	sourceFloats   map[string]sdk.FloatGetter
	sourceDates    map[string]sdk.TimeGetter
	sourceBlob     map[string]sdk.BytesGetter
	sourceBool     map[string]sdk.BoolGetter
}

func (n *NormalizeStructure) Init(provider sdk.Provider) {
	n.sourceStrings = make(map[string]sdk.StringGetter)
	n.sourceInts = make(map[string]sdk.IntGetter)
	n.sourceFloats = make(map[string]sdk.FloatGetter)
	n.sourceDates = make(map[string]sdk.TimeGetter)
	n.sourceBlob = make(map[string]sdk.BytesGetter)
	n.sourceBool = make(map[string]sdk.BoolGetter)
	n.provider = provider
	err := xml.Unmarshal([]byte(provider.ToolConfig()), &n.config)
	if err != nil {
		n.sendError(`error decoding config: %v`, err.Error())
		return
	}
	n.output = provider.GetOutputAnchor(`Output`)
	n.editor = &sdk.EditingRecordInfo{}
	switch strings.ToLower(n.config.StringType) {
	case `string`:
		n.addStringField = n.editor.AddStringField
	case `wstring`:
		n.addStringField = n.editor.AddWStringField
	case `v_string`:
		n.addStringField = n.editor.AddV_StringField
	case `v_wstring`:
		n.addStringField = n.editor.AddV_WStringField
	default:
		n.sendError(`string format %v is not valid, choose String, WString, V_String, or V_WString`, n.config.StringType)
	}
	switch strings.ToLower(n.config.IntType) {
	case `byte`:
		n.addIntField = n.editor.AddByteField
	case `int16`:
		n.addIntField = n.editor.AddInt16Field
	case `int32`:
		n.addIntField = n.editor.AddInt32Field
	case `int64`:
		n.addIntField = n.editor.AddInt64Field
	default:
		n.sendError(`int format %v is not valid, choose Byte, Int16, Int32, or Int64`, n.config.IntType)
	}
	switch strings.ToLower(n.config.FloatType) {
	case `float`:
		n.addFloatField = n.editor.AddFloatField
	case `double`:
		n.addFloatField = n.editor.AddDoubleField
	default:
		n.sendError(`float format %v is not valid, choose Float or Double`, n.config.FloatType)
	}
}

func (n *NormalizeStructure) OnInputConnectionOpened(connection sdk.InputConnection) {
	if n.hasError {
		return
	}

	inInfo := connection.Metadata()

	for _, field := range inInfo.Fields() {
		if listContains(field.Type, stringTypes) {
			n.addStringField(field.Name, source, n.config.StringLength)
			in, _ := inInfo.GetStringField(field.Name)
			n.sourceStrings[field.Name] = in.GetValue
			continue
		}
		if listContains(field.Type, intTypes) {
			n.addIntField(field.Name, source)
			in, _ := inInfo.GetIntField(field.Name)
			n.sourceInts[field.Name] = in.GetValue
			continue
		}
		if listContains(field.Type, floatTypes) {
			n.addFloatField(field.Name, source)
			in, _ := inInfo.GetFloatField(field.Name)
			n.sourceFloats[field.Name] = in.GetValue
			continue
		}
		switch strings.ToLower(field.Type) {
		case `date`:
			n.editor.AddDateField(field.Name, source)
			in, _ := inInfo.GetTimeField(field.Name)
			n.sourceDates[field.Name] = in.GetValue
		case `datetime`:
			n.editor.AddDateTimeField(field.Name, source)
			in, _ := inInfo.GetTimeField(field.Name)
			n.sourceDates[field.Name] = in.GetValue
		case `spatialobj`, `blob`:
			n.editor.AddBlobField(field.Name, source, field.Size)
			in, _ := inInfo.GetBlobField(field.Name)
			n.sourceBlob[field.Name] = in.GetValue
		case `bool`:
			n.editor.AddBoolField(field.Name, source)
			in, _ := inInfo.GetBoolField(field.Name)
			n.sourceBool[field.Name] = in.GetValue
		default:
			n.sendError(`invalid data type for field %v: %v`, field.Name, field.Type)
		}
	}

	n.outInfo = n.editor.GenerateOutgoingRecordInfo()
	n.output.Open(n.outInfo)
	n.outputOpened = true
}

func (n *NormalizeStructure) OnRecordPacket(connection sdk.InputConnection) {
	if n.hasError {
		return
	}

	packet := connection.Read()
	for packet.Next() {
		for name, getter := range n.sourceStrings {
			value, isNull := getter(packet.Record())
			if isNull {
				n.outInfo.StringFields[name].SetNull()
			} else {
				n.outInfo.StringFields[name].SetString(value)
			}
		}
		for name, getter := range n.sourceInts {
			value, isNull := getter(packet.Record())
			if isNull {
				n.outInfo.IntFields[name].SetNull()
			} else {
				n.outInfo.IntFields[name].SetInt(value)
			}
		}
		for name, getter := range n.sourceFloats {
			value, isNull := getter(packet.Record())
			if isNull {
				n.outInfo.FloatFields[name].SetNull()
			} else {
				n.outInfo.FloatFields[name].SetFloat(value)
			}
		}
		for name, getter := range n.sourceBool {
			value, isNull := getter(packet.Record())
			if isNull {
				n.outInfo.BoolFields[name].SetNull()
			} else {
				n.outInfo.BoolFields[name].SetBool(value)
			}
		}
		for name, getter := range n.sourceDates {
			value, isNull := getter(packet.Record())
			if isNull {
				n.outInfo.DateTimeFields[name].SetNull()
			} else {
				n.outInfo.DateTimeFields[name].SetDateTime(value)
			}
		}
		for name, getter := range n.sourceBlob {
			n.outInfo.BlobFields[name].SetBlob(getter(packet.Record()))
		}
		n.output.Write()
	}
	n.output.UpdateProgress(connection.Progress())
	n.provider.Io().UpdateProgress(connection.Progress())
}

func (n *NormalizeStructure) OnComplete() {
	if n.outputOpened {
		n.output.Close()
	}
}

func (n *NormalizeStructure) sendError(format string, a ...interface{}) {
	n.hasError = true
	n.provider.Io().Error(fmt.Sprintf(format, a...))
}

func listContains(value string, values []string) bool {
	for _, item := range values {
		if strings.ToLower(value) == strings.ToLower(item) {
			return true
		}
	}
	return false
}
