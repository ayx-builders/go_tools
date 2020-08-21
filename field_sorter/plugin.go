package field_sorter

import (
	"encoding/xml"
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/output_connection"
	"github.com/tlarsen7572/goalteryx/presort"
)

type Plugin struct {
	toolId int
	config *config
	output output_connection.OutputConnection
}

type config struct {
	Alphabetical bool    `xml:"alphabetical"`
	Fields       []field `xml:",any"`
}

type field struct {
	IsPattern bool   `xml:"isPattern"`
	Text      string `xml:"text"`
}

func (p *Plugin) Init(toolId int, configStr string) bool {
	p.toolId = toolId
	p.config = &config{}
	err := xml.Unmarshal([]byte(configStr), p.config)
	if err != nil {
		api.OutputMessage(toolId, api.Error, err.Error())
		return false
	}
	p.output = output_connection.New(toolId, `Output`, 10)
	return true
}

func (p *Plugin) PushAllRecords(_ int) bool {
	panic("PushAllRecords is invalid")
}

func (p *Plugin) Close(_ bool) {}

func (p *Plugin) AddIncomingConnection(_ string, _ string) (api.IncomingInterface, *presort.PresortInfo) {
	return &Ii{
		toolId: p.toolId,
		config: p.config,
		output: p.output,
	}, nil
}

func (p *Plugin) AddOutgoingConnection(_ string, connectionInterface *api.ConnectionInterfaceStruct) bool {
	p.output.Add(connectionInterface)
	return true
}

func (p *Plugin) GetToolId() int {
	return p.toolId
}
