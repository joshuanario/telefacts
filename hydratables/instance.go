package hydratables

import (
	"fmt"
	"sort"

	"ecksbee.com/telefacts/attr"
	"ecksbee.com/telefacts/serializables"
)

type Instance struct {
	FileName      string
	Contexts      []Context
	Units         []Unit
	FootnoteLinks []FootnoteLink
	Facts         []Fact
}

type Fact struct {
	Href         string
	PrefixedName string
	ID           string
	ContextRef   string
	Context      *Context
	UnitRef      string
	Unit         *Unit
	Decimals     string
	Precision    string
	XMLInner     string
}

type ExplicitMember struct {
	Dimension struct {
		Href  string
		Value string
	}
	Member struct {
		Href     string
		CharData string
	}
}

type TypedMember struct {
	Dimension struct {
		Href  string
		Value string
	}
	XMLInner string
}

type Context struct {
	ID     string
	Entity Entity
	Period struct {
		Instant  *Instant
		Duration *Duration
	}
	Scenario *DimensionContext
}

type Entity struct {
	Identifier struct {
		Scheme   string
		CharData string
	}
	Segment *DimensionContext
}

type DimensionContext struct {
	ExplicitMembers []ExplicitMember
	TypedMembers    []TypedMember
}

type Instant struct {
	CharData string
}

type Duration struct {
	StartDate string
	EndDate   string
}

type Unit struct {
	ID      string
	Measure *UnitMeasure
	Divide  *UnitDivide
}

type UnitMeasure struct {
	Href     string
	CharData string
}

type UnitDivide struct {
	UnitNumerator struct {
		Measure *UnitMeasure
	}
	UnitDenominator struct {
		Measure *UnitMeasure
	}
}

type FootnoteLink struct {
	Title        string
	Footnotes    []Footnote
	Locs         []FootnoteLinkLoc
	FootnoteArcs []FootnoteArc
}

type Footnote struct {
	ID       string
	CharData string
	Lang     string
	Label    string
}

type FootnoteLinkLoc struct {
	Href  string
	Fact  *Fact
	Label string
}

type FootnoteArc struct {
	From     string
	Loc      *FootnoteLinkLoc
	To       string
	Footnote *Footnote
}

func HydrateInstance(file *serializables.InstanceFile, fileName string) (*Instance, error) {
	if len(fileName) <= 0 {
		return nil, fmt.Errorf("Empty file name")
	}
	if file == nil {
		return nil, fmt.Errorf("Empty file")
	}
	ret := Instance{}
	ret.FileName = fileName
	ret.Contexts = hydrateContexts(file)
	ret.Units = hydrateUnits(file)
	ret.Facts = hydrateFacts(file)
	ret.FootnoteLinks = hydrateFootnoteLinks(file)
	return &ret, nil
}

