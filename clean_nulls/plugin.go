package clean_nulls

import (
	"github.com/tlarsendataguy/goalteryx/sdk"
)

type Plugin struct {
	output   sdk.OutputAnchor
	info     *sdk.OutgoingRecordInfo
	provider sdk.Provider
}

func (p *Plugin) Init(provider sdk.Provider) {
	p.provider = provider
	p.output = provider.GetOutputAnchor(`Output`)
}

func (p *Plugin) OnInputConnectionOpened(connection sdk.InputConnection) {
	p.info = connection.Metadata().Clone().GenerateOutgoingRecordInfo()
	p.output.Open(p.info)
	p.provider.Io().UpdateProgress(0.0)
}

func (p *Plugin) OnRecordPacket(connection sdk.InputConnection) {
	packet := connection.Read()
	for packet.Next() {
		p.info.CopyFrom(packet.Record())
		for _, field := range p.info.StringFields {
			if field.GetNull() {
				field.SetString(``)
			}
		}
		for _, field := range p.info.IntFields {
			if field.GetNull() {
				field.SetInt(0)
			}
		}
		for _, field := range p.info.BoolFields {
			if field.GetNull() {
				field.SetBool(false)
			}
		}
		for _, field := range p.info.FloatFields {
			if field.GetNull() {
				field.SetFloat(0)
			}
		}
		p.output.Write()
	}
	p.output.UpdateProgress(connection.Progress())
	p.provider.Io().UpdateProgress(connection.Progress())
}

func (p *Plugin) OnComplete() {
	p.output.UpdateProgress(1.0)
	p.provider.Io().UpdateProgress(1.0)
}
