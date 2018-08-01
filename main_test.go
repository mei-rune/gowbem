package gowbem_test

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"strings"
	"testing"
	"time"

	. "github.com/runner-mei/gowbem"
	"github.com/runner-mei/gowbem/params"
)

var (
	schema = flag.String("scheme", "http", "")
	host   = flag.String("host", "192.168.1.14", "")
	port   = flag.String("port", "5988", "")
	path   = flag.String("path", "/cimom", "")

	username     = flag.String("username", "root", "")
	userpassword = flag.String("password", "", "")
)

func getTestUri() *url.URL {
	return &url.URL{
		Scheme: *schema,
		User:   url.UserPassword(*username, *userpassword),
		Host:   *host + ":" + *port,
		Path:   *path,
	}
}

func TestEnumerateClassNames(t *testing.T) {
	if "" == *userpassword {
		t.Skip("please input password.")
	}
	c, e := NewClientCIMXML(getTestUri(), false)
	if nil != e {
		t.Error(e)
		return
	}

	timerCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	names, e := c.EnumerateClassNames(timerCtx, "root/cimv2", "", false)
	if nil != e {
		t.Error(e)
		return
	}
	if 0 == len(names) {
		t.Error("class list is emtpy")
		return
	}
	for _, name := range names {
		t.Log(name)
	}

	timerCtx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	names2, e := c.EnumerateClassNames(timerCtx, "root/cimv2", "", true)
	if nil != e {
		t.Error(e)
		return
	}
	if 0 == len(names2) {
		t.Error("class list is emtpy")
		return
	}
	for _, name := range names2 {
		t.Log(name)
	}

	if len(names) >= len(names2) {
		t.Error("len(names) >= len(names2)")
	}
}

func TestEnumerateInstanceNames(t *testing.T) {
	if "" == *userpassword {
		t.Skip("please input password.")
	}
	// go func() {
	// 	http.ListenAndServe(":", nil)
	// }()

	c, e := NewClientCIMXML(getTestUri(), false)
	if nil != e {
		t.Error(e)
		return
	}
	timerCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	names, e := c.EnumerateInstanceNames(timerCtx, "root/cimv2", "Linux_UnixProcess")
	if nil != e {
		t.Error(e)
		return
	}
	if 0 == len(names) {
		t.Error("class list is emtpy")
		return
	}

	for idx, name := range names {
		if idx > 10 {
			break
		}

		timerCtx, _ = context.WithTimeout(context.Background(), 10*time.Second)
		instance1, e := c.GetInstanceByInstanceName(timerCtx, "root/cimv2", name, true, true, true, nil)
		if nil != e {
			t.Error(e)
			continue
		}
		t.Log(instance1)

		for i := 0; i < instance1.GetPropertyCount(); i++ {
			pr := instance1.GetPropertyByIndex(i)
			t.Log(pr.GetName(), "=", pr.GetValue())
		}

		namespace, className, keyBindings, e := Parse(name.String())
		if nil != e {
			t.Error(e)
			continue
		}
		if "" == namespace {
			namespace = "root/cimv2"
		}

		timerCtx, _ = context.WithTimeout(context.Background(), 10*time.Second)
		instance2, e := c.GetInstance(timerCtx, namespace, className, keyBindings, true, true, true, nil)
		if nil != e {
			t.Error(e)
			continue
		}

		if instance2.GetPropertyCount() != instance1.GetPropertyCount() {
			t.Error("property count isn't equals.")
			continue
		}
	}
}

func TestEnumerateInstances(t *testing.T) {
	if "" == *userpassword {
		t.Skip("please input password.")
	}
	// go func() {
	// 	http.ListenAndServe(":", nil)
	// }()

	c, e := NewClientCIMXML(getTestUri(), false)
	if nil != e {
		t.Error(e)
		return
	}
	timerCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	instances, e := c.EnumerateInstances(timerCtx, "root/cimv2", "Linux_UnixProcess", false, false, false, false, nil)
	if nil != e {
		t.Error(e)
		return
	}
	if 0 == len(instances) {
		t.Error("instance list is emtpy")
		return
	}

	for _, instance := range instances {
		t.Log(instance.GetName())
		for i := 0; i < instance.GetInstance().GetPropertyCount(); i++ {
			pr := instance.GetInstance().GetPropertyByIndex(i)
			t.Log(pr.GetName(), "=", pr.GetValue())
		}
	}
}