func hydrateContexts(instanceFile *serializables.InstanceFile) []Context {
	ret := make([]Context, len(instanceFile.Context))
	for _, context := range instanceFile.Context {
		nsAttr := attr.FindAttr(context.XMLAttrs, "xmlns")
		if nsAttr == nil || nsAttr.Value != attr.XBRLI {
			continue
		}
		item := Context{}
		idAttr := attr.FindAttr(context.XMLAttrs, "id")
		if idAttr == nil || idAttr.Value == "" {
			continue
		}
		item.ID = idAttr.Value
		newEntity := Entity{}
		entity := context.Entity[0]
		if len(entity.Identifier) <= 0 {
			continue
		}
		if entity.Identifier[0].CharData == "" {
			continue
		}
		schemeAttr := attr.FindAttr(entity.Identifier[0].XMLAttrs, "scheme")
		if schemeAttr == nil {
			continue
		}
		newEntity.Identifier.CharData = entity.Identifier[0].CharData
		newEntity.Identifier.Scheme = schemeAttr.Value
		if len(entity.Segment) > 0 {
			segment := DimensionContext{}
			if len(entity.Segment[0].ExplicitMember) > 0 {
				segment.ExplicitMembers = make([]ExplicitMember, 0, len(entity.Segment[0].ExplicitMember))
				for _, explicitMember := range entity.Segment[0].ExplicitMember {
					dimAttr := attr.FindAttr(explicitMember.XMLAttrs, "dimension")
					if dimAttr == nil {
						continue
					}
					newExplicitMember := ExplicitMember{
						Dimension: struct {
							Href  string
							Value string
						}{
							Href:  "", //todo
							Value: dimAttr.Value,
						},
						Member: struct {
							Href     string
							CharData string
						}{
							Href:     "", //todo
							CharData: explicitMember.CharData,
						},
					}
					segment.ExplicitMembers = append(segment.ExplicitMembers, newExplicitMember)
				}
			}
			if len(entity.Segment[0].TypedMember) > 0 {
				segment.TypedMembers = make([]TypedMember, 0, len(entity.Segment[0].TypedMember))
				for _, typedMember := range entity.Segment[0].TypedMember {
					dimAttr := attr.FindAttr(typedMember.XMLAttrs, "dimension")
					if dimAttr == nil {
						continue
					}
					newTypedMember := TypedMember{
						Dimension: struct {
							Href  string
							Value string
						}{
							Href:  "", //todo
							Value: dimAttr.Value,
						},
						XMLInner: typedMember.XMLInner,
					}
					segment.TypedMembers = append(segment.TypedMembers, newTypedMember)
				}
			}
			newEntity.Segment = &segment
		}
		item.Entity = newEntity
		if len(context.Period) > 0 {
			if len(context.Period[0].Instant) > 0 {
				instant := Instant{
					CharData: context.Period[0].Instant[0].CharData,
				}
				item.Period.Instant = &instant
			}
			if len(context.Period[0].StartDate) > 0 && len(context.Period[0].EndDate) > 0 {
				duration := Duration{
					StartDate: context.Period[0].StartDate[0].CharData,
					EndDate:   context.Period[0].EndDate[0].CharData,
				}
				item.Period.Duration = &duration
			}
		}
		if len(context.Scenario) > 0 {
			scenario := DimensionContext{}
			if len(context.Scenario[0].ExplicitMember) > 0 {
				scenario.ExplicitMembers = make([]ExplicitMember, 0, len(context.Scenario[0].ExplicitMember))
				for _, explicitMember := range context.Scenario[0].ExplicitMember {
					dimAttr := attr.FindAttr(explicitMember.XMLAttrs, "dimension")
					if dimAttr == nil || dimAttr.Value == "" {
						continue
					}
					newExplicitMember := ExplicitMember{
						Dimension: struct {
							Href  string
							Value string
						}{
							Href:  "", //todo
							Value: dimAttr.Value,
						},
						Member: struct {
							Href     string
							CharData string
						}{
							Href:     "", //todo
							CharData: explicitMember.CharData,
						},
					}
					scenario.ExplicitMembers = append(scenario.ExplicitMembers, newExplicitMember)
				}
			}
			if len(context.Scenario[0].TypedMember) > 0 {
				scenario.TypedMembers = make([]TypedMember, 0, len(context.Scenario[0].TypedMember))
				for _, typedMember := range context.Scenario[0].TypedMember {
					dimAttr := attr.FindAttr(typedMember.XMLAttrs, "dimension")
					if dimAttr == nil || dimAttr.Value == "" {
						continue
					}
					newTypedMember := TypedMember{
						Dimension: struct {
							Href  string
							Value string
						}{
							Href:  "", //todo
							Value: dimAttr.Value,
						},
						XMLInner: typedMember.XMLInner,
					}
					scenario.TypedMembers = append(scenario.TypedMembers, newTypedMember)
				}
			}
			item.Scenario = &scenario
		}
		ret = append(ret, item)
	}
	sort.SliceStable(ret, func(i int, j int) bool {
		return ret[i].ID < ret[j].ID
	})
	return ret
}

