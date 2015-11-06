package gowbem

import (
	"errors"
	"net/url"
	"strings"
)

var (
	messageNotExists             = errors.New("CIM.MESSAGE isn't exists.")
	simpleReqNotExists           = errors.New("CIM.MESSAGE.SIMPLERSP isn't exists.")
	imethodResponseNotExists     = errors.New("CIM.MESSAGE.SIMPLERSP.IMETHODRESPONSE isn't exists.")
	ireturnValueNotExists        = errors.New("CIM.MESSAGE.SIMPLERSP.IMETHODRESPONSE.RETURNVALUE isn't exists.")
	instancePathNotExists        = errors.New("CIM.MESSAGE.SIMPLERSP.IMETHODRESPONSE.RETURNVALUE.INSTANCEPATH isn't exists.")
	instanceNamesNotExists       = errors.New("CIM.MESSAGE.SIMPLERSP.IMETHODRESPONSE.RETURNVALUE.INSTANCENAME isn't exists.")
	classNamesNotExists          = errors.New("CIM.MESSAGE.SIMPLERSP.IMETHODRESPONSE.RETURNVALUE.CLASSNAME isn't exists.")
	instancesNotExists           = errors.New("CIM.MESSAGE.SIMPLERSP.IMETHODRESPONSE.RETURNVALUE.INSTANCE isn't exists.")
	instancesMutiChioce          = errors.New("CIM.MESSAGE.SIMPLERSP.IMETHODRESPONSE.RETURNVALUE.INSTANCE is greate one.")
	valueNamedInstancesNotExists = errors.New("CIM.MESSAGE.SIMPLERSP.IMETHODRESPONSE.RETURNVALUE.VALUE.NAMEDINSTANCE isn't exists.")
	classesNotExists             = errors.New("CIM.MESSAGE.SIMPLERSP.IMETHODRESPONSE.RETURNVALUE.CLASS isn't exists.")
	classesMutiChioce            = errors.New("CIM.MESSAGE.SIMPLERSP.IMETHODRESPONSE.RETURNVALUE.CLASS is greate one.")
)

func booleanString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

type ClientCIMXML struct {
	Client

	CimVersion      string
	DtdVersion      string
	ProtocolVersion string
}

func (c *ClientCIMXML) init(u *url.URL, insecure bool) {
	c.Client.init(u, insecure)
	c.CimVersion = "2.0"
	c.DtdVersion = "2.0"
	c.ProtocolVersion = "1.0"

	//fmt.Println(c.Client.u.User)
}

