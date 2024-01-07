package response

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestErrorExtra(t *testing.T) {
	SkipConvey("Test error extra", t, func() {
		// err := apperrors.NewCustomError(apperrors.RequestParamErr, "request param error")
		// So(err, ShouldNotBeNil)
		// // fmt.Printf("err: %+v\n", err.Error())
		// data := ErrCodeResp(err)
		// // fmt.Printf("data: %+v\n", data)
		// So(data, ShouldNotBeNil)
	})
}
