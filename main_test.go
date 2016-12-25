package main


import (
	"testing"
	"io/ioutil"
	. "github.com/smartystreets/goconvey/convey"
	"fmt"
	"os"
)

func Test_clone(t *testing.T) {

	Convey("We can test cloning thisisfine's repos", t, func(){
		var err error
		path, err = ioutil.TempDir("", "fetch_test_")
		So(err, ShouldBeNil)

		err = clone()
		So(err, ShouldBeNil)

		dirs, err := ioutil.ReadDir(path)
		So(err, ShouldBeNil)
		for _, d := range dirs {
			fmt.Println(d)
		}

		err = os.RemoveAll(path)
		So(err, ShouldBeNil)
	})

}