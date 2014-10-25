package wbem

import (
	"flag"
	"net/http"
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

	go func() {
		http.ListenAndServe(":", nil)
	}()

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
		instances, e := c.GetInstance("root/cimv2", name, true, true, true, nil)
		if nil != e {
			t.Error(e)
			continue
		}
		if 0 == len(instances) {
			t.Error("instances list is emtpy")
			continue
		}

		t.Log(name)
		for _, instance := range instances {
			for i := 0; i < instance.GetPropertyCount(); i++ {
				pr := instance.GetPropertyByIndex(i)
				t.Log(pr.GetName(), "=", pr.GetValue())
			}
		}
	}

}