func TestEnumerateInstances2(t *testing.T) {
	if "" == *userpassword {
		t.Skip("please input password.")
	}
	// go func() {
	// 	http.ListenAndServe(":", nil)
	// }()

	c, e := NewClientCIMXML(getTestUri(), false)
	if nil != e {
		t.Error(e)
		return
	}
	timerCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	instances, e := c.EnumerateInstances(timerCtx, "root/cimv2", "CIM_Processor", false, false, false, false, nil)
	if nil != e {
		t.Error(e)
		return
	}
	if 0 == len(instances) {
		t.Error("instance list is emtpy")
		return
	}

	for _, instance := range instances {
		t.Log(instance.GetName())
		for i := 0; i < instance.GetInstance().GetPropertyCount(); i++ {
			pr := instance.GetInstance().GetPropertyByIndex(i)
			t.Log(pr.GetName(), "=", pr.GetValue())
		}

		//instanceName := instance
		//instanceName.GetKeyBindings().Get(0).GetValue()

		timerCtx, _ = context.WithTimeout(context.Background(), 10*time.Second)
		instanceWithNames, e := c.AssociatorInstances(timerCtx, "root/cimv2", instance.GetName(), "", "", "", "", false, nil)
		if nil != e {
			//t.Error(e)
			fmt.Println(e)
			continue
		}
		if 0 == len(instanceWithNames) {
			continue
		}

		for _, instance := range instanceWithNames {
			//t.Log(instance.GetName())
			for i := 0; i < instance.GetInstance().GetPropertyCount(); i++ {
				pr := instance.GetInstance().GetPropertyByIndex(i)
				t.Log("\t", pr.GetName(), "=", pr.GetValue())
			}
		}
	}
}

func TestGetClass(t *testing.T) {
	if "" == *userpassword {
		t.Skip("please input password.")
	}
	// go func() {
	// 	http.ListenAndServe(":", nil)
	// }()

	c, e := NewClientCIMXML(getTestUri(), false)
	if nil != e {
		t.Error(e)
		return
	}
	timerCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	class, e := c.GetClass(timerCtx, "root/cimv2", "Linux_UnixProcess", true, true, true, nil)
	if nil != e {
		t.Error(e)
		return
	}
	if 0 == len(class) {
		t.Error("class is emtpy")
		return
	}

	t.Log(class)
}

func TestEnumerateClasses(t *testing.T) {
	if "" == *userpassword {
		t.Skip("please input password.")
	}
	// go func() {
	//  http.ListenAndServe(":", nil)
	// }()

	c, e := NewClientCIMXML(getTestUri(), false)
	if nil != e {
		t.Error(e)
		return
	}

	timerCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	classes, e := c.EnumerateClasses(timerCtx, "root/cimv2", "CIM_UnixProcess", true, true, true, true)
	if nil != e {
		t.Error(e)
		return
	}
	if 0 == len(classes) {
		t.Error("classes is emtpy")
		return
	}

	t.Log(classes)
}

func TestAssociatorInstances(t *testing.T) {
	t.Skip("timeout.")

	if "" == *userpassword {
		t.Skip("please input password.")
	}
	go func() {
		http.ListenAndServe(":", nil)
	}()

	c, e := NewClientCIMXML(getTestUri(), false)
	if nil != e {
		t.Error(e)
		return
	}

	timerCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	names, e := c.EnumerateClassNames(timerCtx, "root/cimv2", "", true)
	if nil != e {
		t.Error(e)
		return
	}
	if 0 == len(names) {
		t.Error("class list is emtpy")
		return
	}
	fmt.Println(len(names))
	for idx, name := range names {
		if "Linux_SysfsAttribute" == name {
			continue
		}
		fmt.Println(idx, "========", name)

		timerCtx, _ = context.WithTimeout(context.Background(), 10*time.Second)
		instances, e := c.EnumerateInstanceNames(timerCtx, "root/cimv2", name) //, false, false, false, false, nil)
		if nil != e {
			//t.Error(e)
			//return
			continue
		}
		if 0 == len(instances) {
			//t.Error("instance list is emtpy")
			continue
		}

		for _, instance := range instances {

			if strings.Contains(instance.String(), "Linux_SysfsAttribute") {
				continue
			}

			if strings.Contains(instance.String(), "Linux_SysfsBusDevice") {
				continue
			}

			if strings.Contains(instance.String(), "Linux_UnixProcess") {
				continue
			}

			if strings.Contains(instance.String(), "PG_UnixProcess") {
				continue
			}

			fmt.Println(instance)

			if 0 == len(instances) {
				//t.Error("instance list is emtpy")
				continue
			}

			for _, instance := range instances {
				fmt.Println("==========", instance)
				//instanceName := instance
				//instanceName.GetKeyBindings().Get(0).GetValue()

				timerCtx, _ = context.WithTimeout(context.Background(), 10*time.Second)
				instanceWithNames, e := c.AssociatorInstances(timerCtx, "root/cimv2", instance, "", "", "", "", false, nil)
				if nil != e {
					//t.Error(e)
					fmt.Println(e)
					continue
				}
				if 0 == len(instanceWithNames) {
					//t.Error("instance list is emtpy")
					continue
				}

				for _, instance := range instanceWithNames {
					//t.Log(instance.GetName())
					for i := 0; i < instance.GetInstance().GetPropertyCount(); i++ {
						pr := instance.GetInstance().GetPropertyByIndex(i)
						t.Log(pr.GetName(), "=", pr.GetValue())
					}
				}

			}
		}
	}

}

// func TestAssociatorClasses(t *testing.T) {
// 	c, e := NewClientCIMXML(getTestUri(), false)
// 	if nil != e {
// 		t.Error(e)
// 		return
// 	}

