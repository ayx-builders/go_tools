package main

import "C"
import (
	"github.com/ayx-builders/go_tools/clean_nulls"
	"github.com/ayx-builders/go_tools/field_sorter"
	"github.com/tlarsen7572/goalteryx/sdk"
	"unsafe"
)

func main() {}

//export CleanNulls
func CleanNulls(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &clean_nulls.Plugin{}
	return C.long(sdk.RegisterTool(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}

//export FieldSorter
func FieldSorter(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &field_sorter.Plugin{}
	return C.long(sdk.RegisterTool(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}
