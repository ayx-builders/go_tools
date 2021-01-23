package field_sorter_test

import (
	"github.com/ayx-builders/go_tools/field_sorter"
	"reflect"
	"testing"
)

func TestAlphabeticalSort(t *testing.T) {
	fields := []string{`B`, `A`, `C`}
	sortInfo := make([]field_sorter.FieldSortInfo, 0)
	result, err := field_sorter.SortFields(fields, sortInfo, true)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	expected := []string{`A`, `B`, `C`}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf(`expected '%v' but got '%v'`, expected, result)
	}
}

func TestExactFieldNameSort(t *testing.T) {
	fields := []string{`B`, `A`, `C`}
	sortInfo := []field_sorter.FieldSortInfo{
		{Text: `C`, IsPattern: false},
		{Text: `A`, IsPattern: false},
		{Text: `B`, IsPattern: false},
	}
	result, err := field_sorter.SortFields(fields, sortInfo, false)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	expected := []string{`C`, `A`, `B`}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf(`expected '%v' but got '%v'`, expected, result)
	}
}

func TestPatternMatchSort(t *testing.T) {
	fields := []string{`AB`, `AA`, `BC`}
	sortInfo := []field_sorter.FieldSortInfo{
		{Text: `^B.*`, IsPattern: true},
		{Text: `^A.*`, IsPattern: true},
	}
	result, err := field_sorter.SortFields(fields, sortInfo, false)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	expected := []string{`BC`, `AB`, `AA`}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf(`expected '%v' but got '%v'`, expected, result)
	}
}

func TestPatternWithUnsortedFields(t *testing.T) {
	fields := []string{`AB`, `AA`, `BC`}
	sortInfo := []field_sorter.FieldSortInfo{
		{Text: `^B.*`, IsPattern: true},
	}
	result, err := field_sorter.SortFields(fields, sortInfo, true)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	expected := []string{`BC`, `AA`, `AB`}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf(`expected '%v' but got '%v'`, expected, result)
	}
}

func TestOverlappingRules(t *testing.T) {
	fields := []string{`Field1`, `Field2`, `Field3`, `Field4`, `Field5`, `Field6`}
	sortInfo := []field_sorter.FieldSortInfo{
		{Text: `Field3`, IsPattern: false},
		{Text: `Field.*`, IsPattern: true},
	}
	result, err := field_sorter.SortFields(fields, sortInfo, true)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	expected := []string{`Field3`, `Field1`, `Field2`, `Field4`, `Field5`, `Field6`}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf(`expected '%v' but got '%v'`, expected, result)
	}
}
