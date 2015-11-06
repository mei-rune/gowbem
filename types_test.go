package gowbem

import (
	"encoding/xml"
	"reflect"
	"strings"
	"testing"

	"github.com/aryann/difflib"
)

func makeLocalNamespace(ss []string) *CimLocalNamespacePath {
	names := make([]CimNamespace, 0, len(ss))
	for _, s := range ss {
		names = append(names, CimNamespace{Name: s})
	}
	return &CimLocalNamespacePath{Namespaces: names}
}

func makeLocalClass(ss []string, cls string) *CimLocalClassPath {
	return &CimLocalClassPath{NamespacePath: *makeLocalNamespace(ss),
		ClassName: CimClassName{Name: cls}}
}

func makeNamespace(host string, ss []string) *CimNamespacePath {
	return &CimNamespacePath{
		Host:               CimHost{Value: host},
		LocalNamespacePath: *makeLocalNamespace(ss),
	}
}
func makeClass(host string, ss []string, cls string) *CimClassPath {
	return &CimClassPath{NamespacePath: *makeNamespace(host, ss),
		ClassName: CimClassName{Name: cls}}
}

func makeInstanceNameWithKV1(cls, k, v, vt string) *CimInstanceName {
	return &CimInstanceName{
		ClassName: cls,
		KeyValue: &CimKeyValue{
			ValueType: vt,
			Value:     v,
		}}
}

func makeInstanceNameWithKV2(cls, k, v, vt string) *CimInstanceName {
	return &CimInstanceName{
		ClassName: cls,
		KeyValue: &CimKeyValue{
			Type:  vt,
			Value: v,
		}}
}

func makeInstanceNameWithValueRef(cls string, ref *CimValueReference) *CimInstanceName {
	return &CimInstanceName{
		ClassName:      cls,
		ValueReference: ref}
}

func makeValueRef(factor string) *CimValueReference {
	return &CimValueReference{ClassPath: makeClass("192.168.1.23", []string{"a", "b", factor}, "value"+factor)}
}

func makeQualifier(factor string) CimQualifier {
	return CimQualifier{
		Name:       "CimQualifier_" + factor,
		Type:       "string",
		Propagated: true,
		Lang:       "cn",
		Value:      &CimValue{Value: "abc_" + factor},
	}
}

func makeValueArray(values ...CimValueOrNull) CimValueArray {
	return CimValueArray(values)
}
func makeValueRefArray(values ...CimValueReferenceOrNull) CimValueRefArray {
	return CimValueRefArray(values)
}

func makeProperty(factor string) *CimAnyProperty {
	return &CimAnyProperty{
		Property: &CimProperty{
			Name:           "pr_" + factor,
			Type:           "string",
			ClassOrigin:    "class_origin" + factor,
			Propagated:     true,
			EmbeddedObject: "object",
			Lang:           "zh",
			Qualifiers:     []CimQualifier{makeQualifier("pr" + factor + "q1"), makeQualifier("pr" + factor + "q2")},
			Value:          &CimValue{Value: "value" + factor},
		},
	}
}

func makeInstance(factor string) *CimInstance {
	return &CimInstance{
		ClassName:  "abc_class" + factor,
		Lang:       "zh",
		Qualifiers: []CimQualifier{makeQualifier("i1"), makeQualifier("i2"), makeQualifier("i3")},
		Properties: []CimAnyProperty{
			*makeProperty(factor + "_1"),
			*makeProperty(factor + "_2"),
			*makeProperty(factor + "_3"),
		},
	}
}

