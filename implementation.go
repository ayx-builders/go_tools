package main

/*
#include "implementation.h"
*/
import "C"
import (
	"github.com/ayx-builders/go_tools/clean_nulls"
	"github.com/tlarsen7572/goalteryx/api"
	"unsafe"
)

func main() {}

//export CleanNulls
func CleanNulls(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &clean_nulls.Plugin{}
	return C.long(api.ConfigurePlugin(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}
