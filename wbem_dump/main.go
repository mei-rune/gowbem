package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/runner-mei/gowbem"
)

var (
	schema    = flag.String("schema", "http", "")
	host      = flag.String("host", "192.168.1.157", "主机的 IP 地址")
	port      = flag.String("port", "5988", "主机上 CIM 服务的端口号")
	path      = flag.String("path", "/cimom", "CIM 服务访问路径")
	namespace = flag.String("namespace", "root/cimv2", "CIM 的命名空间")
	classname = flag.String("class", "", "CIM 的的类名")

	username     = flag.String("username", "root", "用户名")
	userpassword = flag.String("password", "root", "用户密码")
	output       = flag.String("output", ".", "结果的输出目录")
	debug        = flag.Bool("debug", false, "是不是在调试")
)

func createURI() *url.URL {
	return &url.URL{
		Scheme: *schema,
		User:   url.UserPassword(*username, *userpassword),
		Host:   *host + ":" + *port,
		Path:   *path,
	}
}

func main() {
	flag.Usage = func() {
		fmt.Println("使用方法： wbem_dump -host=192.168.1.157 -port=5988 -username=root -password=rootpwd\r\n" +
			"可用选项")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *debug {
		gowbem.SetDebugProvider(&gowbem.FileDebugProvider{Path: *output})
	}

	c, e := gowbem.NewClientCIMXML(createURI(), true)
	if nil != e {
		log.Fatalln("连接失败，", e)
	}

	if *classname != "" && *namespace != "" {
		dumpClass(c, *namespace, *classname)
		return
	}

	timeCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var defaultList []string
	if "" != *namespace {
		defaultList = []string{*namespace}
	}
	var namespaces, err = c.EnumerateNamespaces(timeCtx, defaultList, 10*time.Second, nil)
	if nil != err {
		log.Fatalln("连接失败，", err)
	}
	for _, ns := range namespaces {
		fmt.Println("开始处理", ns)
		dumpNS(c, ns)
	}
	fmt.Println("导出成功！")

	// timeCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// computerSystems, err := c.EnumerateInstances(timeCtx, *namespace,
	// 	"CIM_ComputerSystem ", true, false, true, true, nil)
	// if err != nil {
	// 	if !gowbem.IsErrNotSupported(err) {
	// 		fmt.Println(fmt.Sprintf("%T %v", err, err))
	// 	}
	// 	return
	// }

	// for _, computer := range computerSystems {
	// 	timeCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// 	instances, err := c.AssociatorInstances(timeCtx, *namespace, computer.GetName(), "CIM_InstalledSoftwareIdentity",
	// 		"CIM_SoftwareIdentity",
	// 		"System", "InstalledSoftware", true, nil)
	// 	if err != nil {
	// 		if !gowbem.IsErrNotSupported(err) {
	// 			fmt.Println(fmt.Sprintf("%T %v", err, err))
	// 		}
	// 		continue
	// 	}

	// 	for _, instance := range instances {
	// 		fmt.Println("-----------------")
	// 		for _, k := range instance.GetProperties() {
	// 			fmt.Println(k.GetName(), k.GetValue())
	// 		}
	// 	}
	// }

	// fmt.Println("测试成功！")
}

func dumpNS(c *gowbem.ClientCIMXML, ns string) {
	classNames, e := c.EnumerateClassNames(context.Background(), ns, "", true)
	if nil != e {
		if !gowbem.IsErrNotSupported(e) && !gowbem.IsEmptyResults(e) {
			fmt.Println("枚举类名失败，", e)
		}
		return
	}
	if 0 == len(classNames) {
		fmt.Println("没有类定义？，")
		return
	}
	for _, className := range classNames {
		timeCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		class, err := c.GetClass(timeCtx, ns, className, true, true, true, nil)
		if err != nil {
			fmt.Println("取数名失败 - ", err)
		}

		nsPath := strings.Replace(ns, "/", "#", -1)
		nsPath = strings.Replace(nsPath, "\\", "@", -1)

		/// @begin 将类定义写到文件
		filename := filepath.Join(*output, nsPath, className+".xml")
		if err := os.MkdirAll(filepath.Join(*output, nsPath), 666); err != nil && !os.IsExist(err) {
			log.Fatalln(err)
		}
		if err := ioutil.WriteFile(filename, []byte(class), 666); err != nil {
			log.Fatalln(err)
		}
		/// @end

		dumpClass(c, ns, class)
	}
}

func dumpClass(c *gowbem.ClientCIMXML, ns, className string) {
	nsPath := strings.Replace(ns, "/", "#", -1)
	nsPath = strings.Replace(nsPath, "\\", "@", -1)

	timeCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	instanceNames, err := c.EnumerateInstanceNames(timeCtx, ns, className)
	if err != nil {
		fmt.Println(className, 0)

		if !gowbem.IsErrNotSupported(err) && !gowbem.IsEmptyResults(err) {
			fmt.Println(fmt.Sprintf("%T %v", err, err))
		}
		return
	}
	fmt.Println(className, len(instanceNames))

	/// @begin 将类定义写到文件
	classPath := filepath.Join(*output, nsPath, className)
	if err := os.MkdirAll(classPath, 666); err != nil && !os.IsExist(err) {
		log.Fatalln(err)
	}
	var buf bytes.Buffer
	for _, instanceName := range instanceNames {
		buf.WriteString(instanceName.String())
		buf.WriteString("\r\n")
	}
	if err := ioutil.WriteFile(filepath.Join(classPath, "instances.txt"), buf.Bytes(), 666); err != nil {
		log.Fatalln(err)
	}
	/// @end

	for idx, instanceName := range instanceNames {
		timeCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		instance, err := c.GetInstanceByInstanceName(timeCtx, ns, instanceName, false, true, true, nil)
		if err != nil {
			if !gowbem.IsErrNotSupported(err) && !gowbem.IsEmptyResults(err) {
				fmt.Println(fmt.Sprintf("%T %v", err, err))
			}
			continue
		}

		/// @begin 将类定义写到文件
		bs, err := xml.MarshalIndent(instance, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}
		if err := ioutil.WriteFile(filepath.Join(classPath, "instance_"+strconv.Itoa(idx)+".xml"), bs, 666); err != nil {
			log.Fatalln(err)
		}
		/// @end

		// fmt.Println()
		// fmt.Println()
		// fmt.Println(instanceName.String())
		//for _, k := range instance.GetProperties() {
		//	fmt.Println(k.GetName(), k.GetValue())
		//}
	}
}