// 	// classes, e := c.EnumerateClasses("root/cimv2", "Linux_UnixProcess", true, true, true, true)
// 	// if nil != e {
// 	// 	t.Error(e)
// 	// 	return
// 	// }
// 	// has_ok := false
// 	// for _, cls := range classes {
// 	// var clsInstance CimClassInnerXml
// 	// if e := xml.Unmarshal([]byte(cls), &clsInstance); nil != e {
// 	// 	t.Log(e)
// 	// 	continue
// 	// }

// 	fmt.Println(clsInstance.Name)
// 	assoc_classes, e := c.AssociatorClasses("root/cimv2", "Linux_UnixProcess", "", "", "", "", true, true, nil)
// 	if nil != e {
// 		t.Log(e)
// 		return
// 	}
// 	if 0 == len(assoc_classes) {
// 		t.Log("classes is emtpy")
// 		return
// 	}
//
// 	//has_ok = true
// 	t.Log(assoc_classes)
// 	fmt.Println(assoc_classes)
// 	//}
// 	// if !has_ok {
// 	// 	t.Error("failed")
// 	// }
// }

// func TestReferenceClasses(t *testing.T) {
// 	c, e := NewClientCIMXML(getTestUri(), false)
// 	if nil != e {
// 		t.Error(e)
// 		return
// 	}

// 	classes, e := c.EnumerateClasses("root/cimv2", "", true, true, true, true)
// 	if nil != e {
// 		t.Error(e)
// 		return
// 	}
// 	has_ok := false
// 	for _, cls := range classes {
// 		var clsInstance CimClassInnerXml
// 		if e := xml.Unmarshal([]byte(cls), &clsInstance); nil != e {
// 			t.Log(e)
// 			continue
// 		}

// 		fmt.Println(clsInstance.Name)
// 		assoc_classes, e := c.ReferenceClasses("root/cimv2", clsInstance.Name, "", "", true, true, nil)
// 		if nil != e {
// 			t.Log(e)
// 			continue
// 		}
// 		if 0 == len(assoc_classes) {
// 			t.Log("classes is emtpy")
// 			continue
// 		}

// 		has_ok = true
// 		t.Log(assoc_classes)
// 		fmt.Println(assoc_classes)
// 	}
// 	if !has_ok {
// 		t.Error("failed")
// 	}
// }

// func TestReferenceNames(t *testing.T) {
// 	c, e := NewClientCIMXML(getTestUri(), false)
// 	if nil != e {
// 		t.Error(e)
// 		return
// 	}
// 	classes, e := c.EnumerateClasses("root/cimv2", "", true, true, true, true)
// 	if nil != e {
// 		t.Error(e)
// 		return
// 	}

// 	for _, cls := range classes {
// 		var clsInstance CimClassInnerXml
// 		if e := xml.Unmarshal([]byte(cls), &clsInstance); nil != e {
// 			t.Log(e)
// 			continue
// 		}

// 		fmt.Println(clsInstance.Name)

// 		instances, e := c.EnumerateInstances("root/cimv2", clsInstance.Name, true, true, true, true, nil)
// 		if nil != e {
// 			t.Error(e)
// 			continue
// 		}
// 		has_ok := false
// 		for _, instance := range instances {
// 			fmt.Println(instance.GetName().String())
// 			names, e := c.ReferenceNames("root/cimv2", instance.GetName(), "", "")
// 			if nil != e {
// 				t.Log(e)
// 				continue
// 			}
// 			if 0 == len(names) {
// 				t.Log("classes is emtpy")
// 				continue
// 			}

// 			has_ok = true
// 			t.Log(names)
// 			fmt.Println(names)
// 		}
// 		if !has_ok {
// 			t.Error("failed")
// 		}
// 	}
// }

func TestInvokeMethod(t *testing.T) {
	if "" == *userpassword {
		t.Skip("please input password.")
	}
	c, e := NewClientCIMXML(getTestUri(), false)
	if nil != e {
		t.Error(e)
		return
	}

	instanceName, e := ParseInstanceName(`PG_ConfigSetting.PropertyName="hostname"`)
	if nil != e {
		t.Error(e)
		return
	}

	timerCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	value, outParams, e := c.InvokeMethod(timerCtx, "root/PG_Internal", instanceName, "UpdatePropertyValue",
		[]CIMParamValue{params.Value("PropertyValue", "test"), params.Value("SetPlannedValue", "true")})
	if nil != e {
		t.Error(e)
		return
	}
	if 0 != len(outParams) {
		for _, param := range outParams {
			t.Log(param)
		}
	}

	if strings.ToLower(fmt.Sprint(value)) != "true" {
		t.Error("except is true got ", value)
	}

	timerCtx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	instanceValue, e := c.GetInstanceByInstanceName(timerCtx, "root/PG_Internal", instanceName, false, false, false, nil)
	if nil != e {
		t.Error(e)
		return
	}
	if instanceValue == nil {
		t.Error("instanceValue is nil")
		return
	}

	value = instanceValue.GetPropertyByName("PlannedValue").GetValue()
	if fmt.Sprint(value) != "test" {
		t.Error("except is test got ", value)
	}
}
