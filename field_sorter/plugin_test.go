package field_sorter

import (
	"encoding/xml"
	"github.com/tlarsen7572/goalteryx/recordcopier"
	"testing"
)

func TestXmlParsing(t *testing.T) {
	config := &config{}
	err := xml.Unmarshal([]byte(`<Configuration>
  <alphabetical>false</alphabetical>
  <field0>
    <text>Symbol</text>
    <isPattern>false</isPattern>
  </field0>
  <field1>
    <text>Buy.*</text>
    <isPattern>true</isPattern>
  </field1>
  <field2>
    <text>Sell.*</text>
    <isPattern>true</isPattern>
  </field2>
</Configuration>`), config)

	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if config.Alphabetical {
		t.Fatalf(`expected not alphabetical but got True`)
	}
	if count := len(config.Fields); count != 3 {
		t.Fatalf(`expected 3 fields but got %v`, count)
	}
	if text := config.Fields[0].Text; text != `Symbol` {
		t.Fatalf(`expected 'Symbol' but got '%v'`, text)
	}
	if text := config.Fields[1].Text; text != `Buy.*` {
		t.Fatalf(`expected 'Buy.*' but got '%v'`, text)
	}
	if text := config.Fields[2].Text; text != `Sell.*` {
		t.Fatalf(`expected 'Sell.*' but got '%v'`, text)
	}
	if isPattern := config.Fields[0].IsPattern; isPattern {
		t.Fatalf(`expected IsPattern of false but got true`)
	}
	if isPattern := config.Fields[1].IsPattern; !isPattern {
		t.Fatalf(`expected IsPattern of true but got false`)
	}
	if isPattern := config.Fields[2].IsPattern; !isPattern {
		t.Fatalf(`expected IsPattern of true but got false`)
	}
}

func TestXmlParsingAlphabetical(t *testing.T) {
	config := &config{}
	err := xml.Unmarshal([]byte(`<Configuration>
  <alphabetical>true</alphabetical>
  <field0>
    <text>Symbol</text>
    <isPattern>false</isPattern>
  </field0>
  <field1>
    <text>Buy.*</text>
    <isPattern>true</isPattern>
  </field1>
  <field2>
    <text>Sell.*</text>
    <isPattern>true</isPattern>
  </field2>
</Configuration>`), config)

	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if !config.Alphabetical {
		t.Fatalf(`expected alphabetical but got false`)
	}
}

func TestSortFields(t *testing.T) {
	sortInfo := []field{
		{Text: `Symbol`, IsPattern: false},
		{Text: `Buy.*`, IsPattern: true},
		{Text: `Sell.*`, IsPattern: true},
	}

	incomingFields := []string{
		`BuyMin`,
		`Symbol`,
		`SellMax`,
		`SellMin`,
		`BuyMax`,
	}

	mapping, err := sortFields(incomingFields, sortInfo, false)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(mapping); count != 5 {
		t.Fatalf(`expected 5 mapping entries but got %v`, count)
	}

	expected := []recordcopier.IndexMap{
		{
			DestinationIndex: 0,
			SourceIndex:      1,
		},
		{
			DestinationIndex: 1,
			SourceIndex:      0,
		},
		{
			DestinationIndex: 2,
			SourceIndex:      4,
		},
		{
			DestinationIndex: 3,
			SourceIndex:      2,
		},
		{
			DestinationIndex: 4,
			SourceIndex:      3,
		},
	}
	for index, entry := range mapping {
		if entry.SourceIndex != expected[index].SourceIndex {
			t.Fatalf(`entry %v expected SourceIndex %v but got %v`, index, expected[index].SourceIndex, entry.SourceIndex)
		}
	}
}