func hydrateUnits(instanceFile *serializables.InstanceFile) []Unit {
	ret := make([]Unit, 0, len(instanceFile.Unit))
	for _, unit := range instanceFile.Unit {
		nsAttr := attr.FindAttr(unit.XMLAttrs, "xmlns")
		if nsAttr == nil || nsAttr.Value != attr.XBRLI {
			continue
		}
		idAttr := attr.FindAttr(unit.XMLAttrs, "id")
		if idAttr == nil || idAttr.Value == "" {
			continue
		}
		item := Unit{}
		if len(unit.Measure) <= 0 {
			if len(unit.Divide) <= 0 {
				continue
			}
			if len(unit.Divide[0].UnitDenominator) <= 0 || len(unit.Divide[0].UnitNumerator) <= 0 {
				continue
			}
			if len(unit.Divide[0].UnitDenominator[0].Measure) <= 0 || len(unit.Divide[0].UnitNumerator[0].Measure) <= 0 {
				continue
			}
			numeratorMeasure := UnitMeasure{
				Href:     "", //todo
				CharData: unit.Divide[0].UnitNumerator[0].Measure[0].CharData,
			}
			denominatorMeasure := UnitMeasure{
				Href:     "", //todo
				CharData: unit.Divide[0].UnitDenominator[0].Measure[0].CharData,
			}
			divide := UnitDivide{
				UnitNumerator: struct{ Measure *UnitMeasure }{
					Measure: &numeratorMeasure,
				},
				UnitDenominator: struct{ Measure *UnitMeasure }{
					Measure: &denominatorMeasure,
				},
			}
			item.Divide = &divide
		} else {
			item.Measure = &UnitMeasure{
				Href:     "",
				CharData: unit.Measure[0].CharData,
			}
		}
		item.ID = idAttr.Value
		ret = append(ret, item)
	}

	sort.SliceStable(ret, func(i int, j int) bool {
		return ret[i].ID < ret[j].ID
	})
	return ret
}

func hydrateFacts(instanceFile *serializables.InstanceFile) []Fact {
	ret := make([]Fact, 0, len(instanceFile.Facts))
	for _, fact := range instanceFile.Facts {
		idAttr := attr.FindAttr(fact.XMLAttrs, "id")
		if idAttr == nil || idAttr.Value == "" {
			continue
		}
		if fact.XMLName.Local == "" || fact.XMLName.Space == "" {
			continue
		}
		contextRefAttr := attr.FindAttr(fact.XMLAttrs, "contextRef")
		if contextRefAttr == nil || contextRefAttr.Value == "" {
			continue
		}
		unitRefAttr := attr.FindAttr(fact.XMLAttrs, "unitRef")
		if unitRefAttr == nil || unitRefAttr.Value == "" {
			continue
		}
		decimalsAttr := attr.FindAttr(fact.XMLAttrs, "decimals")
		if decimalsAttr == nil || decimalsAttr.Value == "" {
			continue
		}
		precisionAttr := attr.FindAttr(fact.XMLAttrs, "precision")
		if precisionAttr == nil || precisionAttr.Value == "" {
			continue
		}
		newFact := Fact{
			ID:           idAttr.Value,
			Href:         "", //todo
			PrefixedName: "", //todo
			ContextRef:   contextRefAttr.Value,
			Context:      nil, //todo
			UnitRef:      unitRefAttr.Value,
			Unit:         nil, //todo
			Decimals:     decimalsAttr.Value,
			Precision:    precisionAttr.Value,
		}
		ret = append(ret, newFact)
	}
	sort.SliceStable(ret, func(i int, j int) bool {
		return ret[i].ID < ret[j].ID
	})
	return ret
}

