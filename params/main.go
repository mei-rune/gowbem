package params

import (
	"fmt"

	"github.com/runner-mei/gowbem"
)

func Value(name, value string) gowbem.CIMParamValue {
	return &gowbem.CimParamValue{
		Name:  name,
		Value: &gowbem.CimValue{Value: value},
	}
}

func LocalClassPath(name, namespaceName, className string) gowbem.CIMParamValue {
	names := gowbem.SplitNamespaces(namespaceName)
	namespaces := make([]gowbem.CimNamespace, len(names))
	for idx, name := range names {
		namespaces[idx].Name = name
	}

	return &gowbem.CimParamValue{
		Name: name,
		ValueReference: &gowbem.CimValueReference{
			LocalClassPath: &gowbem.CimLocalClassPath{
				NamespacePath: gowbem.CimLocalNamespacePath{Namespaces: namespaces},
				ClassName:     gowbem.CimClassName{Name: className},
			},
		},
	}
}

func ClassName(name, className string) gowbem.CIMParamValue {
	return &gowbem.CimParamValue{
		Name: name,
		ValueReference: &gowbem.CimValueReference{
			ClassName: &gowbem.CimClassName{Name: className},
		},
	}
}

func LocalInstancePath(name, namespaceName string, insName interface{}) gowbem.CIMParamValue {
	names := gowbem.SplitNamespaces(namespaceName)
	namespaces := make([]gowbem.CimNamespace, len(names))
	for idx, name := range names {
		namespaces[idx].Name = name
	}

	switch v := insName.(type) {
	case string:
		instanceName, e := gowbem.ParseInstanceName(v)
		if e != nil {
			panic(e)
		}
		return &gowbem.CimParamValue{
			Name: name,
			ValueReference: &gowbem.CimValueReference{
				LocalInstancePath: &gowbem.CimLocalInstancePath{
					LocalNamespacePath: gowbem.CimLocalNamespacePath{Namespaces: namespaces},
					InstanceName:       *instanceName,
				},
			},
		}

	case *gowbem.CimInstanceName:
		return &gowbem.CimParamValue{
			Name: name,
			ValueReference: &gowbem.CimValueReference{
				LocalInstancePath: &gowbem.CimLocalInstancePath{
					LocalNamespacePath: gowbem.CimLocalNamespacePath{Namespaces: namespaces},
					InstanceName:       *v,
				},
			},
		}

	case gowbem.CimInstanceName:
		return &gowbem.CimParamValue{
			Name: name,
			ValueReference: &gowbem.CimValueReference{
				LocalInstancePath: &gowbem.CimLocalInstancePath{
					LocalNamespacePath: gowbem.CimLocalNamespacePath{Namespaces: namespaces},
					InstanceName:       v,
				},
			},
		}

	default:
		panic(fmt.Errorf("unsupport type '%T' - %#v", insName, insName))
	}
}

func InstanceName(name string, insName interface{}) gowbem.CIMParamValue {
	switch v := insName.(type) {
	case string:
		instanceName, e := gowbem.ParseInstanceName(v)
		if e != nil {
			panic(e)
		}
		return &gowbem.CimParamValue{
			Name: name,
			ValueReference: &gowbem.CimValueReference{
				InstanceName: instanceName,
			},
		}

	case *gowbem.CimInstanceName:
		return &gowbem.CimParamValue{
			Name: name,
			ValueReference: &gowbem.CimValueReference{
				InstanceName: v,
			},
		}

	case gowbem.CimInstanceName:
		return &gowbem.CimParamValue{
			Name: name,
			ValueReference: &gowbem.CimValueReference{
				InstanceName: &v,
			},
		}

	default:
		panic(fmt.Errorf("unsupport type '%T' - %#v", insName, insName))
	}
}

// type CimValueReference struct {
//   XMLName           xml.Name              `xml:"VALUE.REFERENCE"`
//   ClassPath         *CimClassPath         `xml:"CLASSPATH,omitempty"`
//   LocalClassPath    *CimLocalClassPath    `xml:"LOCALCLASSPATH,omitempty"`
//   ClassName         *CimClassName         `xml:"CLASSNAME,omitempty"`
//   InstancePath      *CimInstancePath      `xml:"INSTANCEPATH,omitempty"`
//   LocalInstancePath *CimLocalInstancePath `xml:"LOCALINSTANCEPATH,omitempty"`
//   InstanceName      *CimInstanceName      `xml:"INSTANCENAME,omitempty"`
// }
