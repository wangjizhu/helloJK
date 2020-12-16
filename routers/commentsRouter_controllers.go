package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["helloprecision/controllers:GPController"] = append(beego.GlobalControllerRouter["helloprecision/controllers:GPController"],
        beego.ControllerComments{
            Method: "BorrowSampleShelf",
            Router: "/BorrowSampleShelf/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["helloprecision/controllers:GPController"] = append(beego.GlobalControllerRouter["helloprecision/controllers:GPController"],
        beego.ControllerComments{
            Method: "GetLengthOfThread",
            Router: "/GetLengthOfThread/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["helloprecision/controllers:GPController"] = append(beego.GlobalControllerRouter["helloprecision/controllers:GPController"],
        beego.ControllerComments{
            Method: "ReturnSampleShelf",
            Router: "/ReturnSampleShelf/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["helloprecision/controllers:GPController"] = append(beego.GlobalControllerRouter["helloprecision/controllers:GPController"],
        beego.ControllerComments{
            Method: "SetResourceSample",
            Router: "/SetResourceSample/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["helloprecision/controllers:GPController"] = append(beego.GlobalControllerRouter["helloprecision/controllers:GPController"],
        beego.ControllerComments{
            Method: "StartMultipleThreads",
            Router: "/StartMultipleThreads/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["helloprecision/controllers:GPController"] = append(beego.GlobalControllerRouter["helloprecision/controllers:GPController"],
        beego.ControllerComments{
            Method: "StartSingleThread",
            Router: "/StartSingleThread/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
