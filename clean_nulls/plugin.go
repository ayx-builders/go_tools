package clean_nulls

import (
	"github.com/tlarsen7572/goalteryx/sdk"
)

type Plugin struct {
	output sdk.OutputAnchor
	info   *sdk.OutgoingRecordInfo
}

func (p *Plugin) Init(provider sdk.Provider) {
	p.output = provider.GetOutputAnchor(`Output`)
}

func (p *Plugin) OnInputConnectionOpened(connection sdk.InputConnection) {
	p.info = connection.Metadata().Clone().GenerateOutgoingRecordInfo()
	p.output.Open(p.info)
}

func (p *Plugin) OnRecordPacket(connection sdk.InputConnection) {
	packet := connection.Read()
	for packet.Next() {
		p.info.CopyFrom(packet.Record())
		for _, field := range p.info.StringFields {
			if _, isNull := field.GetCurrentString(); isNull {
				field.SetString(``)
			}
		}
		for _, field := range p.info.IntFields {
			if _, isNull := field.GetCurrentInt(); isNull {
				field.SetInt(0)
			}
		}
		for _, field := range p.info.BoolFields {
			if _, isNull := field.GetCurrentBool(); isNull {
				field.SetBool(false)
			}
		}
		for _, field := range p.info.FloatFields {
			if _, isNull := field.GetCurrentFloat(); isNull {
				field.SetFloat(0)
			}
		}
		p.output.Write()
	}
}

func (p *Plugin) OnComplete() {}
