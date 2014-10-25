package wbem

type CIMTypeCode int

const (
	INVALID CIMTypeCode = iota
	BOOLEAN
	STRING
	CHAR16
	UINT8
	SINT8
	UINT16
	SINT16
	UINT32
	SINT32
	UINT64
	SINT64
	DATETIME
	REAL32
	REAL64
	NUMERIC
	REFERENCE
)

const (
	NON_ARRAY       = 0
	UNBOUNDED_ARRAY = -1
)

var (
	INVALID_TYPE = CIMType{typeCode: INVALID}

	types = map[string]CIMType{
		"":          CIMType{typeCode: INVALID},
		"boolean":   CIMType{typeCode: BOOLEAN},
		"string":    CIMType{typeCode: STRING},
		"char16":    CIMType{typeCode: CHAR16},
		"uint8":     CIMType{typeCode: UINT8},
		"sint8":     CIMType{typeCode: SINT8},
		"uint16":    CIMType{typeCode: UINT16},
		"sint16":    CIMType{typeCode: SINT16},
		"uint32":    CIMType{typeCode: UINT32},
		"sint32":    CIMType{typeCode: SINT32},
		"uint64":    CIMType{typeCode: UINT64},
		"sint64":    CIMType{typeCode: SINT64},
		"datetime":  CIMType{typeCode: DATETIME},
		"real32":    CIMType{typeCode: REAL32},
		"real64":    CIMType{typeCode: REAL64},
		"numeric":   CIMType{typeCode: NUMERIC},
		"reference": CIMType{typeCode: REFERENCE},
	}
)

type CIMType struct {
	typeCode CIMTypeCode
	/**
	 * non array if =0  <br>
	 * unbounded if <0  <br>
	 * bounded if >0
	 */
	arraySize    int
	refClassName string
}

/**
 * Returns the class name of the CIM REFERENCE data type.
 *
 * @return The CIM REFERENCE class name.
 */
func (self *CIMType) GetClassName() string {
	return self.refClassName
}

/**
 * Returns the size of the maximum number of elements an array data type may
 * hold.
 *
 * @return The maximum size of the array data type.
 */
func (self *CIMType) GetSize() int {
	return self.arraySize
}

/**
 * Returns the data type.
 *
 * @return The data type.
 */
func (self *CIMType) GetType() CIMTypeCode {
	return self.typeCode
}

/**
 * Checks if the data type is an array type.
 *
 * @return <code>true</code> if the data type is an array type,
 *         <code>false</code> otherwise.
 */
func (self *CIMType) IsArray() bool {
	return self.arraySize == 0
}

func CreateCIMType(t string) CIMType {
	if ct, ok := types[t]; ok {
		return ct
	}
	return CIMType{typeCode: INVALID}
}

func CreateCIMArrayType(t string, arraySize int) CIMType {
	if ct, ok := types[t]; ok {
		ct.arraySize = arraySize
		return ct
	}
	return CIMType{typeCode: INVALID, arraySize: arraySize}
}

func CreateCIMReferenceType(class string) CIMType {
	return CIMType{typeCode: REFERENCE, refClassName: class}
}
