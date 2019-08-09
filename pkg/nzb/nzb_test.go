package nzb

import (
	"reflect"
	"testing"
)

func TestSubject(t *testing.T) {
	f, err := FromFile("./testdata/test.nzb")
	if err != nil {
		t.Fatalf("failed to create NZB from file: %v", err)
	}
	if f.Subject != "ezNZB-01-09-2013 Test.mp3 - \"test.mp3\" yEnc (1/10)" {
		t.Errorf("expected subject not equal to actual subject '%s'", f.Subject)
	}
}

func TestSegments(t *testing.T) {
	f, err := FromFile("./testdata/test.nzb")
	if err != nil {
		t.Fatalf("failed to create NZB from file: %v", err)
	}
	ids := make([]string, len(f.Segments))
	for i, s := range f.Segments {
		ids[i] = s.ID
	}
	if !reflect.DeepEqual(ids, []string {
		"NewzToolz_Rulz!_www.techsono.com_3443298495_2277527@giganews.com",
		"NewzToolz_Rulz!_www.techsono.com_3443298506_2278166@giganews.com",
		"NewzToolz_Rulz!_www.techsono.com_3443298516_2278767@giganews.com",
		"NewzToolz_Rulz!_www.techsono.com_3443298528_2279477@giganews.com",
		"NewzToolz_Rulz!_www.techsono.com_3443298538_2280086@giganews.com",
		"NewzToolz_Rulz!_www.techsono.com_3443298548_2280689@giganews.com",
		"NewzToolz_Rulz!_www.techsono.com_3443298558_2281304@giganews.com",
		"NewzToolz_Rulz!_www.techsono.com_3443298568_2281932@giganews.com",
		"NewzToolz_Rulz!_www.techsono.com_3443298581_2282669@giganews.com",
		"NewzToolz_Rulz!_www.techsono.com_3443298591_2283280@giganews.com",
	}) {
		t.Errorf("ids not as expected: %v", ids)
	}
}