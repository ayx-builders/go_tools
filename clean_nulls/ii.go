package clean_nulls

import (
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/output_connection"
	"github.com/tlarsen7572/goalteryx/recordblob"
	"github.com/tlarsen7572/goalteryx/recordcopier"
	"github.com/tlarsen7572/goalteryx/recordinfo"
)

type Ii struct {
	toolId int
	output output_connection.OutputConnection
	info   recordinfo.RecordInfo
	copier *recordcopier.RecordCopier
	fields int
}

func (i *Ii) Init(recordInfoIn string) bool {
	var err error
	i.info, err = recordinfo.FromXml(recordInfoIn)
	if err != nil {
		api.OutputMessage(i.toolId, api.Error, err.Error())
		return false
	}
	err = i.output.Init(i.info)
	if err != nil {
		api.OutputMessage(i.toolId, api.Error, err.Error())
		return false
	}

	i.fields = i.info.NumFields()
	copierMap := make([]recordcopier.IndexMap, i.fields)
	for index := 0; index < i.fields; index++ {
		copierMap[index] = recordcopier.IndexMap{
			DestinationIndex: index,
			SourceIndex:      index,
		}
	}
	i.copier, err = recordcopier.New(i.info, i.info, copierMap)
	if err != nil {
		api.OutputMessage(i.toolId, api.Error, err.Error())
		return false
	}
	return true
}

func (i *Ii) PushRecord(record recordblob.RecordBlob) bool {
	err := i.copier.Copy(record)
	if err != nil {
		api.OutputMessage(i.toolId, api.Error, err.Error())
		return false
	}

	for index := 0; index < i.fields; index++ {
		field, _ := i.info.GetFieldByIndex(index)
		switch field.Type {
		case recordinfo.Byte, recordinfo.Int16, recordinfo.Int32, recordinfo.Int64:
			_, isNull, _ := i.info.GetCurrentInt(field.Name)
			if isNull {
				_ = i.info.SetIntField(field.Name, 0)
			}
		case recordinfo.Float, recordinfo.Double, recordinfo.FixedDecimal:
			_, isNull, _ := i.info.GetCurrentFloat(field.Name)
			if isNull {
				_ = i.info.SetFloatField(field.Name, 0)
			}
		case recordinfo.String, recordinfo.WString, recordinfo.V_WString, recordinfo.V_String:
			_, isNull, _ := i.info.GetCurrentString(field.Name)
			if isNull {
				_ = i.info.SetStringField(field.Name, ``)
			}
		case recordinfo.Bool:
			_, isNull, _ := i.info.GetCurrentBool(field.Name)
			if isNull {
				_ = i.info.SetBoolField(field.Name, false)
			}
		default:
			continue
		}
	}
	blob, err := i.info.GenerateRecord()
	if err != nil {
		api.OutputMessage(i.toolId, api.Error, err.Error())
		return false
	}
	i.output.PushRecord(blob)
	return true
}

func (i *Ii) UpdateProgress(percent float64) {
	api.OutputToolProgress(i.toolId, percent)
	i.output.UpdateProgress(percent)
}

func (i *Ii) Close() {
	i.output.Close()
}

func (i *Ii) CacheSize() int {
	return 10
}