func hydrateFootnoteLinks(instanceFile *serializables.InstanceFile) []FootnoteLink {
	ret := make([]FootnoteLink, 0, len(instanceFile.FootnoteLink))
	for _, footnoteLink := range instanceFile.FootnoteLink {
		item := FootnoteLink{}
		nsAttr := attr.FindAttr(footnoteLink.XMLAttrs, "xmlns")
		if nsAttr == nil || nsAttr.Value != attr.XLINK {
			continue
		}
		typeAttr := attr.FindAttr(footnoteLink.XMLAttrs, "type")
		if typeAttr == nil || typeAttr.Value != "extended" {
			continue
		}
		roleAttr := attr.FindAttr(footnoteLink.XMLAttrs, "xmlns")
		if roleAttr == nil || roleAttr.Value != attr.ROLELINK {
			continue
		}
		titleAttr := attr.FindAttr(footnoteLink.XMLAttrs, "title")
		if titleAttr != nil {
			item.Title = titleAttr.Value
		}
		item.Locs = make([]FootnoteLinkLoc, 0, len(footnoteLink.Loc))
		for _, loc := range footnoteLink.Loc {
			loctypeAttr := attr.FindAttr(loc.XMLAttrs, "type")
			if loctypeAttr == nil || loctypeAttr.Value != "locator" {
				continue
			}
			locnsAttr := attr.FindAttr(loc.XMLAttrs, "xmlns")
			if locnsAttr == nil || locnsAttr.Value != attr.XLINK {
				continue
			}
			hrefAttr := attr.FindAttr(loc.XMLAttrs, "href")
			if hrefAttr == nil || hrefAttr.Value == "" {
				continue
			}
			labelAttr := attr.FindAttr(loc.XMLAttrs, "label")
			if labelAttr == nil || labelAttr.Value == "" {
				continue
			}
			newLoc := FootnoteLinkLoc{
				Href:  hrefAttr.Value,
				Fact:  nil, //todo
				Label: labelAttr.Value,
			}
			item.Locs = append(item.Locs, newLoc)
		}
		item.Footnotes = make([]Footnote, 0, len(footnoteLink.Footnote))
		for _, footnote := range footnoteLink.Footnote {
			footnotetypeAttr := attr.FindAttr(footnote.XMLAttrs, "type")
			if footnotetypeAttr == nil || footnotetypeAttr.Value != "resource" {
				continue
			}
			footnoteidAttr := attr.FindAttr(footnote.XMLAttrs, "id")
			if footnoteidAttr == nil || footnoteidAttr.Value == "" {
				continue
			}
			footnotensAttr := attr.FindAttr(footnote.XMLAttrs, "xmlns")
			if footnotensAttr == nil || footnotensAttr.Value != attr.LINK {
				continue
			}
			footnotelabelAttr := attr.FindAttr(footnote.XMLAttrs, "label")
			if footnotelabelAttr == nil || footnotelabelAttr.Value == "" {
				continue
			}
			footnotelangAttr := attr.FindAttr(footnote.XMLAttrs, "lang")
			if footnotelangAttr == nil || footnotelangAttr.Value == "" {
				continue
			}
			newFootnote := Footnote{
				ID:       footnoteidAttr.Value,
				Label:    footnotelabelAttr.Value,
				Lang:     footnotelangAttr.Value,
				CharData: footnote.CharData,
			}
			item.Footnotes = append(item.Footnotes, newFootnote)
		}
		item.FootnoteArcs = make([]FootnoteArc, 0, len(footnoteLink.FootnoteArc))
		for _, footnoteArc := range footnoteLink.FootnoteArc {
			footnoteArcnsAttr := attr.FindAttr(footnoteArc.XMLAttrs, "xmlns")
			if footnoteArcnsAttr == nil || footnoteArcnsAttr.Value != attr.LINK {
				continue
			}
			footnoteArcarcroleAttr := attr.FindAttr(footnoteArc.XMLAttrs, "arcrole")
			if footnoteArcarcroleAttr == nil || footnoteArcarcroleAttr.Value != attr.FactFootnoteArcrole {
				continue
			}
			footnoteArctypeAttr := attr.FindAttr(footnoteArc.XMLAttrs, "type")
			if footnoteArctypeAttr == nil || footnoteArctypeAttr.Value != "arc" {
				continue
			}
			footnoteArcfromAttr := attr.FindAttr(footnoteArc.XMLAttrs, "from")
			if footnoteArcfromAttr == nil || footnoteArcfromAttr.Value == "" {
				continue
			}
			footnoteArctoAttr := attr.FindAttr(footnoteArc.XMLAttrs, "to")
			if footnoteArctoAttr == nil || footnoteArctoAttr.Value == "" {
				continue
			}
			var locPtr *FootnoteLinkLoc
			for _, lloc := range item.Locs {
				if lloc.Label == footnoteArcfromAttr.Value {
					locPtr = &lloc
					break
				}
			}
			var footnotePtr *Footnote
			for _, ffootnote := range item.Footnotes {
				if ffootnote.Label == footnoteArctoAttr.Value {
					footnotePtr = &ffootnote
					break
				}
			}
			if locPtr == nil || footnotePtr == nil {
				continue
			}
			newFootnoteArc := FootnoteArc{
				From:     footnoteArcfromAttr.Value,
				To:       footnoteArctoAttr.Value,
				Loc:      locPtr,
				Footnote: footnotePtr,
			}
			item.FootnoteArcs = append(item.FootnoteArcs, newFootnoteArc)
		}
		ret = append(ret, item)
	}
	sort.SliceStable(ret, func(i int, j int) bool {
		return ret[i].Title < ret[j].Title
	})
	return ret
}