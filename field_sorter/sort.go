package field_sorter

import (
	"github.com/tlarsen7572/goalteryx/recordcopier"
	"regexp"
	"sort"
)

type sourceField struct {
	name  string
	index int
}

func sortFields(fields []string, sortInfo []field, sortAlphabetical bool) ([]recordcopier.IndexMap, error) {
	sourceFields := make([]sourceField, len(fields))
	for index, field := range fields {
		sourceFields[index] = sourceField{
			name:  field,
			index: index,
		}
	}
	if sortAlphabetical {
		sort.SliceStable(sourceFields, func(i, j int) bool {
			return sourceFields[i].name < sourceFields[j].name
		})
	}
	mapping := make([]recordcopier.IndexMap, len(sourceFields))
	destinationIndex := 0
	for _, sortRule := range sortInfo {
		if len(sourceFields) == 0 {
			return mapping, nil
		}
		if sortRule.IsPattern {
			re, err := regexp.Compile(sortRule.Text)
			if err != nil {
				return nil, err
			}
			for index := 0; index < len(sourceFields); index++ {
				field := sourceFields[index]
				if re.MatchString(field.name) {
					mapping[destinationIndex] = recordcopier.IndexMap{
						DestinationIndex: destinationIndex,
						SourceIndex:      field.index,
					}
					destinationIndex++
					sourceFields = append(sourceFields[:index], sourceFields[index+1:]...)
					index--
				}
			}
			continue
		}
		for index := 0; index < len(sourceFields); index++ {
			field := sourceFields[index]
			if field.name == sortRule.Text {
				mapping[destinationIndex] = recordcopier.IndexMap{
					DestinationIndex: destinationIndex,
					SourceIndex:      field.index,
				}
				destinationIndex++
				sourceFields = append(sourceFields[:index], sourceFields[index+1:]...)
				index--
			}
		}
	}

	for _, field := range sourceFields {
		mapping[destinationIndex] = recordcopier.IndexMap{
			DestinationIndex: destinationIndex,
			SourceIndex:      field.index,
		}
		destinationIndex++
	}

	return mapping, nil
}