var paramValues = []CimParamValue{

	CimParamValue{
		Name:           "p1",
		ParamType:      "string",
		EmbeddedObject: "instance",
		Value:          &CimValue{Value: "value1"},
	},

	CimParamValue{
		Name:           "p2",
		ParamType:      "string",
		EmbeddedObject: "instance",
		ValueReference: &CimValueReference{
			ClassPath: makeClass("127.9.2.1", []string{"a", "bc"}, "test_class1"),
		},
	},

	CimParamValue{
		Name:           "p3",
		ParamType:      "string",
		EmbeddedObject: "instance",
		ValueArray: makeValueArray(
			CimValueOrNull{Value: &CimValue{Value: "abc"}},
			CimValueOrNull{Null: &CimValueNull{}},
			CimValueOrNull{Value: &CimValue{Value: "abc"}},
		),
	},

	CimParamValue{
		Name:           "p4",
		ParamType:      "string",
		EmbeddedObject: "instance",
		ValueRefArray: makeValueRefArray(
			CimValueReferenceOrNull{Value: &CimValueReference{ClassPath: makeClass("127.9.2.1", []string{"a", "bc"}, "test_class1")}},
			CimValueReferenceOrNull{Null: &CimValueNull{}},
			CimValueReferenceOrNull{Value: &CimValueReference{
				InstanceName: &CimInstanceName{
					ClassName: "abc_test",
					KeyBindings: []CimKeyBinding{
						CimKeyBinding{
							Name:     "kb1",
							KeyValue: &CimKeyValue{Type: "string", Value: "kb_value_34"},
						},
					},
				},
			},
			},
		),
	},

	CimParamValue{
		Name:           "p5",
		ParamType:      "string",
		EmbeddedObject: "instance",
		ClassName:      &CimClassName{Name: "a.b.class_test_p5"},
	},

	CimParamValue{
		Name:           "p6",
		ParamType:      "string",
		EmbeddedObject: "instance",
		InstanceName:   makeInstanceNameWithKV1("a.b.c.class_test_p6", "k_p6", "v_p6", "string"),
	},

	CimParamValue{
		Name:           "p7",
		ParamType:      "string",
		EmbeddedObject: "instance",
		Class: &CimClass{
			Name:       "a.b.c.class_test_p7",
			SuperClass: "a.b.c.class_test_p7_base",
			Qualifiers: []CimQualifier{CimQualifier{

				Name:       "CimQualifier_p7_1",
				Type:       "abc",
				Propagated: true,
				Lang:       "cn",
				Value:      &CimValue{Value: "abc"},
				CimQualifierFlavor: CimQualifierFlavor{Overridable: true,
					ToSubclass:   true,
					ToInstance:   true,
					Translatable: true},
			},

				CimQualifier{
					Name:       "CimQualifier_p7_2",
					Type:       "abc",
					Propagated: true,
					Lang:       "cn",
					ValueArray: makeValueArray(
						CimValueOrNull{Value: &CimValue{Value: "abc"}},
					),
				}},
			Properties: []CimAnyProperty{
				CimAnyProperty{Property: &CimProperty{
					Name:           "pr_p7_1_abc_1",
					Type:           "string",
					ClassOrigin:    "pr_p7_1_abc_1_origin",
					Propagated:     true,
					EmbeddedObject: "object",
					Lang:           "cn",
					Value:          &CimValue{Value: "vvvvv"},
				}},
				CimAnyProperty{PropertyArray: &CimPropertyArray{
					Name:           "pr_p7_1_abc_2",
					Type:           "string",
					ArraySize:      23,
					ClassOrigin:    "pr_p7_1_abc_2_origin",
					Propagated:     true,
					EmbeddedObject: "object",
					Lang:           "cn",
					ValueArray: makeValueArray(
						CimValueOrNull{Value: &CimValue{Value: "vvvvv"}},
					),
				}},
				CimAnyProperty{PropertyReference: &CimPropertyReference{
					Name:           "pr_p7_1_abc_3",
					ReferenceClass: "ref_class",
					ClassOrigin:    "pr_p7_1_abc_3_origin",

					Propagated: true,
					Qualifiers: []CimQualifier{CimQualifier{
						Name:       "CimQualifier_p7_q_1",
						Type:       "abc",
						Propagated: true,
						Lang:       "cn",
						Value:      &CimValue{Value: "abc"},
						CimQualifierFlavor: CimQualifierFlavor{Overridable: true,
							ToSubclass:   true,
							ToInstance:   true,
							Translatable: true},
					},

						CimQualifier{
							Name:       "CimQualifier_p7_q_1",
							Type:       "abc",
							Propagated: true,
							Lang:       "cn",
							ValueArray: makeValueArray(
								CimValueOrNull{Value: &CimValue{Value: "abc"}},
							),
						}},
					ValueReference: &CimValueReference{
						LocalInstancePath: &CimLocalInstancePath{
							LocalNamespacePath: *makeLocalNamespace([]string{"f", "t"}),
							InstanceName:       *makeInstanceNameWithKV2("F.4.A.class_abc", "ref_class", "abcss", "bool"),
						},
					},
				}},
			},
			Methods: []CimMethod{
				CimMethod{
					Name:        "method_1",
					Type:        "string",
					ClassOrigin: "method_1_origin",
					Propagated:  true,
					Qualifiers: []CimQualifier{
						makeQualifier("m1_q_1"),
						makeQualifier("m1_q_2"),
					},
					Parameters: []CimAnyParameter{
						CimAnyParameter{Parameter: &CimParameter{
							Name: "method_1_p1",
							Type: "string",
							Qualifiers: []CimQualifier{
								makeQualifier("m1_q_p1_1"),
								makeQualifier("m1_q_p1_2"),
							}}},
						CimAnyParameter{Parameter: &CimParameter{
							Name: "method_1_p2",
							Type: "string"}},
						CimAnyParameter{ParameterReference: &CimParameterReference{
							Name:           "method_1_p3",
							ReferenceClass: "string",
							Qualifiers: []CimQualifier{
								makeQualifier("m1_q_p3_1"),
								makeQualifier("m1_q_p3_2"),
							}}},
						CimAnyParameter{ParameterReference: &CimParameterReference{
							Name:           "method_1_p4",
							ReferenceClass: "string"}},
						CimAnyParameter{ParameterArray: &CimParameterArray{
							Name:      "method_1_p5",
							Type:      "string",
							ArraySize: 5,
							Qualifiers: []CimQualifier{
								makeQualifier("m1_q_p5_1"),
								makeQualifier("m1_q_p5_2"),
							}}},
						CimAnyParameter{ParameterArray: &CimParameterArray{
							Name:      "method_1_p6",
							Type:      "string",
							ArraySize: 5}},

						CimAnyParameter{ParameterRefArray: &CimParameterRefArray{
							Name:           "method_1_p6",
							ReferenceClass: "string",
							ArraySize:      6,
							Qualifiers: []CimQualifier{
								makeQualifier("m1_q_p6_1"),
								makeQualifier("m1_q_p6_2"),
							}}},

						CimAnyParameter{ParameterRefArray: &CimParameterRefArray{
							Name:           "method_1_p6",
							ReferenceClass: "string",
							ArraySize:      6}},
					},
				},
			},
		},
	},

	CimParamValue{
		Name:           "p7",
		ParamType:      "string",
		EmbeddedObject: "instance",
		Instance: &CimInstance{
			ClassName:  "a.b.c.class_test_p7",
			Properties: []CimAnyProperty{},
		},
	},
}