func (c *ClientCIMXML) EnumerateClassNames(namespaceName, className string, deep bool) ([]string, error) {
	// obtain data
	if "" == namespaceName {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"namespace name is empty.")
	}

	names := strings.Split(namespaceName, "/")
	namespaces := make([]CimNamespace, len(names))
	for idx, name := range names {
		namespaces[idx].Name = name
	}

	var paramValues []CimIParamValue
	if deep {
		paramValues = []CimIParamValue{
			CimIParamValue{Name: "DeepInheritance", Value: &CimValue{Value: "true"}},
		}
	}

	if "" != className {
		paramValues = append(paramValues, CimIParamValue{
			Name:      "ClassName",
			ClassName: &CimClassName{Name: className},
		})
	}

	simpleReq := &CimSimpleReq{IMethodCall: &CimIMethodCall{
		Name:               "EnumerateClassNames",
		LocalNamespacePath: CimLocalNamespacePath{Namespaces: namespaces},
		ParamValues:        paramValues,
	}}

	req := &CIM{
		CimVersion: c.CimVersion,
		DtdVersion: c.DtdVersion,
		Message: &CimMessage{
			Id:              c.generateId(),
			ProtocolVersion: c.ProtocolVersion,
			SimpleReq:       simpleReq,
		},
		//Declaration: &CimDeclaration,
	}

	resp := &CIM{hasFault: func(cim *CIM) error {
		if nil == cim.Message {
			return messageNotExists
		}
		if nil == cim.Message.SimpleRsp {
			return simpleReqNotExists
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse {
			return imethodResponseNotExists
		}
		if nil != cim.Message.SimpleRsp.IMethodResponse.Error {
			e := cim.Message.SimpleRsp.IMethodResponse.Error
			return WBEMException(CIMStatusCode(e.Code), e.Description)
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse.ReturnValue {
			return ireturnValueNotExists
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse.ReturnValue.ClassNames {
			return classNamesNotExists
		}
		return nil
	}}

	if err := c.RoundTrip("POST", map[string]string{"CIMProtocolVersion": c.ProtocolVersion,
		"CIMOperation": "MethodCall",
		"CIMMethod":    "EnumerateClassNames",
		"CIMObject":    url.QueryEscape(namespaceName)}, req, resp); nil != err {
		return nil, err
	}

	results := make([]string, len(resp.Message.SimpleRsp.IMethodResponse.ReturnValue.ClassNames))
	for idx, name := range resp.Message.SimpleRsp.IMethodResponse.ReturnValue.ClassNames {
		results[idx] = name.Name
	}
	return results, nil
}

func (c *ClientCIMXML) EnumerateInstanceNames(namespaceName, className string) ([]CIMInstanceName, error) {
	if "" == namespaceName {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"namespace name is empty.")
	}

	if "" == className {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"class name is empty.")
	}

	names := strings.Split(namespaceName, "/")
	namespaces := make([]CimNamespace, len(names))
	for idx, name := range names {
		namespaces[idx].Name = name
	}

	paramValues := []CimIParamValue{
		CimIParamValue{
			Name:      "ClassName",
			ClassName: &CimClassName{Name: className},
		},
	}

	simpleReq := &CimSimpleReq{IMethodCall: &CimIMethodCall{
		Name:               "EnumerateInstanceNames",
		LocalNamespacePath: CimLocalNamespacePath{Namespaces: namespaces},
		ParamValues:        paramValues,
	}}

	req := &CIM{
		CimVersion: c.CimVersion,
		DtdVersion: c.DtdVersion,
		Message: &CimMessage{
			Id:              c.generateId(),
			ProtocolVersion: c.ProtocolVersion,
			SimpleReq:       simpleReq,
		},
		//Declaration: &CimDeclaration,
	}

	resp := &CIM{hasFault: func(cim *CIM) error {
		if nil == cim.Message {
			return messageNotExists
		}
		if nil == cim.Message.SimpleRsp {
			return simpleReqNotExists
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse {
			return imethodResponseNotExists
		}
		if nil != cim.Message.SimpleRsp.IMethodResponse.Error {
			e := cim.Message.SimpleRsp.IMethodResponse.Error
			return WBEMException(CIMStatusCode(e.Code), e.Description)
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse.ReturnValue {
			return ireturnValueNotExists
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse.ReturnValue.InstanceNames {
			return instanceNamesNotExists
		}
		return nil
	}}

	// CIMProtocolVersion: 1.0
	// CIMOperation: MethodCall
	// CIMMethod: EnumerateClassNames
	// CIMObject: root%2Fcimv2

	if err := c.RoundTrip("POST", map[string]string{"CIMProtocolVersion": c.ProtocolVersion,
		"CIMOperation": "MethodCall",
		"CIMMethod":    "EnumerateInstanceNames",
		"CIMObject":    url.QueryEscape(namespaceName)}, req, resp); nil != err {
		return nil, err
	}

	results := make([]CIMInstanceName, len(resp.Message.SimpleRsp.IMethodResponse.ReturnValue.InstanceNames))
	for idx, name := range resp.Message.SimpleRsp.IMethodResponse.ReturnValue.InstanceNames {
		results[idx] = name
	}
	return results, nil
}

func (c *ClientCIMXML) GetInstance(namespaceName, className string, keyBindings CIMKeyBindings, localOnly bool,
	includeQualifiers bool, includeClassOrigin bool, propertyList []string) (CIMInstance, error) {
	instanceName := &CimInstanceName{
		ClassName: className,
	}

	switch keyBindings.Len() {
	case 0:
		return nil, errors.New("keyBindings is empty.")
	case 1:
		kb := keyBindings.Get(0)
		if "_" == kb.GetName() {
			instanceName.KeyValue = kb.(*CimKeyBinding).KeyValue
			instanceName.ValueReference = kb.(*CimKeyBinding).ValueReference
			break
		}
		fallthrough
	default:
		instanceName.KeyBindings = keyBindings.(CimKeyBindings)
	}
	return c.GetInstanceByInstanceName(namespaceName, instanceName, localOnly, includeQualifiers, includeClassOrigin, propertyList)
}

func (c *ClientCIMXML) GetInstanceByInstanceName(namespaceName string, instanceName CIMInstanceName, localOnly bool,
	includeQualifiers bool, includeClassOrigin bool, propertyList []string) (CIMInstance, error) {
	if "" == namespaceName {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"namespace name is empty.")
	}

	if "" == instanceName.GetClassName() {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"class name is empty.")
	}

	names := strings.Split(namespaceName, "/")
	namespaces := make([]CimNamespace, len(names))
	for idx, name := range names {
		namespaces[idx].Name = name
	}

	paramValues := []CimIParamValue{
		CimIParamValue{
			Name:         "InstanceName",
			InstanceName: instanceName.(*CimInstanceName),
		},

		CimIParamValue{
			Name:  "LocalOnly",
			Value: &CimValue{Value: booleanString(localOnly)},
		},

		CimIParamValue{
			Name:  "IncludeQualifiers",
			Value: &CimValue{Value: booleanString(includeQualifiers)},
		},
		CimIParamValue{
			Name:  "IncludeClassOrigin",
			Value: &CimValue{Value: booleanString(includeClassOrigin)},
		},
	}
	if 0 != len(propertyList) {
		properties := CimValueArray(make([]CimValueOrNull, len(propertyList)))
		for idx, s := range propertyList {
			properties[idx] = CimValueOrNull{Value: &CimValue{Value: s}}
		}
		paramValues = append(paramValues,
			CimIParamValue{
				Name:       "PropertyList",
				ValueArray: properties,
			})
	}

	simpleReq := &CimSimpleReq{IMethodCall: &CimIMethodCall{
		Name:               "GetInstance",
		LocalNamespacePath: CimLocalNamespacePath{Namespaces: namespaces},
		ParamValues:        paramValues,
	}}

	req := &CIM{
		CimVersion: c.CimVersion,
		DtdVersion: c.DtdVersion,
		Message: &CimMessage{
			Id:              c.generateId(),
			ProtocolVersion: c.ProtocolVersion,
			SimpleReq:       simpleReq,
		},
		//Declaration: &CimDeclaration,
	}

	resp := &CIM{hasFault: func(cim *CIM) error {
		if nil == cim.Message {
			return messageNotExists
		}
		if nil == cim.Message.SimpleRsp {
			return simpleReqNotExists
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse {
			return imethodResponseNotExists
		}
		if nil != cim.Message.SimpleRsp.IMethodResponse.Error {
			e := cim.Message.SimpleRsp.IMethodResponse.Error
			return WBEMException(CIMStatusCode(e.Code), e.Description)
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse.ReturnValue {
			return ireturnValueNotExists
		}
		if 0 == len(cim.Message.SimpleRsp.IMethodResponse.ReturnValue.Instances) {
			return instancesNotExists
		}

		if 1 < len(cim.Message.SimpleRsp.IMethodResponse.ReturnValue.Instances) {
			return instancesMutiChioce
		}
		return nil
	}}

	// CIMProtocolVersion: 1.0
	// CIMOperation: MethodCall
	// CIMMethod: EnumerateClassNames
	// CIMObject: root%2Fcimv2

	if err := c.RoundTrip("POST", map[string]string{"CIMProtocolVersion": c.ProtocolVersion,
		"CIMOperation": "MethodCall",
		"CIMMethod":    "GetInstance",
		"CIMObject":    url.QueryEscape(namespaceName)}, req, resp); nil != err {
		return nil, err
	}

	if 0 == len(resp.Message.SimpleRsp.IMethodResponse.ReturnValue.Instances) {
		return nil, nil
	}
	return &resp.Message.SimpleRsp.IMethodResponse.ReturnValue.Instances[0], nil
}

func (c *ClientCIMXML) EnumerateInstances(namespaceName, className string, deepInheritance bool,
	localOnly bool, includeQualifiers bool, includeClassOrigin bool, propertyList []string) ([]CIMInstanceWithName, error) {

	if "" == namespaceName {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"namespace name is empty.")
	}

	if "" == className {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"class name is empty.")
	}

	names := strings.Split(namespaceName, "/")
	namespaces := make([]CimNamespace, len(names))
	for idx, name := range names {
		namespaces[idx].Name = name
	}

	paramValues := []CimIParamValue{
		CimIParamValue{
			Name:      "ClassName",
			ClassName: &CimClassName{Name: className},
		},

		CimIParamValue{
			Name:  "LocalOnly",
			Value: &CimValue{Value: booleanString(localOnly)},
		},

		CimIParamValue{
			Name:  "DeepInheritance",
			Value: &CimValue{Value: booleanString(deepInheritance)},
		},

		CimIParamValue{
			Name:  "IncludeQualifiers",
			Value: &CimValue{Value: booleanString(includeQualifiers)},
		},
		CimIParamValue{
			Name:  "IncludeClassOrigin",
			Value: &CimValue{Value: booleanString(includeClassOrigin)},
		},
	}
	if 0 != len(propertyList) {
		properties := CimValueArray(make([]CimValueOrNull, len(propertyList)))
		for idx, s := range propertyList {
			properties[idx] = CimValueOrNull{Value: &CimValue{Value: s}}
		}
		paramValues = append(paramValues,
			CimIParamValue{
				Name:       "PropertyList",
				ValueArray: properties,
			})
	}

	simpleReq := &CimSimpleReq{IMethodCall: &CimIMethodCall{
		Name:               "EnumerateInstances",
		LocalNamespacePath: CimLocalNamespacePath{Namespaces: namespaces},
		ParamValues:        paramValues,
	}}

	req := &CIM{
		CimVersion: c.CimVersion,
		DtdVersion: c.DtdVersion,
		Message: &CimMessage{
			Id:              c.generateId(),
			ProtocolVersion: c.ProtocolVersion,
			SimpleReq:       simpleReq,
		},
		//Declaration: &CimDeclaration,
	}

	resp := &CIM{hasFault: func(cim *CIM) error {
		if nil == cim.Message {
			return messageNotExists
		}
		if nil == cim.Message.SimpleRsp {
			return simpleReqNotExists
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse {
			return imethodResponseNotExists
		}
		if nil != cim.Message.SimpleRsp.IMethodResponse.Error {
			e := cim.Message.SimpleRsp.IMethodResponse.Error
			return WBEMException(CIMStatusCode(e.Code), e.Description)
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse.ReturnValue {
			return ireturnValueNotExists
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse.ReturnValue.ValueNamedInstances {
			return valueNamedInstancesNotExists
		}
		return nil
	}}

	// CIMProtocolVersion: 1.0
	// CIMOperation: MethodCall
	// CIMMethod: EnumerateClassNames
	// CIMObject: root%2Fcimv2

	if err := c.RoundTrip("POST", map[string]string{"CIMProtocolVersion": c.ProtocolVersion,
		"CIMOperation": "MethodCall",
		"CIMMethod":    "EnumerateInstances",
		"CIMObject":    url.QueryEscape(namespaceName)}, req, resp); nil != err {
		return nil, err
	}

	results := make([]CIMInstanceWithName, len(resp.Message.SimpleRsp.IMethodResponse.ReturnValue.ValueNamedInstances))
	for idx, _ := range resp.Message.SimpleRsp.IMethodResponse.ReturnValue.ValueNamedInstances {
		results[idx] = &resp.Message.SimpleRsp.IMethodResponse.ReturnValue.ValueNamedInstances[idx]
	}
	return results, nil
}

func (c *ClientCIMXML) GetClass(namespaceName string, className string, localOnly bool,
	includeQualifiers bool, includeClassOrigin bool, propertyList []string) (string, error) {
	if "" == namespaceName {
		return "", WBEMException(CIM_ERR_INVALID_PARAMETER,
			"namespace name is empty.")
	}

	if "" == className {
		return "", WBEMException(CIM_ERR_INVALID_PARAMETER,
			"class name is empty.")
	}

	names := strings.Split(namespaceName, "/")
	namespaces := make([]CimNamespace, len(names))
	for idx, name := range names {
		namespaces[idx].Name = name
	}

	paramValues := []CimIParamValue{
		CimIParamValue{
			Name:      "ClassName",
			ClassName: &CimClassName{Name: className},
		},

		CimIParamValue{
			Name:  "LocalOnly",
			Value: &CimValue{Value: booleanString(localOnly)},
		},

		CimIParamValue{
			Name:  "IncludeQualifiers",
			Value: &CimValue{Value: booleanString(includeQualifiers)},
		},
		CimIParamValue{
			Name:  "IncludeClassOrigin",
			Value: &CimValue{Value: booleanString(includeClassOrigin)},
		},
	}
	if 0 != len(propertyList) {
		properties := CimValueArray(make([]CimValueOrNull, len(propertyList)))
		for idx, s := range propertyList {
			properties[idx] = CimValueOrNull{Value: &CimValue{Value: s}}
		}
		paramValues = append(paramValues,
			CimIParamValue{
				Name:       "PropertyList",
				ValueArray: properties,
			})
	}

	simpleReq := &CimSimpleReq{IMethodCall: &CimIMethodCall{
		Name:               "GetClass",
		LocalNamespacePath: CimLocalNamespacePath{Namespaces: namespaces},
		ParamValues:        paramValues,
	}}

	req := &CIM{
		CimVersion: c.CimVersion,
		DtdVersion: c.DtdVersion,
		Message: &CimMessage{
			Id:              c.generateId(),
			ProtocolVersion: c.ProtocolVersion,
			SimpleReq:       simpleReq,
		},
		//Declaration: &CimDeclaration,
	}

	resp := &CIM{hasFault: func(cim *CIM) error {
		if nil == cim.Message {
			return messageNotExists
		}
		if nil == cim.Message.SimpleRsp {
			return simpleReqNotExists
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse {
			return imethodResponseNotExists
		}
		if nil != cim.Message.SimpleRsp.IMethodResponse.Error {
			e := cim.Message.SimpleRsp.IMethodResponse.Error
			return WBEMException(CIMStatusCode(e.Code), e.Description)
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse.ReturnValue {
			return ireturnValueNotExists
		}
		if 0 == len(cim.Message.SimpleRsp.IMethodResponse.ReturnValue.Classes) {
			return classesNotExists
		}

		if 1 < len(cim.Message.SimpleRsp.IMethodResponse.ReturnValue.Classes) {
			return classesMutiChioce
		}
		return nil
	}}

	// CIMProtocolVersion: 1.0
	// CIMOperation: MethodCall
	// CIMMethod: EnumerateClassNames
	// CIMObject: root%2Fcimv2

	if err := c.RoundTrip("POST", map[string]string{"CIMProtocolVersion": c.ProtocolVersion,
		"CIMOperation": "MethodCall",
		"CIMMethod":    "GetClass",
		"CIMObject":    url.QueryEscape(namespaceName)}, req, resp); nil != err {
		return "", err
	}

	return resp.Message.SimpleRsp.IMethodResponse.ReturnValue.Classes[0].String(), nil
}

func (c *ClientCIMXML) EnumerateClasses(namespaceName string, className string, deepInheritance bool,
	localOnly bool, includeQualifiers bool, includeClassOrigin bool) ([]string, error) {
	if "" == namespaceName {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"namespace name is empty.")
	}

	names := strings.Split(namespaceName, "/")
	namespaces := make([]CimNamespace, len(names))
	for idx, name := range names {
		namespaces[idx].Name = name
	}

	paramValues := []CimIParamValue{
		CimIParamValue{
			Name:  "LocalOnly",
			Value: &CimValue{Value: booleanString(localOnly)},
		},
		CimIParamValue{
			Name:  "DeepInheritance",
			Value: &CimValue{Value: booleanString(deepInheritance)},
		},

		CimIParamValue{
			Name:  "IncludeQualifiers",
			Value: &CimValue{Value: booleanString(includeQualifiers)},
		},
		CimIParamValue{
			Name:  "IncludeClassOrigin",
			Value: &CimValue{Value: booleanString(includeClassOrigin)},
		},
	}

	if "" != className {
		paramValues = append(paramValues, CimIParamValue{
			Name:      "ClassName",
			ClassName: &CimClassName{Name: className},
		})
	}

	simpleReq := &CimSimpleReq{IMethodCall: &CimIMethodCall{
		Name:               "EnumerateClasses",
		LocalNamespacePath: CimLocalNamespacePath{Namespaces: namespaces},
		ParamValues:        paramValues,
	}}

	req := &CIM{
		CimVersion: c.CimVersion,
		DtdVersion: c.DtdVersion,
		Message: &CimMessage{
			Id:              c.generateId(),
			ProtocolVersion: c.ProtocolVersion,
			SimpleReq:       simpleReq,
		},
		//Declaration: &CimDeclaration,
	}

	resp := &CIM{hasFault: func(cim *CIM) error {
		if nil == cim.Message {
			return messageNotExists
		}
		if nil == cim.Message.SimpleRsp {
			return simpleReqNotExists
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse {
			return imethodResponseNotExists
		}
		if nil != cim.Message.SimpleRsp.IMethodResponse.Error {
			e := cim.Message.SimpleRsp.IMethodResponse.Error
			return WBEMException(CIMStatusCode(e.Code), e.Description)
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse.ReturnValue {
			return ireturnValueNotExists
		}
		// if 0 == len(cim.Message.SimpleRsp.IMethodResponse.ReturnValue.Classes) &&
		//    0 == len(resp.Message.SimpleRsp.IMethodResponse.ReturnValue.ClassNames) {
		// 	return classesNotExists
		// }
		return nil
	}}

	// CIMProtocolVersion: 1.0
	// CIMOperation: MethodCall
	// CIMMethod: EnumerateClassNames
	// CIMObject: root%2Fcimv2

	if err := c.RoundTrip("POST", map[string]string{"CIMProtocolVersion": c.ProtocolVersion,
		"CIMOperation": "MethodCall",
		"CIMMethod":    "EnumerateClasses",
		"CIMObject":    url.QueryEscape(namespaceName)}, req, resp); nil != err {
		return nil, err
	}

	results := make([]string, len(resp.Message.SimpleRsp.IMethodResponse.ReturnValue.Classes))
	for idx, class := range resp.Message.SimpleRsp.IMethodResponse.ReturnValue.Classes {
		results[idx] = class.String()
	}
	for _, name := range resp.Message.SimpleRsp.IMethodResponse.ReturnValue.ClassNames {
		results = append(results, name.Name)
	}
	return results, nil
}

func (c *ClientCIMXML) AssociatorNames(namespaceName string, instanceName CIMInstanceName,
	assocClass, resultClass, role, resultRole string) ([]CIMInstanceName, error) {

	if "" == namespaceName {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"namespace name is empty.")
	}

	if "" == instanceName.GetClassName() {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"class name is empty.")
	}

	names := strings.Split(namespaceName, "/")
	namespaces := make([]CimNamespace, len(names))
	for idx, name := range names {
		namespaces[idx].Name = name
	}

	paramValues := []CimIParamValue{
		CimIParamValue{
			Name:         "ObjectName",
			InstanceName: instanceName.(*CimInstanceName),
		},
	}

	if "" != assocClass {
		paramValues = append(paramValues, CimIParamValue{
			Name:      "AssocClass",
			ClassName: &CimClassName{Name: assocClass},
		})
	}

	if "" != resultClass {
		paramValues = append(paramValues, CimIParamValue{
			Name:      "ResultClass",
			ClassName: &CimClassName{Name: resultClass},
		})
	}

	if "" != role {
		paramValues = append(paramValues, CimIParamValue{
			Name:  "Role",
			Value: &CimValue{Value: role},
		})
	}

	if "" != resultRole {
		paramValues = append(paramValues, CimIParamValue{
			Name:  "ResultRole",
			Value: &CimValue{Value: resultRole},
		})
	}

	simpleReq := &CimSimpleReq{IMethodCall: &CimIMethodCall{
		Name:               "AssociatorNames",
		LocalNamespacePath: CimLocalNamespacePath{Namespaces: namespaces},
		ParamValues:        paramValues,
	}}

	req := &CIM{
		CimVersion: c.CimVersion,
		DtdVersion: c.DtdVersion,
		Message: &CimMessage{
			Id:              c.generateId(),
			ProtocolVersion: c.ProtocolVersion,
			SimpleReq:       simpleReq,
		},
		//Declaration: &CimDeclaration,
	}

	resp := &CIM{hasFault: func(cim *CIM) error {
		if nil == cim.Message {
			return messageNotExists
		}
		if nil == cim.Message.SimpleRsp {
			return simpleReqNotExists
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse {
			return imethodResponseNotExists
		}
		if nil != cim.Message.SimpleRsp.IMethodResponse.Error {
			e := cim.Message.SimpleRsp.IMethodResponse.Error
			return WBEMException(CIMStatusCode(e.Code), e.Description)
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse.ReturnValue {
			return ireturnValueNotExists
		}
		// if 0 == len(cim.Message.SimpleRsp.IMethodResponse.ReturnValue.InstanceNames) {
		// 	return InstanceNamesNotExists
		// }
		return nil
	}}

	// CIMProtocolVersion: 1.0
	// CIMOperation: MethodCall
	// CIMMethod: EnumerateClassNames
	// CIMObject: root%2Fcimv2

	if err := c.RoundTrip("POST", map[string]string{"CIMProtocolVersion": c.ProtocolVersion,
		"CIMOperation": "MethodCall",
		"CIMMethod":    "AssociatorNames",
		"CIMObject":    url.QueryEscape(namespaceName)}, req, resp); nil != err {
		return nil, err
	}

	results := make([]CIMInstanceName, len(resp.Message.SimpleRsp.IMethodResponse.ReturnValue.InstanceNames))
	for idx, name := range resp.Message.SimpleRsp.IMethodResponse.ReturnValue.InstanceNames {
		results[idx] = name
	}
	return results, nil
}

func (c *ClientCIMXML) AssociatorInstances(namespaceName string, instanceName CIMInstanceName,
	assocClass, resultClass, role, resultRole string, includeClassOrigin bool, propertyList []string) ([]CIMInstance, error) {
	if "" == namespaceName {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"namespace name is empty.")
	}

	if "" == instanceName.GetClassName() {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"class name is empty.")
	}

	names := strings.Split(namespaceName, "/")
	namespaces := make([]CimNamespace, len(names))
	for idx, name := range names {
		namespaces[idx].Name = name
	}

	paramValues := []CimIParamValue{
		CimIParamValue{
			Name:         "ObjectName",
			InstanceName: instanceName.(*CimInstanceName),
		},
		CimIParamValue{
			Name:  "IncludeClassOrigin",
			Value: &CimValue{Value: booleanString(includeClassOrigin)},
		},
	}

	if "" != assocClass {
		paramValues = append(paramValues, CimIParamValue{
			Name:      "AssocClass",
			ClassName: &CimClassName{Name: assocClass},
		})
	}

	if "" != resultClass {
		paramValues = append(paramValues, CimIParamValue{
			Name:      "ResultClass",
			ClassName: &CimClassName{Name: resultClass},
		})
	}

	if "" != role {
		paramValues = append(paramValues, CimIParamValue{
			Name:  "Role",
			Value: &CimValue{Value: role},
		})
	}

	if "" != resultRole {
		paramValues = append(paramValues, CimIParamValue{
			Name:  "ResultRole",
			Value: &CimValue{Value: resultRole},
		})
	}

	if 0 != len(propertyList) {
		properties := CimValueArray(make([]CimValueOrNull, len(propertyList)))
		for idx, s := range propertyList {
			properties[idx] = CimValueOrNull{Value: &CimValue{Value: s}}
		}
		paramValues = append(paramValues,
			CimIParamValue{
				Name:       "PropertyList",
				ValueArray: properties,
			})
	}

	simpleReq := &CimSimpleReq{IMethodCall: &CimIMethodCall{
		Name:               "Associators",
		LocalNamespacePath: CimLocalNamespacePath{Namespaces: namespaces},
		ParamValues:        paramValues,
	}}

	req := &CIM{
		CimVersion: c.CimVersion,
		DtdVersion: c.DtdVersion,
		Message: &CimMessage{
			Id:              c.generateId(),
			ProtocolVersion: c.ProtocolVersion,
			SimpleReq:       simpleReq,
		},
		//Declaration: &CimDeclaration,
	}

	resp := &CIM{hasFault: func(cim *CIM) error {
		if nil == cim.Message {
			return messageNotExists
		}
		if nil == cim.Message.SimpleRsp {
			return simpleReqNotExists
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse {
			return imethodResponseNotExists
		}
		if nil != cim.Message.SimpleRsp.IMethodResponse.Error {
			e := cim.Message.SimpleRsp.IMethodResponse.Error
			return WBEMException(CIMStatusCode(e.Code), e.Description)
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse.ReturnValue {
			return ireturnValueNotExists
		}
		// if 0 == len(cim.Message.SimpleRsp.IMethodResponse.ReturnValue.Instances) {
		// 	return classesNotExists
		// }
		return nil
	}}

	// CIMProtocolVersion: 1.0
	// CIMOperation: MethodCall
	// CIMMethod: EnumerateClassNames
	// CIMObject: root%2Fcimv2

	if err := c.RoundTrip("POST", map[string]string{"CIMProtocolVersion": c.ProtocolVersion,
		"CIMOperation": "MethodCall",
		"CIMMethod":    "Associators",
		"CIMObject":    url.QueryEscape(namespaceName)}, req, resp); nil != err {
		return nil, err
	}

	results := make([]CIMInstance, len(resp.Message.SimpleRsp.IMethodResponse.ReturnValue.Instances))
	for idx, _ := range resp.Message.SimpleRsp.IMethodResponse.ReturnValue.Instances {
		results[idx] = &resp.Message.SimpleRsp.IMethodResponse.ReturnValue.Instances[idx]
	}
	return results, nil
}

func (c *ClientCIMXML) AssociatorClasses(namespaceName, className, assocClass, resultClass, role, resultRole string,
	includeQualifiers, includeClassOrigin bool, propertyList []string) ([]string, error) {
	if "" == namespaceName {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"namespace name is empty.")
	}

	if "" == className {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"class name is empty.")
	}

	names := strings.Split(namespaceName, "/")
	namespaces := make([]CimNamespace, len(names))
	for idx, name := range names {
		namespaces[idx].Name = name
	}

	paramValues := []CimIParamValue{
		CimIParamValue{
			Name:         "ObjectName",
			InstanceName: &CimInstanceName{ClassName: className},
		},
		CimIParamValue{
			Name:  "IncludeQualifiers",
			Value: &CimValue{Value: booleanString(includeQualifiers)},
		},
		CimIParamValue{
			Name:  "IncludeClassOrigin",
			Value: &CimValue{Value: booleanString(includeClassOrigin)},
		},
	}

	if "" != assocClass {
		paramValues = append(paramValues, CimIParamValue{
			Name:      "AssocClass",
			ClassName: &CimClassName{Name: assocClass},
		})
	}

	if "" != resultClass {
		paramValues = append(paramValues, CimIParamValue{
			Name:      "ResultClass",
			ClassName: &CimClassName{Name: resultClass},
		})
	}

	if "" != role {
		paramValues = append(paramValues, CimIParamValue{
			Name:  "Role",
			Value: &CimValue{Value: role},
		})
	}

	if "" != resultRole {
		paramValues = append(paramValues, CimIParamValue{
			Name:  "ResultRole",
			Value: &CimValue{Value: resultRole},
		})
	}

	if 0 != len(propertyList) {
		properties := CimValueArray(make([]CimValueOrNull, len(propertyList)))
		for idx, s := range propertyList {
			properties[idx] = CimValueOrNull{Value: &CimValue{Value: s}}
		}
		paramValues = append(paramValues,
			CimIParamValue{
				Name:       "PropertyList",
				ValueArray: properties,
			})
	}

	simpleReq := &CimSimpleReq{IMethodCall: &CimIMethodCall{
		Name:               "Associators",
		LocalNamespacePath: CimLocalNamespacePath{Namespaces: namespaces},
		ParamValues:        paramValues,
	}}

	req := &CIM{
		CimVersion: c.CimVersion,
		DtdVersion: c.DtdVersion,
		Message: &CimMessage{
			Id:              c.generateId(),
			ProtocolVersion: c.ProtocolVersion,
			SimpleReq:       simpleReq,
		},
		//Declaration: &CimDeclaration,
	}

	resp := &CIM{hasFault: func(cim *CIM) error {
		if nil == cim.Message {
			return messageNotExists
		}
		if nil == cim.Message.SimpleRsp {
			return simpleReqNotExists
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse {
			return imethodResponseNotExists
		}
		if nil != cim.Message.SimpleRsp.IMethodResponse.Error {
			e := cim.Message.SimpleRsp.IMethodResponse.Error
			return WBEMException(CIMStatusCode(e.Code), e.Description)
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse.ReturnValue {
			return ireturnValueNotExists
		}
		// if 0 == len(cim.Message.SimpleRsp.IMethodResponse.ReturnValue.Classes) {
		// 	return classesNotExists
		// }
		return nil
	}}

	// CIMProtocolVersion: 1.0
	// CIMOperation: MethodCall
	// CIMMethod: Associators
	// CIMObject: root%2Fcimv2

	if err := c.RoundTrip("POST", map[string]string{"CIMProtocolVersion": c.ProtocolVersion,
		"CIMOperation": "MethodCall",
		"CIMMethod":    "Associators",
		"CIMObject":    url.QueryEscape(namespaceName)}, req, resp); nil != err {
		return nil, err
	}

	results := make([]string, len(resp.Message.SimpleRsp.IMethodResponse.ReturnValue.Classes))
	for idx, class := range resp.Message.SimpleRsp.IMethodResponse.ReturnValue.Classes {
		results[idx] = class.String()
	}

	for _, name := range resp.Message.SimpleRsp.IMethodResponse.ReturnValue.ClassNames {
		results = append(results, name.Name)
	}
	return results, nil
}

func (c *ClientCIMXML) ReferenceNames(namespaceName string, instanceName CIMInstanceName,
	resultClass, role string) ([]CIMInstanceName, error) {

	if "" == namespaceName {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"namespace name is empty.")
	}

	if "" == instanceName.GetClassName() {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"class name is empty.")
	}

	if 0 == instanceName.GetKeyBindings().Len() {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"key bindings is empty.")
	}

	names := strings.Split(namespaceName, "/")
	namespaces := make([]CimNamespace, len(names))
	for idx, name := range names {
		namespaces[idx].Name = name
	}

	paramValues := []CimIParamValue{
		CimIParamValue{
			Name:         "ObjectName",
			InstanceName: instanceName.(*CimInstanceName),
		},
	}

	if "" != resultClass {
		paramValues = append(paramValues, CimIParamValue{
			Name:      "ResultClass",
			ClassName: &CimClassName{Name: resultClass},
		})
	}

	if "" != role {
		paramValues = append(paramValues, CimIParamValue{
			Name:  "Role",
			Value: &CimValue{Value: role},
		})
	}

	simpleReq := &CimSimpleReq{IMethodCall: &CimIMethodCall{
		Name:               "ReferenceNames",
		LocalNamespacePath: CimLocalNamespacePath{Namespaces: namespaces},
		ParamValues:        paramValues,
	}}

	req := &CIM{
		CimVersion: c.CimVersion,
		DtdVersion: c.DtdVersion,
		Message: &CimMessage{
			Id:              c.generateId(),
			ProtocolVersion: c.ProtocolVersion,
			SimpleReq:       simpleReq,
		},
		//Declaration: &CimDeclaration,
	}

	resp := &CIM{hasFault: func(cim *CIM) error {
		if nil == cim.Message {
			return messageNotExists
		}
		if nil == cim.Message.SimpleRsp {
			return simpleReqNotExists
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse {
			return imethodResponseNotExists
		}
		if nil != cim.Message.SimpleRsp.IMethodResponse.Error {
			e := cim.Message.SimpleRsp.IMethodResponse.Error
			return WBEMException(CIMStatusCode(e.Code), e.Description)
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse.ReturnValue {
			return ireturnValueNotExists
		}
		// if 0 == len(cim.Message.SimpleRsp.IMethodResponse.ReturnValue.InstanceNames) {
		// 	return InstanceNamesNotExists
		// }
		return nil
	}}

	// CIMProtocolVersion: 1.0
	// CIMOperation: MethodCall
	// CIMMethod: ReferenceNames
	// CIMObject: root%2Fcimv2

	if err := c.RoundTrip("POST", map[string]string{"CIMProtocolVersion": c.ProtocolVersion,
		"CIMOperation": "MethodCall",
		"CIMMethod":    "ReferenceNames",
		"CIMObject":    url.QueryEscape(namespaceName)}, req, resp); nil != err {
		return nil, err
	}

	results := make([]CIMInstanceName, len(resp.Message.SimpleRsp.IMethodResponse.ReturnValue.InstanceNames))
	for idx, name := range resp.Message.SimpleRsp.IMethodResponse.ReturnValue.InstanceNames {
		results[idx] = name
	}
	return results, nil
}

func (c *ClientCIMXML) ReferenceInstances(namespaceName string, instanceName CIMInstanceName,
	resultClass, role string, includeClassOrigin bool, propertyList []string) ([]CIMInstance, error) {
	if "" == namespaceName {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"namespace name is empty.")
	}

	if "" == instanceName.GetClassName() {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"class name is empty.")
	}

	names := strings.Split(namespaceName, "/")
	namespaces := make([]CimNamespace, len(names))
	for idx, name := range names {
		namespaces[idx].Name = name
	}

	paramValues := []CimIParamValue{
		CimIParamValue{
			Name:         "ObjectName",
			InstanceName: instanceName.(*CimInstanceName),
		},
		CimIParamValue{
			Name:  "IncludeClassOrigin",
			Value: &CimValue{Value: booleanString(includeClassOrigin)},
		},
	}

	if "" != resultClass {
		paramValues = append(paramValues, CimIParamValue{
			Name:      "ResultClass",
			ClassName: &CimClassName{Name: resultClass},
		})
	}

	if "" != role {
		paramValues = append(paramValues, CimIParamValue{
			Name:  "Role",
			Value: &CimValue{Value: role},
		})
	}

	if 0 != len(propertyList) {
		properties := CimValueArray(make([]CimValueOrNull, len(propertyList)))
		for idx, s := range propertyList {
			properties[idx] = CimValueOrNull{Value: &CimValue{Value: s}}
		}
		paramValues = append(paramValues,
			CimIParamValue{
				Name:       "PropertyList",
				ValueArray: properties,
			})
	}

	simpleReq := &CimSimpleReq{IMethodCall: &CimIMethodCall{
		Name:               "References",
		LocalNamespacePath: CimLocalNamespacePath{Namespaces: namespaces},
		ParamValues:        paramValues,
	}}

	req := &CIM{
		CimVersion: c.CimVersion,
		DtdVersion: c.DtdVersion,
		Message: &CimMessage{
			Id:              c.generateId(),
			ProtocolVersion: c.ProtocolVersion,
			SimpleReq:       simpleReq,
		},
		//Declaration: &CimDeclaration,
	}

	resp := &CIM{hasFault: func(cim *CIM) error {
		if nil == cim.Message {
			return messageNotExists
		}
		if nil == cim.Message.SimpleRsp {
			return simpleReqNotExists
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse {
			return imethodResponseNotExists
		}
		if nil != cim.Message.SimpleRsp.IMethodResponse.Error {
			e := cim.Message.SimpleRsp.IMethodResponse.Error
			return WBEMException(CIMStatusCode(e.Code), e.Description)
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse.ReturnValue {
			return ireturnValueNotExists
		}
		// if 0 == len(cim.Message.SimpleRsp.IMethodResponse.ReturnValue.Instances) {
		// 	return classesNotExists
		// }
		return nil
	}}

	// CIMProtocolVersion: 1.0
	// CIMOperation: MethodCall
	// CIMMethod: EnumerateClassNames
	// CIMObject: root%2Fcimv2

	if err := c.RoundTrip("POST", map[string]string{"CIMProtocolVersion": c.ProtocolVersion,
		"CIMOperation": "MethodCall",
		"CIMMethod":    "References",
		"CIMObject":    url.QueryEscape(namespaceName)}, req, resp); nil != err {
		return nil, err
	}

	results := make([]CIMInstance, len(resp.Message.SimpleRsp.IMethodResponse.ReturnValue.Instances))
	for idx, _ := range resp.Message.SimpleRsp.IMethodResponse.ReturnValue.Instances {
		results[idx] = &resp.Message.SimpleRsp.IMethodResponse.ReturnValue.Instances[idx]
	}
	return results, nil
}

func (c *ClientCIMXML) ReferenceClasses(namespaceName, className, resultClass, role string,
	includeQualifiers, includeClassOrigin bool, propertyList []string) ([]string, error) {
	if "" == namespaceName {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"namespace name is empty.")
	}

	if "" == className {
		return nil, WBEMException(CIM_ERR_INVALID_PARAMETER,
			"class name is empty.")
	}

	names := strings.Split(namespaceName, "/")
	namespaces := make([]CimNamespace, len(names))
	for idx, name := range names {
		namespaces[idx].Name = name
	}

	paramValues := []CimIParamValue{
		CimIParamValue{
			Name:      "ObjectName",
			ClassName: &CimClassName{Name: className},
		},
		CimIParamValue{
			Name:  "IncludeQualifiers",
			Value: &CimValue{Value: booleanString(includeQualifiers)},
		},
		CimIParamValue{
			Name:  "IncludeClassOrigin",
			Value: &CimValue{Value: booleanString(includeClassOrigin)},
		},
	}

	if "" != resultClass {
		paramValues = append(paramValues, CimIParamValue{
			Name:      "ResultClass",
			ClassName: &CimClassName{Name: resultClass},
		})
	}

	if "" != role {
		paramValues = append(paramValues, CimIParamValue{
			Name:  "Role",
			Value: &CimValue{Value: role},
		})
	}

	if 0 != len(propertyList) {
		properties := CimValueArray(make([]CimValueOrNull, len(propertyList)))
		for idx, s := range propertyList {
			properties[idx] = CimValueOrNull{Value: &CimValue{Value: s}}
		}
		paramValues = append(paramValues,
			CimIParamValue{
				Name:       "PropertyList",
				ValueArray: properties,
			})
	}

	simpleReq := &CimSimpleReq{IMethodCall: &CimIMethodCall{
		Name:               "References",
		LocalNamespacePath: CimLocalNamespacePath{Namespaces: namespaces},
		ParamValues:        paramValues,
	}}

	req := &CIM{
		CimVersion: c.CimVersion,
		DtdVersion: c.DtdVersion,
		Message: &CimMessage{
			Id:              c.generateId(),
			ProtocolVersion: c.ProtocolVersion,
			SimpleReq:       simpleReq,
		},
		//Declaration: &CimDeclaration,
	}

	resp := &CIM{hasFault: func(cim *CIM) error {
		if nil == cim.Message {
			return messageNotExists
		}
		if nil == cim.Message.SimpleRsp {
			return simpleReqNotExists
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse {
			return imethodResponseNotExists
		}
		if nil != cim.Message.SimpleRsp.IMethodResponse.Error {
			e := cim.Message.SimpleRsp.IMethodResponse.Error
			return WBEMException(CIMStatusCode(e.Code), e.Description)
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse.ReturnValue {
			return ireturnValueNotExists
		}
		// if 0 == len(cim.Message.SimpleRsp.IMethodResponse.ReturnValue.Classes)
		//    0 == len(cim.Message.SimpleRsp.IMethodResponse.ReturnValue.ClassNames) {
		// 	return classesNotExists
		// }
		return nil
	}}

	// CIMProtocolVersion: 1.0
	// CIMOperation: MethodCall
	// CIMMethod: References
	// CIMObject: root%2Fcimv2

	if err := c.RoundTrip("POST", map[string]string{"CIMProtocolVersion": c.ProtocolVersion,
		"CIMOperation": "MethodCall",
		"CIMMethod":    "References",
		"CIMObject":    url.QueryEscape(namespaceName)}, req, resp); nil != err {
		return nil, err
	}

	results := make([]string, len(resp.Message.SimpleRsp.IMethodResponse.ReturnValue.Classes))
	for idx, class := range resp.Message.SimpleRsp.IMethodResponse.ReturnValue.Classes {
		results[idx] = class.String()
	}
	for _, name := range resp.Message.SimpleRsp.IMethodResponse.ReturnValue.ClassNames {
		results = append(results, name.Name)
	}
	return results, nil
}

func NewClientCIMXML(u *url.URL, insecure bool) (*ClientCIMXML, error) {
	c := &ClientCIMXML{}
	c.init(u, insecure)
	return c, nil
}
