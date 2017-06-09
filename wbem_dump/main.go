package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
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
	flag.Parse()

	if *debug {
		gowbem.SetDebugProvider(&gowbem.FileDebugProvider{Path: "."})
	}

	c, e := gowbem.NewClientCIMXML(createURI(), true)
	if nil != e {
		log.Fatalln("连接失败，", e)
	}
	{
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
	}

	timeCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var namespaces, err = c.EnumerateNamespaces(timeCtx, nil)
	if nil != err {
		log.Fatalln("连接失败，", err)
	}
	for _, ns := range namespaces {
		fmt.Println("开始处理", ns)
		names, e := c.EnumerateClassNames(context.Background(), ns, "", true)
		if nil != e {
			if !gowbem.IsErrNotSupported(err) && !gowbem.IsEmptyResults(err) {
				fmt.Println("枚举类名失败，", e)
			}
			continue
		}
		if 0 == len(names) {
			fmt.Println("没有类定义？，")
			continue
		}
		for _, name := range names {
			timeCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
			instanceNames, err := c.EnumerateInstanceNames(timeCtx, ns, name)
			if err != nil {
				fmt.Println(name, 0)

				if !gowbem.IsErrNotSupported(err) && !gowbem.IsEmptyResults(err) {
					fmt.Println(fmt.Sprintf("%T %v", err, err))
				}
				continue
			}

			fmt.Println(name, len(instanceNames))

			for _, instanceName := range instanceNames {
				timeCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
				_, err := c.GetInstanceByInstanceName(timeCtx, ns, instanceName, false, true, true, nil)
				if err != nil {
					if !gowbem.IsErrNotSupported(err) && !gowbem.IsEmptyResults(err) {
						fmt.Println(fmt.Sprintf("%T %v", err, err))
					}
					continue
				}

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