var simple_req1 = CimSimpleReq{
	Correlators: []CimCorrelator{CimCorrelator{Name: "cor1", Type: "string", Value: CimValue{Value: "cor1Value"}},
		CimCorrelator{Name: "cor2", Type: "string", Value: CimValue{Value: "cor2Value"}}},
	MethodCall: &CimMethodCall{
		Name:           "abc",
		LocalClassPath: makeLocalClass([]string{"a", "b"}, "class1"),
		ParamValues:    paramValues},
}

var simple_req2 = CimSimpleReq{
	Correlators: []CimCorrelator{CimCorrelator{Name: "cor1", Type: "string", Value: CimValue{Value: "cor1Value"}},
		CimCorrelator{Name: "cor2", Type: "string", Value: CimValue{Value: "cor2Value"}}},
	MethodCall: &CimMethodCall{
		Name: "abc",
		LocalInstancePath: &CimLocalInstancePath{
			LocalNamespacePath: *makeLocalNamespace([]string{"a", "b"}),
			InstanceName: CimInstanceName{
				ClassName: "cls_test23"},
		},
	},
}

var simple_req3 = CimSimpleReq{
	IMethodCall: &CimIMethodCall{},
}

var req = &CIM{CimVersion: "1.2.3.4",
	DtdVersion: "4.5.6.7",
	Message: &CimMessage{Id: "12",
		ProtocolVersion: "1.2.3.7",
		MultiReq: &CimMultiReq{SimpleReqs: []CimSimpleReq{simple_req1,
			simple_req2,
			simple_req3},
		},
	}}

func TestMultiReq(t *testing.T) {
	bs, e := xml.MarshalIndent(req, "", "  ")
	if nil != e {
		t.Error(e)
		return
	} else {
		t.Log(string(bs))
	}

	var req2 CIM
	if e := xml.Unmarshal(bs, &req2); nil != e {
		t.Error(e)
		return
	}

	if !reflect.DeepEqual(req, req2) {

		bs2, e := xml.MarshalIndent(req2, "", "  ")
		if nil != e {
			t.Error(e)
			return
		} else {

			if string(bs) != string(bs2) {
				t.Errorf("excepted is %#v", req)
				t.Errorf("actual is %#v", req2)
				//t.Log(string(bs))

				results := difflib.Diff(strings.Split(string(bs), "\n"), strings.Split(string(bs2), "\n"))
				if 0 != len(results) {
					for _, rec := range results {
						t.Error(rec.String())
					}
				}
			}

		}
	}
}

