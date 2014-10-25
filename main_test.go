package wbem

import (
	"flag"
	_ "net/http/pprof"
	"net/url"
	"testing"
)

var (
	schema = flag.String("schema", "http", "")
	host   = flag.String("host", "192.168.1.14", "")
	port   = flag.String("port", "5988", "")
	path   = flag.String("path", "/cimom", "")

	username     = flag.String("username", "root", "")
	userpassword = flag.String("password", "root", "")
)

func getTestUri() url.URL {
	return url.URL{
		Scheme: *schema,
		User:   url.UserPassword(*username, *userpassword),
		Host:   *host + ":" + *port,
		Path:   *path,
	}
}

func TestEnumerateClassNames(t *testing.T) {
	c, e := NewClientCIMXML(getTestUri(), false)
	if nil != e {
		t.Error(e)
		return
	}
	names, e := c.EnumerateClassNames("root/cimv2", "", false)
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

	names2, e := c.EnumerateClassNames("root/cimv2", "", true)
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
	// go func() {
	// 	http.ListenAndServe(":", nil)
	// }()

	c, e := NewClientCIMXML(getTestUri(), false)
	if nil != e {
		t.Error(e)
		return
	}
	names, e := c.EnumerateInstanceNames("root/cimv2", "Linux_UnixProcess")
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
		//t.Log(name)
		instance, e := c.GetInstanceByInstanceName("root/cimv2", name, true, true, true, nil)
		if nil != e {
			t.Error(e)
			continue
		}
		// if 0 == len(instances) {
		// 	t.Error("instances list is emtpy")
		// 	continue
		// }

		t.Log(name)

		for i := 0; i < instance.GetPropertyCount(); i++ {
			pr := instance.GetPropertyByIndex(i)
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

		instance2, e := c.GetInstance(namespace, className, keyBindings, true, true, true, nil)
		if nil != e {
			t.Error(e)
			continue
		}
		if instance2.GetPropertyCount() != instance.GetPropertyCount() {
			t.Error("property count isn't equals.")
			continue
		}
	}
}

func TestEnumerateInstances(t *testing.T) {
	// go func() {
	// 	http.ListenAndServe(":", nil)
	// }()

	c, e := NewClientCIMXML(getTestUri(), false)
	if nil != e {
		t.Error(e)
		return
	}
	instances, e := c.EnumerateInstances("root/cimv2", "Linux_UnixProcess", false, false, false, false, nil)
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

func TestGetClass(t *testing.T) {
	// go func() {
	// 	http.ListenAndServe(":", nil)
	// }()

	c, e := NewClientCIMXML(getTestUri(), false)
	if nil != e {
		t.Error(e)
		return
	}
	class, e := c.GetClass("root/cimv2", "Linux_UnixProcess", true, true, true, nil)
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
	// go func() {
	//  http.ListenAndServe(":", nil)
	// }()

	c, e := NewClientCIMXML(getTestUri(), false)
	if nil != e {
		t.Error(e)
		return
	}
	classes, e := c.EnumerateClasses("root/cimv2", "CIM_UnixProcess", true, true, true, true)
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
