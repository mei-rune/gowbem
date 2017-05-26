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
	host      = flag.String("host", "192.168.1.14", "")
	port      = flag.String("port", "5988", "")
	path      = flag.String("path", "/cimom", "")
	namespace = flag.String("namespace", "root/cimv2", "")

	username     = flag.String("username", "root", "")
	userpassword = flag.String("password", "", "")
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
	names, e := c.EnumerateClassNames(context.Background(), *namespace, "", false)
	if nil != e {
		log.Fatalln("枚举类名失败，", e)
	}
	if 0 == len(names) {
		log.Fatalln("没有类定义？，")
	}
	for _, name := range names {
		fmt.Println(name)

		// timeCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		// instances, err := c.EnumerateInstances(timeCtx,
		// 	*namespace, name, true, false, true, true, nil)
		// if err != nil {
		// 	if !gowbem.IsErrNotSupported(err) {
		// 		fmt.Println(fmt.Sprintf("%T %v", err, err))
		// 	}
		// 	continue
		// }
		// for _, instance := range instances {
		// 	fmt.Println(spew.Sprint(instance))
		// }
	}

	fmt.Println("测试成功！")

	timeCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	computerSystems, err := c.EnumerateInstances(timeCtx, *namespace,
		"CIM_ComputerSystem ", true, false, true, true, nil)
	if err != nil {
		if !gowbem.IsErrNotSupported(err) {
			fmt.Println(fmt.Sprintf("%T %v", err, err))
		}
		return
	}

	for _, computer := range computerSystems {
		timeCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		instances, err := c.AssociatorInstances(timeCtx, *namespace, computer.GetName(), "CIM_InstalledSoftwareIdentity",
			"CIM_SoftwareIdentity",
			"System", "InstalledSoftware", true, nil)
		if err != nil {
			if !gowbem.IsErrNotSupported(err) {
				fmt.Println(fmt.Sprintf("%T %v", err, err))
			}
			continue
		}

		for _, instance := range instances {
			fmt.Println("-----------------")
			for _, k := range instance.GetProperties() {
				fmt.Println(k.GetName(), k.GetValue())
			}
		}
	}

	fmt.Println("测试成功！")
}
