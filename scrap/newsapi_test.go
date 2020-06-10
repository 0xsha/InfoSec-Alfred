package scrap

import (
	"testing"
)

func TestFetchNewsAPI(t *testing.T) {

	got  , err:= FetchNewsAPI()

	if err !=nil{

		t.Errorf("Error Fetching NewsAPI %s" ,err)
	}

	if len(got) < 10{
		t.Errorf("Expected at least 10 items got %d" , len(got))
	}
}