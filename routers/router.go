// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"dataAnalysis/controllers"
	_ "dataAnalysis/utils"

	"github.com/astaxie/beego"
)

func init() {
	beego.ErrorController(&controllers.ErrorController{})
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/object", beego.NSInclude(&controllers.ObjectController{})),
		beego.NSNamespace("/user", beego.NSInclude(&controllers.UserController{})),
		beego.NSNamespace("/excel", beego.NSInclude(&controllers.ExcelController{})),
	)
	beego.AddNamespace(ns)
}