var simple_rsp1 = CimSimpleRsp{
	IMethodResponse: &CimIMethodResponse{
		Name:        "abc",
		ParamValues: paramValues,
		ReturnValue: &CimIReturnValue{
			ClassNames: []CimClassName{CimClassName{Name: "abc"}},
		},
	},
}

var simple_error_rsp = CimSimpleRsp{
	IMethodResponse: &CimIMethodResponse{
		Name: "err_rsp",
		//ParamValues: paramValues,
		Error: &CimError{Code: 123, Description: "test message"},
	},
}

func TestSimpleRsp(t *testing.T) {
	for _, rsp := range []CimSimpleRsp{simple_rsp1, simple_error_rsp} {
		bs, e := xml.MarshalIndent(rsp, "", "  ")
		if nil != e {
			t.Error(e)
			return
		} else {
			t.Log(string(bs))
		}

		var unmarshal_rsp CimSimpleRsp
		if e := xml.Unmarshal(bs, &unmarshal_rsp); nil != e {
			t.Error(e)
			return
		}

		if !reflect.DeepEqual(rsp, unmarshal_rsp) {

			bs2, e := xml.MarshalIndent(unmarshal_rsp, "", "  ")
			if nil != e {
				t.Error(e)
				return
			} else {

				if string(bs) != string(bs2) {
					t.Errorf("excepted is %#v", rsp)
					t.Errorf("actual is %#v", unmarshal_rsp)
					//t.Log(string(bs))

					results := difflib.Diff(strings.Split(string(bs), "\n"), strings.Split(string(bs2), "\n"))
					if 0 != len(results) {
						for _, rec := range results {
							t.Error(rec.String())
						}
					}
				}
			}
		}
	}
}

