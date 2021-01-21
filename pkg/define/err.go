package define

import (
	"errors"
	"fmt"
)

const (
	InnerServerErr = "服务器内部错误"
)

type TpaasError interface {
	// error info to response.Detail
	Error() string
	// return code to response.Code
	RetCode() int
	// response message to response.Msg
	Message() string
}

///////////

var _ TpaasError = new(errNotFound)
var _ TpaasError = new(errSvcNotFound)
var _ TpaasError = new(errToken)
var _ TpaasError = new(errInnerServer)
var _ TpaasError = new(errSvcTmout)

var _ error = new(errNotFound)
var _ error = new(errSvcNotFound)
var _ error = new(errToken)
var _ error = new(errInnerServer)
var _ error = new(errSvcTmout)

type errSvcTmout struct {
	Key  string
	Code int
}

func ErrSvcTmout(servicename string) error {
	return &errSvcTmout{Key: servicename, Code: StatusGatewayTimeout}
}

func (e *errSvcTmout) Error() string {
	return fmt.Sprintf("service[%s] is timeout", e.Key)
}

func (e *errSvcTmout) RetCode() int {
	return StatusGatewayTimeout
}

func (e *errSvcTmout) Message() string {
	return fmt.Sprintf("%s超时", e.Key)
}

type errNotFound struct {
	Key   string
	Value string
	Code  int
}

func ErrNotFound(res, name string) error {
	return &errNotFound{Key: res, Value: name, Code: StatusNotFound}
}

func (e *errNotFound) Error() string {
	return fmt.Sprintf("%s[%s] not exists in es", e.Key, e.Value)
}

func (e *errNotFound) RetCode() int {
	return StatusNotFound
}

func (e *errNotFound) Message() string {
	return fmt.Sprintf("%s不存在", e.Key)
}

type errSvcNotFound struct {
	Key   string
	Value string
	Code  int
}

func ErrSvcNotFound(svc, name string) error {
	return &errSvcNotFound{Key: svc, Value: name, Code: StatusServiceUnavailable}
}

func (e *errSvcNotFound) Message() string {
	return "服务未开通"
}

func (e *errSvcNotFound) Error() string {
	return fmt.Sprintf("%s[%s] not exists in es", e.Key, e.Value)
}

func (e *errSvcNotFound) RetCode() int {
	return StatusServiceUnavailable
}

////////
type errNoPermission struct {
	Resource string
	Action   string
	Who      string
	Code     int
}

/////
type errToken struct {
	Code int
}

func ErrToken() error {
	return &errToken{Code: StatusUnauthorized}
}

func (e *errToken) Error() string {
	return fmt.Sprintf("token error")
}

func (e *errToken) RetCode() int {
	return StatusUnauthorized
}
func (e *errToken) Message() string {
	return fmt.Sprintf("token失效")
}

//////
type errInnerServer struct {
	Code int
	Func string
	Err  error
}

func ErrInnerServer(f string, err error) error {
	if err == nil {
		err = errors.New("内部错误")
	}
	return &errInnerServer{Code: StatusInternalServerError, Func: f, Err: err}
}

func (e *errInnerServer) Error() string {
	return fmt.Sprintf("%s err:%s", e.Func, e.Err.Error())
}

func (e *errInnerServer) RetCode() int {
	return StatusInternalServerError
}

func (e *errInnerServer) Message() string {
	return fmt.Sprintf("服务器内部错误")
}

/////
func ErrorMsg(e error) (int, string, string) {
	te, ok := e.(TpaasError)
	if !ok {
		return StatusInternalServerError, InnerServerErr, e.Error()
	}
	return te.RetCode(), te.Message(), te.Error()
}

// 数据库错误
type errDB struct {
	Table string
	Func  string
	Err   error
}

func ErrDB(err error, table, funcname string) error {
	return &errDB{Table: table, Func: funcname, Err: err}
}

func (e *errDB) Error() string {
	return fmt.Sprintf("Table:%s, Func:%s, Err:%s", e.Table, e.Func, e.Err)
}
