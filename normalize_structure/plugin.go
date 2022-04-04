package normalize_structure

import (
	"github.com/tlarsen7572/goalteryx/sdk"
)

type NormalizeStructure struct {
	provider sdk.Provider
	output   sdk.OutputAnchor
	outInfo  *sdk.OutgoingRecordInfo
}

func (n *NormalizeStructure) Init(provider sdk.Provider) {
	n.provider = provider
	n.output = provider.GetOutputAnchor(`Output`)
}

func (n *NormalizeStructure) OnInputConnectionOpened(connection sdk.InputConnection) {
	inInfo := connection.Metadata()
	n.outInfo = inInfo.Clone().GenerateOutgoingRecordInfo()
	n.output.Open(n.outInfo)
}

func (n *NormalizeStructure) OnRecordPacket(connection sdk.InputConnection) {
}

func (n *NormalizeStructure) OnComplete() {
	n.output.Close()
}
