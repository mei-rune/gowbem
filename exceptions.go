package gowbem

import (
	"strconv"
	"strings"
)

type CIMStatusCode int

const (
	// CIMError error code constants
	CIM_ERR_FAILED                              CIMStatusCode = 1  // A general error occurred
	CIM_ERR_ACCESS_DENIED                       CIMStatusCode = 2  // Resource not available
	CIM_ERR_INVALID_NAMESPACE                   CIMStatusCode = 3  // The target namespace does not exist
	CIM_ERR_INVALID_PARAMETER                   CIMStatusCode = 4  // Parameter value(s) invalid
	CIM_ERR_INVALID_CLASS                       CIMStatusCode = 5  // The specified Class does not exist
	CIM_ERR_NOT_FOUND                           CIMStatusCode = 6  // Requested object could not be found
	CIM_ERR_NOT_SUPPORTED                       CIMStatusCode = 7  // Operation not supported
	CIM_ERR_CLASS_HAS_CHILDREN                  CIMStatusCode = 8  // Class has subclasses
	CIM_ERR_CLASS_HAS_INSTANCES                 CIMStatusCode = 9  // Class has instances
	CIM_ERR_INVALID_SUPERCLASS                  CIMStatusCode = 10 // Superclass does not exist
	CIM_ERR_ALREADY_EXISTS                      CIMStatusCode = 11 // Object already exists
	CIM_ERR_NO_SUCH_PROPERTY                    CIMStatusCode = 12 // Property does not exist
	CIM_ERR_TYPE_MISMATCH                       CIMStatusCode = 13 // Value incompatible with type
	CIM_ERR_QUERY_LANGUAGE_NOT_SUPPORTED        CIMStatusCode = 14 // Query language not supported
	CIM_ERR_INVALID_QUERY                       CIMStatusCode = 15 // Query not valid
	CIM_ERR_METHOD_NOT_AVAILABLE                CIMStatusCode = 16 // Extrinsic method not executed
	CIM_ERR_METHOD_NOT_FOUND                    CIMStatusCode = 17 // Extrinsic method does not exist
	CIM_ERR_NAMESPACE_NOT_EMPTY                 CIMStatusCode = 20
	CIM_ERR_INVALID_ENUMERATION_CONTEXT         CIMStatusCode = 21
	CIM_ERR_INVALID_OPERATION_TIMEOUT           CIMStatusCode = 22
	CIM_ERR_PULL_HAS_BEEN_ABANDONED             CIMStatusCode = 23
	CIM_ERR_PULL_CANNOT_BE_ABANDONED            CIMStatusCode = 24
	CIM_ERR_FILTERED_ENUMERATION_NOT_SUPPORTED  CIMStatusCode = 25
	CIM_ERR_CONTINUATION_ON_ERROR_NOT_SUPPORTED CIMStatusCode = 26
	CIM_ERR_SERVER_LIMITS_EXCEEDED              CIMStatusCode = 27
	CIM_ERR_SERVER_IS_SHUTTING_DOWN             CIMStatusCode = 28
)

var (
	CIMStatusCodes = []string{
		"COM_ERR_OK",
		"CIM_ERR_FAILED",
		"CIM_ERR_ACCESS_DENIED",
		"CIM_ERR_INVALID_NAMESPACE",
		"CIM_ERR_INVALID_PARAMETER",
		"CIM_ERR_INVALID_CLASS",
		"CIM_ERR_NOT_FOUND",
		"CIM_ERR_NOT_SUPPORTED",
		"CIM_ERR_CLASS_HAS_CHILDREN",
		"CIM_ERR_CLASS_HAS_INSTANCES",
		"CIM_ERR_INVALID_SUPERCLASS",
		"CIM_ERR_ALREADY_EXISTS",
		"CIM_ERR_NO_SUCH_PROPERTY",
		"CIM_ERR_TYPE_MISMATCH",
		"CIM_ERR_QUERY_LANGUAGE_NOT_SUPPORTED",
		"CIM_ERR_INVALID_QUERY",
		"CIM_ERR_METHOD_NOT_AVAILABLE",
		"CIM_ERR_METHOD_NOT_FOUND",
		"CIM_ERR_18",
		"CIM_ERR_19",
		"CIM_ERR_NAMESPACE_NOT_EMPTY",
		"CIM_ERR_INVALID_ENUMERATION_CONTEXT",
		"CIM_ERR_INVALID_OPERATION_TIMEOUT",
		"CIM_ERR_PULL_HAS_BEEN_ABANDONED",
		"CIM_ERR_PULL_CANNOT_BE_ABANDONED",
		"CIM_ERR_FILTERED_ENUMERATION_NOT_SUPPORTED",
		"CIM_ERR_CONTINUATION_ON_ERROR_NOT_SUPPORTED",
		"CIM_ERR_SERVER_LIMITS_EXCEEDED",
		"CIM_ERR_SERVER_IS_SHUTTING_DOWN",
	}
)

func (code CIMStatusCode) String() string {
	if code > 0 && code <= 28 {
		return CIMStatusCodes[code]
	}
	return "COM_ERR_" + strconv.FormatInt(int64(code), 10)
}

type WbemError struct {
	code CIMStatusCode
	msg  string
}

func (e *WbemError) Error() string {
	code_str := e.code.String()
	if strings.Contains(e.msg, code_str) {
		return e.msg
	}
	return code_str + ":" + e.msg
}

func WBEMException(code CIMStatusCode, msg string) error {
	return &WbemError{code: code, msg: msg}
}

func IsErrNotSupported(e error) bool {
	if we, ok := e.(*WbemError); ok {
		return we.code == CIM_ERR_NOT_SUPPORTED
	}
	if fe, ok := e.(*FaultError); ok {
		return IsErrNotSupported(fe.err)
	}
	return false
}
