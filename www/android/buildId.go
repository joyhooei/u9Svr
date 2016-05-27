package androidPackage

import (
	"github.com/astaxie/beego"
	"os"
	"u9/models"
	"u9/tool"
	"u9/tool/android"
	"strings"
	"io/ioutil"
)

const (
	projectPath = "package/project"
)

type BuildId struct {
	product        *models.Product
	productVersion *models.ProductVersion
	channel        *models.Channel
	packageParam   *models.PackageParam
	copyToPath     string
	buildIdPath    string
	packagePath    string
	channelPath    string
}

func NewBuildId(packageTaskId int,channel *models.Channel, product *models.Product,
	productVersion *models.ProductVersion , packageParam *models.PackageParam) *BuildId {
	ret := new(BuildId)
	ret.product, ret.productVersion = product, productVersion
	ret.channel = channel
	ret.packageParam = packageParam
	ret.channelPath = GetChannelPath(channel)
	apkName := GetApkName(product, productVersion)
	ret.packagePath = GetPackagePath(packageTaskId, apkName)
	ret.buildIdPath = GetBuildIdPath(packageTaskId, "/project")
	ret.copyToPath = GetBuildIdPath(packageTaskId, "")
	return ret
}

//流程：
// 初始化：
// 1、拷贝project 模板目录
// 2、拷贝打包目录 res、AndroidManifest.xml
// 反编译：

// 1、ant打包 "ant release"
// 2、反编译 classes.dex 输出至打包目录下(smali文件夹下)

func (this *BuildId) Handle() {
	switch this.channel.Id{
		case 106:
			fallthrough
		case 107:
			fallthrough
		case 122:
 			fallthrough
 		case 130:
 			fallthrough
 		// case 126:
 		// 	fallthrough
 		case 136:
			this.init()
			this.ant()
			this.dex()
		case 139:
			this.init()
			this.tencent()
		default :
			
	}
}

func (this *BuildId) init() {
	os.RemoveAll(this.copyToPath)
	tool.CreateDir(this.copyToPath)
	tool.CopyDir(projectPath, this.copyToPath)
	tool.CopyDir(this.packagePath+"/res", this.buildIdPath)
	tool.CopyFile(this.packagePath+"/AndroidManifest.xml", this.buildIdPath+"/AndroidManifest.xml")
}
func (this *BuildId) ant() {
	beego.Trace("buildIdPath:", this.buildIdPath)
	beego.Trace("packagePath:", this.packagePath)
	if err := android.Ant(this.buildIdPath, "release"); err != nil {
		beego.Trace("ant release err:", err)
		beego.Trace("ant release err:", err)
		panic(err)
	}
}

func(this *BuildId) dex(){
		if err := android.UnCompileSmallDex(this.buildIdPath+"/bin/classes.dex", this.packagePath+"/smali"); err != nil {
		// "package/build_id/out"
		beego.Trace("UnCompileSmallDex erro:", err)
		panic(err)
	}
}

func (this *BuildId) tencent() {
	//1、准备环境 可以直接使用BuildId init
	//2、拷贝YSDK jar包  
	//3、创建目录 
	//4、生成java 文件 编译
	//5、拷贝文件 项目\bin\classes\目录
	var err error
	content := "package " + this.packageParam.PackageName + ".wxapi"
	filePath := strings.Replace(this.packageParam.PackageName + ".wxapi", ".", "/", -1)
	smaliPath := this.packagePath +"/smali/"+filePath
	javaPath := this.buildIdPath + "/src/"+filePath
	classesFile := this.buildIdPath + "/bin/classes/"+filePath + "/WXEntryActivity.class"
	smaliFile := this.packagePath +"/smali/"+filePath + "/WXEntryActivity.class"
	tencentJar := this.channelPath+"/YSDK_Android_1.1.1_235.jar"
	cpTencetnJar := this.buildIdPath + "/libs/YSDK_Android_1.1.1_235.jar"
	content = content + ";\r\npublic class WXEntryActivity extends com.tencent.ysdk.module.user.impl.wx.YSDKWXEntryActivity{ }"
	d1 := []byte(content)
	if err = os.MkdirAll(javaPath, 0777);err != nil{
		beego.Trace(err)
		panic(err)
	}
	if err := ioutil.WriteFile(javaPath+"/WXEntryActivity.java", d1, 0644);err !=nil{
		beego.Trace(err)
		panic(err)
	}
	
	if _,err = tool.CopyFile(tencentJar,cpTencetnJar);err != nil {
		beego.Trace(err)
		panic(err)
	}
	this.ant()
	if err = os.MkdirAll(smaliPath, 0777);err != nil{
		beego.Trace(err)
		panic(err)
	}
	if _,err = tool.CopyFile(classesFile,smaliFile);err != nil {
		beego.Trace(err)
		panic(err)
	}
	
}
