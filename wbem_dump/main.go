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
	"time"

	"github.com/runner-mei/gowbem"
)

var (
	schema    = flag.String("schema", "http", "")
	host      = flag.String("host", "192.168.1.157", "")
	port      = flag.String("port", "5988", "")
	path      = flag.String("path", "/cimom", "")
	namespace = flag.String("namespace", "root/cimv2", "")

	username     = flag.String("username", "root", "")
	userpassword = flag.String("password", "root", "")
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
	debug := flag.Bool("debug", false, "")
	output := flag.String("output", ".", "")
	flag.Parse()

	if *debug {
		gowbem.SetDebugProvider(&gowbem.FileDebugProvider{Path: *output})
	}

	c, e := gowbem.NewClientCIMXML(createURI(), true)
	if nil != e {
		log.Fatalln("连接失败，", e)
	}
	/*{
		timeCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		instances, err := c.EnumerateInstances(timeCtx, *namespace, "CIM_DiskDrive", true, false, false, false, nil)
		if err != nil {
			if !gowbem.IsErrNotSupported(err) && !gowbem.IsEmptyResults(err) {
				fmt.Println(fmt.Sprintf("%T %v", err, err))
			}
			return
		}

		fmt.Println()
		fmt.Println()
		for _, instance := range instances {
			for _, k := range instance.GetInstance().GetProperties() {
				fmt.Println(k.GetName(), k.GetValue())
			}
		}
		return
	}*/

	timeCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var namespaces, err = c.EnumerateNamespaces(timeCtx, nil)
	if nil != err {
		log.Fatalln("连接失败，", err)
	}
	for _, ns := range namespaces {
		fmt.Println("开始处理", ns)
		classNames, e := c.EnumerateClassNames(context.Background(), ns, "", true)
		if nil != e {
			if !gowbem.IsErrNotSupported(err) && !gowbem.IsEmptyResults(err) {
				fmt.Println("枚举类名失败，", e)
			}
			continue
		}
		if 0 == len(classNames) {
			fmt.Println("没有类定义？，")
			continue
		}
		for _, className := range classNames {
			timeCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
			class, err := c.GetClass(timeCtx, ns, className, true, true, true, nil)
			if err != nil {
				fmt.Println("取数名失败 - ", err)
			}

			/// @begin 将类定义写到文件
			filename := filepath.Join(*output, ns, className+".xml")
			if err := os.MkdirAll(filepath.Join(*output, ns), 666); err != nil && !os.IsExist(err) {
				log.Fatalln(err)
			}
			if err := ioutil.WriteFile(filename, []byte(class), 666); err != nil {
				log.Fatalln(err)
			}
			/// @end

			timeCtx, _ = context.WithTimeout(context.Background(), 30*time.Second)
			instanceNames, err := c.EnumerateInstanceNames(timeCtx, ns, className)
			if err != nil {
				fmt.Println(className, 0)

				if !gowbem.IsErrNotSupported(err) && !gowbem.IsEmptyResults(err) {
					fmt.Println(fmt.Sprintf("%T %v", err, err))
				}
				continue
			}
			fmt.Println(className, len(instanceNames))

			/// @begin 将类定义写到文件
			classPath := filepath.Join(*output, ns, className)
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
