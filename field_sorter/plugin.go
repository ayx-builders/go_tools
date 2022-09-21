package field_sorter

import (
	"encoding/xml"
	"fmt"
	"github.com/tlarsendataguy/goalteryx/sdk"
)

type config struct {
	Alphabetical bool            `xml:"alphabetical"`
	Fields       []FieldSortInfo `xml:",any"`
}

type FieldSortInfo struct {
	IsPattern bool   `xml:"isPattern"`
	Text      string `xml:"text"`
}

type Plugin struct {
	config   *config
	output   sdk.OutputAnchor
	info     *sdk.OutgoingRecordInfo
	provider sdk.Provider
}

func (p *Plugin) Init(provider sdk.Provider) {
	p.provider = provider
	p.config = &config{}
	err := xml.Unmarshal([]byte(provider.ToolConfig()), p.config)
	if err != nil {
		provider.Io().Error(err.Error())
	}
	p.output = provider.GetOutputAnchor(`Output`)
}

func (p *Plugin) OnInputConnectionOpened(connection sdk.InputConnection) {
	metadata := connection.Metadata()
	editor := metadata.Clone()

	incomingFields := metadata.Fields()
	fieldNames := make([]string, len(incomingFields))
	for index, field := range incomingFields {
		fieldNames[index] = field.Name
	}
	sortedFields, err := SortFields(fieldNames, p.config.Fields, p.config.Alphabetical)
	if err != nil {
		p.provider.Io().Error(fmt.Sprintf(`error sorting fields: %v`, err.Error()))
		return
	}
	for index, field := range sortedFields {
		err = editor.MoveField(field, index)
		if err != nil {
			p.provider.Io().Error(fmt.Sprintf(`error moving field %v to position %v: %v`, field, index, err.Error()))
			return
		}
	}
	p.info = editor.GenerateOutgoingRecordInfo()
	p.output.Open(p.info)
}

func (p *Plugin) OnRecordPacket(connection sdk.InputConnection) {
	packet := connection.Read()
	for packet.Next() {
		p.info.CopyFrom(packet.Record())
		p.output.Write()
	}
}

func (p *Plugin) OnComplete() {}
