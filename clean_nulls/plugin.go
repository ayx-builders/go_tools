package clean_nulls

import (
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/output_connection"
	"github.com/tlarsen7572/goalteryx/presort"
)

type Plugin struct {
	toolId int
	output output_connection.OutputConnection
}

func (p *Plugin) Init(toolId int, _ string) bool {
	p.toolId = toolId
	p.output = output_connection.New(toolId, `Output`, 10)
	return true
}

func (p *Plugin) PushAllRecords(_ int) bool {
	panic("not valid: this is not an input tool")
}

func (p *Plugin) Close(_ bool) {}

func (p *Plugin) AddIncomingConnection(_ string, _ string) (api.IncomingInterface, *presort.PresortInfo) {
	return &Ii{toolId: p.toolId, output: p.output}, nil
}

func (p *Plugin) AddOutgoingConnection(_ string, connectionInterface *api.ConnectionInterfaceStruct) bool {
	p.output.Add(connectionInterface)
	return true
}

func (p *Plugin) GetToolId() int {
	return p.toolId
}
