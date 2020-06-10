package scrap

import (
	"testing"
)

func TestFetchGithubAdvisory(t *testing.T) {

	got , err := FetchGithubAdvisory()
	want := 10

	if err != nil{
		t.Errorf("Error Fetching Github %s" , err)
	}

	if len(got) != want{

		t.Errorf("Excepted 10 items got %d",len(got))
	}

}
