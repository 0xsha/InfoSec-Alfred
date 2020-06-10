package scrap

import (
	"testing"
)

func TestFetchNetSecReddit(t *testing.T) {

	got , err := FetchNetSecReddit()
	want := 10

	if err != nil{
		t.Errorf("Error Fetching NetSec %s",err)
	}

	if len(got) != want{

		t.Errorf("Expected 10 items got %d",len(got))

	}
}