func TestSortFieldsAlphabetical(t *testing.T) {
	sortInfo := []field{
		{Text: `Symbol`, IsPattern: false},
		{Text: `Buy.*`, IsPattern: true},
		{Text: `Sell.*`, IsPattern: true},
	}

	incomingFields := []string{
		`BuyMin`,
		`Symbol`,
		`SellMax`,
		`SellMin`,
		`BuyMax`,
	}

	mapping, err := sortFields(incomingFields, sortInfo, true)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(mapping); count != 5 {
		t.Fatalf(`expected 5 mapping entries but got %v`, count)
	}

	expected := []recordcopier.IndexMap{
		{
			DestinationIndex: 0,
			SourceIndex:      1,
		},
		{
			DestinationIndex: 1,
			SourceIndex:      4,
		},
		{
			DestinationIndex: 2,
			SourceIndex:      0,
		},
		{
			DestinationIndex: 3,
			SourceIndex:      2,
		},
		{
			DestinationIndex: 4,
			SourceIndex:      3,
		},
	}
	for index, entry := range mapping {
		if entry.SourceIndex != expected[index].SourceIndex {
			t.Fatalf(`entry %v expected SourceIndex %v but got %v`, index, expected[index].SourceIndex, entry.SourceIndex)
		}
	}
}

func TestSortFieldsNoDuplicates(t *testing.T) {
	sortInfo := []field{
		{Text: `SellMin`, IsPattern: false},
		{Text: `Symbol`, IsPattern: false},
		{Text: `Buy.*`, IsPattern: true},
		{Text: `Sell.*`, IsPattern: true},
	}

	incomingFields := []string{
		`BuyMin`,
		`Symbol`,
		`SellMax`,
		`SellMin`,
		`BuyMax`,
	}

	mapping, err := sortFields(incomingFields, sortInfo, true)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(mapping); count != 5 {
		t.Fatalf(`expected 5 mapping entries but got %v`, count)
	}

	expected := []recordcopier.IndexMap{
		{
			DestinationIndex: 0,
			SourceIndex:      3,
		},
		{
			DestinationIndex: 1,
			SourceIndex:      1,
		},
		{
			DestinationIndex: 2,
			SourceIndex:      4,
		},
		{
			DestinationIndex: 3,
			SourceIndex:      0,
		},
		{
			DestinationIndex: 4,
			SourceIndex:      2,
		},
	}
	for index, entry := range mapping {
		if entry.SourceIndex != expected[index].SourceIndex {
			t.Fatalf(`entry %v expected SourceIndex %v but got %v`, index, expected[index].SourceIndex, entry.SourceIndex)
		}
	}
}

func TestSortUnmappedFields(t *testing.T) {
	sortInfo := []field{
		{Text: `Symbol`, IsPattern: false},
		{Text: `Buy.*`, IsPattern: true},
		{Text: `Sell.*`, IsPattern: true},
	}

	incomingFields := []string{
		`BuyMin`,
		`zzz`,
		`Blah blah`,
		`Symbol`,
		`SellMax`,
		`SellMin`,
		`BuyMax`,
	}

	mapping, err := sortFields(incomingFields, sortInfo, true)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(mapping); count != 7 {
		t.Fatalf(`expected 7 mapping entries but got %v`, count)
	}

	expected := []recordcopier.IndexMap{
		{
			DestinationIndex: 0,
			SourceIndex:      3,
		},
		{
			DestinationIndex: 1,
			SourceIndex:      6,
		},
		{
			DestinationIndex: 2,
			SourceIndex:      0,
		},
		{
			DestinationIndex: 3,
			SourceIndex:      4,
		},
		{
			DestinationIndex: 4,
			SourceIndex:      5,
		},
		{
			DestinationIndex: 5,
			SourceIndex:      2,
		},
		{
			DestinationIndex: 6,
			SourceIndex:      1,
		},
	}
	for index, entry := range mapping {
		if entry.SourceIndex != expected[index].SourceIndex {
			t.Fatalf(`entry %v expected SourceIndex %v but got %v`, index, expected[index].SourceIndex, entry.SourceIndex)
		}
	}
}
