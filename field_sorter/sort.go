package field_sorter

import (
	"regexp"
	"sort"
)

func SortFields(fields []string, sortInfo []FieldSortInfo, sortAlphabetical bool) ([]string, error) {
	sorter := &fieldSorter{
		fields:             fields,
		sortInfo:           sortInfo,
		sortAlphabetical:   sortAlphabetical,
		destinationIndex:   0,
		destinationIndices: make([]string, len(fields)),
	}
	return sorter.Sort()
}

type fieldSorter struct {
	fields             []string
	sortInfo           []FieldSortInfo
	sortAlphabetical   bool
	destinationIndex   int
	destinationIndices []string
}

func (f *fieldSorter) Sort() ([]string, error) {
	f.sortAlphabeticallyIfNeeded()
	err := f.processRules()
	if err != nil {
		return nil, err
	}
	f.appendUnsortedFields()
	return f.destinationIndices, nil
}

func (f *fieldSorter) sortAlphabeticallyIfNeeded() {
	if f.sortAlphabetical {
		sort.Strings(f.fields)
	}
}

func (f *fieldSorter) processRules() error {
	for _, sortRule := range f.sortInfo {
		if sortRule.IsPattern {
			err := f.processPatternMatch(sortRule)
			if err != nil {
				return err
			}
		} else {
			f.processExactMatch(sortRule)
		}
	}
	return nil
}

func (f *fieldSorter) processExactMatch(sortRule FieldSortInfo) {
	for sourceIndex := 0; sourceIndex < len(f.fields); sourceIndex++ {
		sourceField := f.fields[sourceIndex]
		if sourceField == sortRule.Text {
			f.markFieldDestination(sourceField)
			f.fields = removeItem(f.fields, sourceIndex)
			return
		}
	}
}

func (f *fieldSorter) processPatternMatch(sortRule FieldSortInfo) error {
	regex, err := regexp.Compile(sortRule.Text)
	if err != nil {
		return err
	}
	for sourceIndex := 0; sourceIndex < len(f.fields); sourceIndex++ {
		sourceField := f.fields[sourceIndex]
		if regex.MatchString(sourceField) {
			f.markFieldDestination(sourceField)
			f.fields = removeItem(f.fields, sourceIndex)
			sourceIndex--
		}
	}
	return nil
}

func (f *fieldSorter) appendUnsortedFields() {
	for _, field := range f.fields {
		f.markFieldDestination(field)
	}
}

func (f *fieldSorter) markFieldDestination(fieldName string) {
	f.destinationIndices[f.destinationIndex] = fieldName
	f.destinationIndex++
}

func removeItem(items []string, index int) []string {
	return append(items[:index], items[index+1:]...)
}
