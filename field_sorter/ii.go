package field_sorter

import (
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/output_connection"
	"github.com/tlarsen7572/goalteryx/recordblob"
	"github.com/tlarsen7572/goalteryx/recordcopier"
	"github.com/tlarsen7572/goalteryx/recordinfo"
)

type Ii struct {
	config  *config
	toolId  int
	output  output_connection.OutputConnection
	inInfo  recordinfo.RecordInfo
	outInfo recordinfo.RecordInfo
	copier  *recordcopier.RecordCopier
}

func (i *Ii) Init(recordInfoIn string) bool {
	var err error
	i.inInfo, err = recordinfo.FromXml(recordInfoIn)
	if err != nil {
		api.OutputMessage(i.toolId, api.Error, err.Error())
		return false
	}
	fieldNames := make([]string, i.inInfo.NumFields())
	for index := 0; index < i.inInfo.NumFields(); index++ {
		fieldInfo, _ := i.inInfo.GetFieldByIndex(index)
		fieldNames[index] = fieldInfo.Name
	}
	mapping, err := sortFields(fieldNames, i.config.Fields, i.config.Alphabetical)
	if err != nil {
		api.OutputMessage(i.toolId, api.Error, err.Error())
		return false
	}
	generator := recordinfo.NewGenerator()
	for _, entry := range mapping {
		fieldInfo, _ := i.inInfo.GetFieldByIndex(entry.SourceIndex)
		generator.AddField(fieldInfo, `Field Sorter`)
	}
	i.outInfo = generator.GenerateRecordInfo()
	i.copier, err = recordcopier.New(i.outInfo, i.inInfo.GenerateRecordBlobReader(), mapping)
	if err != nil {
		api.OutputMessage(i.toolId, api.Error, err.Error())
		return false
	}
	_ = i.output.Init(i.outInfo)
	return true
}

func (i *Ii) PushRecord(record recordblob.RecordBlob) bool {
	err := i.copier.Copy(record)
	if err != nil {
		api.OutputMessage(i.toolId, api.Error, err.Error())
		return false
	}
	blob, err := i.outInfo.GenerateRecord()
	if err != nil {
		api.OutputMessage(i.toolId, api.Error, err.Error())
		return false
	}
	i.output.PushRecord(blob)
	return true
}

func (i *Ii) UpdateProgress(percent float64) {
	api.OutputToolProgress(i.toolId, percent)
}

func (i *Ii) Close() {
	i.output.Close()
}

func (i *Ii) CacheSize() int {
	return 10
}
