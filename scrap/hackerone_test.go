package scrap

import (
	"testing"
)

func TestFetchHackerOne(t *testing.T) {

	got , err := FetchHackerOne()
	want := 25
	if err != nil{

		t.Errorf("Error Fetching HackerOne %s",err)
	}

	if len(got.HacktivityItems.Edges) !=want{

		t.Errorf("Expected 25 Items got %d", len(got.HacktivityItems.Edges) )

	}


}