func TestCimError(t *testing.T) {
	err_txt := `<CIM CIMVERSION="2.0" DTDVERSION="2.0">
<MESSAGE ID="1-563c39eab1c2802414000002" PROTOCOLVERSION="1.0">
<SIMPLERSP>
<IMETHODRESPONSE NAME="Associators">
<ERROR CODE="100" DESCRIPTION="Unrecognized CIM status code &quot;100&quot;: Cannot connect to local CIM server. Connection failed."/></IMETHODRESPONSE>
</SIMPLERSP>
</MESSAGE>
</CIM>`

	var unmarshal_rsp = CIM{hasFault: func(cim *CIM) error {
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

	if e := xml.Unmarshal([]byte(err_txt), &unmarshal_rsp); nil != e {
		t.Error(e)
		return
	}

	//fmt.Println(unmarshal_rsp.Fault())
	//if unmarshal_rsp.Fault() {
	//}

	if 100 != unmarshal_rsp.Message.SimpleRsp.IMethodResponse.Error.Code {
		t.Error("code is error")
		return
	}

	t.Log(unmarshal_rsp.Message.SimpleRsp.IMethodResponse.Error.Description)
}

var class = &CimClass{
	Name:       "a.b.c.class_test_p7",
	SuperClass: "a.b.c.class_test_p7_base",

	Properties: []CimAnyProperty{
		CimAnyProperty{Property: &CimProperty{
			Name:           "pr_p7_1_abc_1",
			Type:           "string",
			ClassOrigin:    "pr_p7_1_abc_1_origin",
			Propagated:     true,
			EmbeddedObject: "object",
			Lang:           "cn",
			Value:          &CimValue{Value: "vvvvv"},
		}},
		CimAnyProperty{PropertyArray: &CimPropertyArray{
			Name:           "pr_p7_1_abc_2",
			Type:           "string",
			ArraySize:      23,
			ClassOrigin:    "pr_p7_1_abc_2_origin",
			Propagated:     true,
			EmbeddedObject: "object",
			Lang:           "cn",
			ValueArray: makeValueArray(
				CimValueOrNull{Value: &CimValue{Value: "vvvvv"}},
			),
		}},
		CimAnyProperty{PropertyReference: &CimPropertyReference{
			Name:           "pr_p7_1_abc_3",
			ReferenceClass: "ref_class",
			ClassOrigin:    "pr_p7_1_abc_3_origin",

			Propagated: true,
			ValueReference: &CimValueReference{
				LocalInstancePath: &CimLocalInstancePath{
					LocalNamespacePath: *makeLocalNamespace([]string{"f", "t"}),
					InstanceName:       *makeInstanceNameWithKV2("F.4.A.class_abc", "ref_class", "abcss", "bool"),
				},
			},
		}},
	},
}

func TestCimClass(t *testing.T) {
	bs, e := xml.MarshalIndent(class, "", "  ")
	if nil != e {
		t.Error(e)
		return
	} else {
		t.Log(string(bs))
	}

	var cls2 CimClass
	if e := xml.Unmarshal(bs, &cls2); nil != e {
		t.Error(e)
		return
	}

	if !reflect.DeepEqual(class, cls2) {

		bs2, e := xml.MarshalIndent(cls2, "", "  ")
		if nil != e {

			t.Errorf("excepted is %#v", class)
			t.Errorf("actual is %#v", cls2)
			t.Error(e)
			return
		} else {

			if string(bs) != string(bs2) {

				t.Errorf("excepted is %#v", class)
				t.Errorf("actual is %#v", cls2)
				//t.Log(string(bs))

				results := difflib.Diff(strings.Split(string(bs), "\n"), strings.Split(string(bs2), "\n"))
				if 0 != len(results) {
					for _, rec := range results {
						t.Error(rec.String())
					}
				}
			}
		}
	}
}

var method = CimMethod{
	Name:        "method_1",
	Type:        "string",
	ClassOrigin: "method_1_origin",
	Propagated:  true,
	Qualifiers: []CimQualifier{
		makeQualifier("m1_q_1"),
		makeQualifier("m1_q_2"),
	},
	Parameters: []CimAnyParameter{
		CimAnyParameter{Parameter: &CimParameter{
			Name: "method_1_p1",
			Type: "string",
			Qualifiers: []CimQualifier{
				makeQualifier("m1_q_p1_1"),
				makeQualifier("m1_q_p1_2"),
			}}},
		CimAnyParameter{Parameter: &CimParameter{
			Name: "method_1_p2",
			Type: "string"}},
		CimAnyParameter{ParameterReference: &CimParameterReference{
			Name:           "method_1_p3",
			ReferenceClass: "string",
			Qualifiers: []CimQualifier{
				makeQualifier("m1_q_p3_1"),
				makeQualifier("m1_q_p3_2"),
			}}},
		CimAnyParameter{ParameterReference: &CimParameterReference{
			Name:           "method_1_p4",
			ReferenceClass: "string"}},
		CimAnyParameter{ParameterArray: &CimParameterArray{
			Name:      "method_1_p5",
			Type:      "string",
			ArraySize: 5,
			Qualifiers: []CimQualifier{
				makeQualifier("m1_q_p5_1"),
				makeQualifier("m1_q_p5_2"),
			}}},
		CimAnyParameter{ParameterArray: &CimParameterArray{
			Name:      "method_1_p6",
			Type:      "string",
			ArraySize: 5}},

		CimAnyParameter{ParameterRefArray: &CimParameterRefArray{
			Name:           "method_1_p6",
			ReferenceClass: "string",
			ArraySize:      6,
			Qualifiers: []CimQualifier{
				makeQualifier("m1_q_p6_1"),
				makeQualifier("m1_q_p6_2"),
			}}},

		CimAnyParameter{ParameterRefArray: &CimParameterRefArray{
			Name:           "method_1_p6",
			ReferenceClass: "string",
			ArraySize:      6}},
	},
}

func TestCimMethod(t *testing.T) {
	bs, e := xml.MarshalIndent(method, "", "  ")
	if nil != e {
		t.Error(e)
		return
	} else {
		t.Log(string(bs))
	}

	var method2 CimMethod
	if e := xml.Unmarshal(bs, &method2); nil != e {
		t.Error(e)
		return
	}

	if !reflect.DeepEqual(method, method2) {
		bs2, e := xml.MarshalIndent(method2, "", "  ")
		if nil != e {

			t.Errorf("excepted is %#v", method)
			t.Errorf("actual is %#v", method2)
			t.Error(e)
			return
		} else {

			if string(bs) != string(bs2) {
				t.Errorf("excepted is %#v", method)
				t.Errorf("actual is %#v", method2)
				//t.Log(string(bs))

				results := difflib.Diff(strings.Split(string(bs), "\n"), strings.Split(string(bs2), "\n"))
				if 0 != len(results) {
					for _, rec := range results {
						t.Error(rec.String())
					}
				}
			}
		}
	}
}

var declaration = CimDeclaration{
	DeclGroups: []CimAnyDeclGroup{
		CimAnyDeclGroup{
			DeclGroup: &CimDeclGroup{
				LocalNamespacePath: makeLocalNamespace([]string{"A", "b"}),
				QualifierDeclarations: []CimQualifierDeclaration{
					CimQualifierDeclaration{Name: "abc"},
					CimQualifierDeclaration{Name: "abc1",
						Type:      "string",
						IsArray:   true,
						ArraySize: 12,
						Scope: &CimScope{Class: true,
							Association: true,
							Reference:   true,
							Property:    true,
							Method:      true,
							Parameter:   true,
							Indication:  true},
						Value: &CimValue{Value: "abcvvvv"},
					},
				},
			},
		},
		CimAnyDeclGroup{
			DeclGroup: &CimDeclGroup{
				NamespacePath: makeNamespace("192.168.1.2", []string{"A", "b"}),
				QualifierDeclarations: []CimQualifierDeclaration{
					CimQualifierDeclaration{Name: "abc"},
					CimQualifierDeclaration{Name: "abc1",
						Type:      "string",
						IsArray:   true,
						ArraySize: 12,
						Scope: &CimScope{Class: true,
							Association: true,
							Reference:   true,
							Property:    true,
							Method:      true,
							Parameter:   true,
							Indication:  true},
						ValueArray: &CimValueArray{
							CimValueOrNull{Value: &CimValue{Value: "abcvvvv1"}},
							CimValueOrNull{Null: &CimValueNull{}},
							CimValueOrNull{Value: &CimValue{Value: "abcvvvv3"}},
						},
					},
				},
				ValueObjects: []CimValueObject{
					CimValueObject{Instance: &CimInstance{
						ClassName: "a_class",
						Lang:      "zh",
					}},
				},
			},
		},
		CimAnyDeclGroup{
			DeclGroupWithName: &CimDeclGroupWithName{
				LocalNamespacePath: makeLocalNamespace([]string{"A", "b"}),
				QualifierDeclarations: []CimQualifierDeclaration{
					CimQualifierDeclaration{Name: "abc"},
					CimQualifierDeclaration{Name: "abc1",
						Type:      "string",
						IsArray:   true,
						ArraySize: 12,
						Scope: &CimScope{Class: true,
							Association: true,
							Reference:   true,
							Property:    true,
							Method:      true,
							Parameter:   true,
							Indication:  true},
						Value: &CimValue{Value: "abcvvvv"},
					},
				},
				ValueNamedObjects: []CimValueNamedObject{
					CimValueNamedObject{InstanceName: makeInstanceNameWithValueRef("abc", makeValueRef("delar_1")),
						Instance: makeInstance("declar_vn_2")},
					CimValueNamedObject{Class: class},
					CimValueNamedObject{},
				},
			},
		},
		CimAnyDeclGroup{
			DeclGroupWithPath: &CimDeclGroupWithPath{},
		},
	},
}

func TestCimDeclaration(t *testing.T) {
	bs, e := xml.MarshalIndent(declaration, "", "  ")
	if nil != e {
		t.Error(e)
		return
	} else {
		t.Log(string(bs))
	}

	var declaration2 CimDeclaration
	if e := xml.Unmarshal(bs, &declaration2); nil != e {
		t.Error(e)
		return
	}

	if !reflect.DeepEqual(declaration, declaration2) {

		// t.Errorf("excepted is %#v", declaration)
		// t.Errorf("actual is %#v", declaration2)

		bs2, e := xml.MarshalIndent(declaration2, "", "  ")
		if nil != e {

			t.Errorf("excepted is %#v", declaration)
			t.Errorf("actual is %#v", declaration2)
			t.Error(e)
			return
		} else {

			if string(bs) != string(bs2) {

				//t.Log(string(bs))

				results := difflib.Diff(strings.Split(string(bs), "\n"), strings.Split(string(bs2), "\n"))
				if 0 != len(results) {
					for _, rec := range results {
						t.Error(rec.String())
					}
				}
			}
		}
	}
}
