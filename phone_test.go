package phone

import (
	"fmt"
	"testing"
)

func Test_Find_1(t *testing.T) {
	phoneInfo, err := Find("1888888")
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("%+v\n", phoneInfo)
		t.Log("first test passed")
	}

}

func Test_Find_2(t *testing.T) {
	_, err := Find("1999999")
	if err == nil {
		t.Error("Find did not work as expected")
	} else {
		t.Log("second test passed ", err)
	}

}
