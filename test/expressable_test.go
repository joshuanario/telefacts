package telefacts_test

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"ecksbee.com/telefacts/pkg/hydratables"
	"ecksbee.com/telefacts/pkg/renderables"
	"ecksbee.com/telefacts/pkg/serializables"
	gocache "github.com/patrickmn/go-cache"
)

func TestCatalog_Expressables(t *testing.T) {
	hcache := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	serializables.WorkingDirectoryPath = filepath.Join(".", "wd")
	serializables.GlobalTaxonomySetPath = filepath.Join(".", "gts")
	hydratables.InjectCache(hcache)
	workingDir := filepath.Join(".", "wd", "folders", "test_ix")
	_, err := os.Stat(workingDir)
	if os.IsNotExist(err) {
		os.MkdirAll(workingDir, fs.FileMode(0700))
	}
	defer func() {
		os.RemoveAll(workingDir)
	}()
	zipFile := filepath.Join(".", "wd", "test_ix.zip")
	err = unZipTestData(workingDir, zipFile)
	if err != nil {
		t.Fatalf("Error: " + err.Error())
		return
	}
	f, err := serializables.Discover("test_ix")
	if err != nil {
		t.Fatalf("Error: " + err.Error())
	}
	h, err := hydratables.Hydrate(f)
	if err != nil {
		t.Fatalf("Error: " + err.Error())
	}
	data, err := renderables.MarshalCatalog(h)
	if err != nil {
		t.Fatalf("Error: " + err.Error())
	}
	c := renderables.Catalog{}
	err = json.Unmarshal(data, &c)
	if err != nil {
		t.Fatalf("Error: " + err.Error())
	}
	if c.DocumentName != "cmg-20200331x10q.htm" {
		t.Fatalf("expected cmg-20200331x10q.htm; outcome %s;\n", c.DocumentName)
	}
	data, err = renderables.MarshalExpressable("us-gaap:EffectOfExchangeRateOnCashCashEquivalentsRestrictedCashAndRestrictedCashEquivalents", "Duration_1_1_2020_To_3_31_2020", h)
	if err != nil {
		t.Fatalf("Error: " + err.Error())
	}
	e := renderables.Expressable{}
	err = json.Unmarshal(data, &e)
	if err != nil {
		t.Fatalf("Error: " + err.Error())
	}
	if e.Href != "http://xbrl.fasb.org/us-gaap/2019/elts/us-gaap-2019-01-31.xsd#us-gaap_EffectOfExchangeRateOnCashCashEquivalentsRestrictedCashAndRestrictedCashEquivalents" {
		t.Fatalf("expected http://xbrl.fasb.org/us-gaap/2019/elts/us-gaap-2019-01-31.xsd#us-gaap_EffectOfExchangeRateOnCashCashEquivalentsRestrictedCashAndRestrictedCashEquivalents; outcome %s;\n", e.Href)
	}
	if e.Context.Period[renderables.PureLabel] != "2020-01-01/2020-03-31" {
		t.Fatalf("expected 2020-01-01/2020-03-31; outcome %s;\n", e.Context.Period[renderables.PureLabel])
	}
	if e.Context.Period[renderables.BriefLabel] != "2020-01-01/2020-03-31" {
		t.Fatalf("expected 2020-01-01/2020-03-31; outcome %s;\n", e.Context.Period[renderables.BriefLabel])
	}
	if e.Context.Period[renderables.English] != "3 months ended March 31, 2020" {
		t.Fatalf("expected 3 months ended March 31, 2020; outcome %s;\n", e.Context.Period[renderables.English])
	}
	data, err = renderables.MarshalExpressable("us-gaap:StockholdersEquity", "As_Of_12_31_2018_us-gaap_StatementEquityComponentsAxis_us-gaap_AccumulatedNetUnrealizedInvestmentGainLossMember", h)
	if err != nil {
		t.Fatalf("Error: " + err.Error())
	}
	e2 := renderables.Expressable{}
	err = json.Unmarshal(data, &e2)
	if err != nil {
		t.Fatalf("Error: " + err.Error())
	}
	if e2.Href != "http://xbrl.fasb.org/us-gaap/2019/elts/us-gaap-2019-01-31.xsd#us-gaap_StockholdersEquity" {
		t.Fatalf("expected http://xbrl.fasb.org/us-gaap/2019/elts/us-gaap-2019-01-31.xsd#us-gaap_StockholdersEquity; outcome %s;\n", e2.Href)
	}
	if e2.Context.Period[renderables.PureLabel] != "2018-12-31" {
		t.Fatalf("expected 2018-12-31; outcome %s;\n", e2.Context.Period[renderables.PureLabel])
	}
	if e2.Context.Period[renderables.BriefLabel] != "2018-12-31" {
		t.Fatalf("expected 2018-12-31; outcome %s;\n", e2.Context.Period[renderables.BriefLabel])
	}
	if e2.Context.Period[renderables.English] != "as of December 31, 2018" {
		t.Fatalf("expected as of December 31, 2018; outcome %s;\n", e2.Context.Period[renderables.English])
	}
}

func Test485BPOS_Expressables(t *testing.T) {
	hcache := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	serializables.WorkingDirectoryPath = filepath.Join(".", "wd")
	serializables.GlobalTaxonomySetPath = filepath.Join(".", "gts")
	hydratables.InjectCache(hcache)
	workingDir := filepath.Join(".", "wd", "folders", "test_485_ix")
	_, err := os.Stat(workingDir)
	if os.IsNotExist(err) {
		os.MkdirAll(workingDir, fs.FileMode(0700))
	}
	defer func() {
		os.RemoveAll(workingDir)
	}()
	zipFile := filepath.Join(".", "wd", "test_485_ix.zip")
	err = unZipTestData(workingDir, zipFile)
	if err != nil {
		t.Fatalf("Error: " + err.Error())
		return
	}
	f, err := serializables.Discover("test_485_ix")
	if err != nil {
		t.Fatalf("Error: " + err.Error())
	}
	h, err := hydratables.Hydrate(f)
	if err != nil {
		t.Fatalf("Error: " + err.Error())
	}
	data, err := renderables.MarshalCatalog(h)
	if err != nil {
		t.Fatalf("Error: " + err.Error())
	}
	c := renderables.Catalog{}
	err = json.Unmarshal(data, &c)
	if err != nil {
		t.Fatalf("Error: " + err.Error())
	}
	if c.DocumentName != "d394191d485bpos.htm" {
		t.Fatalf("expected d394191d485bpos.htm; outcome %s;\n", c.DocumentName)
	}
	data, err = renderables.MarshalExpressable("rr:PortfolioTurnoverRate", "S000002724Member_InvestorACInstitutionalAndClassRMember", h)
	if err != nil {
		t.Fatalf("Error: " + err.Error())
	}
	e := renderables.Expressable{}
	err = json.Unmarshal(data, &e)
	if err != nil {
		t.Fatalf("Error: " + err.Error())
	}
	unlabelledFact := (*e.Expression)[renderables.PureLabel]
	expression := unlabelledFact.Head + unlabelledFact.Core + unlabelledFact.Tail
	if expression != "1.0200 pure" {
		t.Fatalf("expected 1.0200; outcome %s;\n", expression)
	}
}
