/*
Copyright (c) 2015 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// this code is copy from https://github.com/vmware/govmomi/blob/master/vim25/soap/client.go
package gowbem

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync/atomic"
)

type HasFault interface {
	Fault() error
}

type DecodeError struct {
	bytes []byte
	err   error
}

func (f *DecodeError) Error() string {
	return f.err.Error() + ", xml as follow:\r\n" + string(f.bytes)
}

type FaultError struct {
	bytes []byte
	err   error
}

func (f *FaultError) Error() string {
	return f.err.Error() + ", xml as follow:\r\n" + string(f.bytes)
}

type RoundTripper interface {
	RoundTrip(action string, reqBody interface{}, resBody HasFault, cached *bytes.Buffer) error
}

var cn uint64 // Client counter

type Client struct {
	rn uint64 // Request counter

	http.Client

	u        url.URL
	insecure bool

	cn_str string // Client counter
	cn     uint64 // Client counter
	cached *bytes.Buffer
}

func NewClient(u *url.URL, insecure bool) *Client {
	c := &Client{}
	c.init(u, insecure)
	return c
}

func (c *Client) init(u *url.URL, insecure bool) {
	c.u = *u
	c.insecure = insecure
	c.cn = atomic.AddUint64(&cn, 1)
	c.rn = 0
	c.cached = bytes.NewBuffer(make([]byte, 0, 8*1024*1024))

	c.cn_str = strconv.FormatUint(c.cn, 10)

	if c.u.Scheme == "https" {
		c.Client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: c.insecure}}
	}

	c.Jar, _ = cookiejar.New(nil)
	//c.u.User = nil
}

func (c *Client) generateId() string {
	return c.cn_str + "-" + GenerateId()
}

func (c *Client) URL() url.URL {
	return c.u
}

type marshaledClient struct {
	Cookies  []*http.Cookie
	URL      *url.URL
	Insecure bool
}

func (c *Client) MarshalJSON() ([]byte, error) {
	m := marshaledClient{
		Cookies:  c.Jar.Cookies(&c.u),
		URL:      &c.u,
		Insecure: c.insecure,
	}

	return json.Marshal(m)
}

func (c *Client) UnmarshalJSON(b []byte) error {
	var m marshaledClient

	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	*c = *NewClient(m.URL, m.Insecure)
	c.Jar.SetCookies(m.URL, m.Cookies)
	return nil
}

func (c *Client) RoundTrip(ctx context.Context, action string, headers map[string]string, reqBody interface{}, resBody HasFault) error {
	var httpreq *http.Request
	var httpres *http.Response
	var dumpWriter io.WriteCloser
	var err error

	num := atomic.AddUint64(&c.rn, 1)
	c.cached.Reset()

	c.cached.WriteString(xml.Header)

	if err = xml.NewEncoder(c.cached).Encode(reqBody); err != nil {
		panic(err)
	}
	rawreqbody := c.cached

	httpreq, err = http.NewRequest(action, c.u.String(), rawreqbody)
	if err != nil {
		panic(err)
	}
	if ctx != nil {
		httpreq = httpreq.WithContext(ctx)
	}

	//if pwd, ok := c.u.User.Password(); ok {
	//	httpreq.SetBasicAuth(c.u.User.Username(), pwd)
	//}

	httpreq.Header.Set(`Content-Type`, `text/xml; charset="utf-8"`)
	if 0 != len(headers) {
		for k, v := range headers {
			httpreq.Header.Set(k, v)
		}
	}

	if DebugEnabled() {
		b, _ := httputil.DumpRequest(httpreq, false)
		dumpWriter = DebugNewFile(fmt.Sprintf("%d-%04d.log", c.cn, num))
		defer dumpWriter.Close()
		dumpWriter.Write(b)
		dumpWriter.Write(c.cached.Bytes())
	}

	//tstart := time.Now()
	httpres, err = c.Client.Do(httpreq)
	//tstop := time.Now()

	// if DebugEnabled() {
	// 	now := time.Now().Format("2006-01-02T15-04-05.999999999")
	// 	ms := tstop.Sub(tstart) / time.Millisecond
	// 	fmt.Fprintf(c.log, "%s: %4d took %6dms\n", now, num, ms)
	// }

	if err != nil {

		if DebugEnabled() {
			dumpWriter.Write([]byte("\r\n"))
			dumpWriter.Write([]byte(err.Error()))
		}

		return err
	}

	//var rawresbody io.Reader = httpres.Body
	defer httpres.Body.Close()

	if httpres.ContentLength <= 0 && (httpres.StatusCode < http.StatusOK || httpres.StatusCode >= http.StatusMultipleChoices) {

		if DebugEnabled() {
			b, _ := httputil.DumpResponse(httpres, false)
			dumpWriter.Write([]byte("\r\n"))
			dumpWriter.Write(b)
			dumpWriter.Write(c.cached.Bytes())
		}

		// 修复 pg 导到一个问题， 当pg出错时返回错误响应时，没有 ContentLength， tcp 连接也不关闭。
		// 这时读 httpres.Body 时会导致本方法挂起。
		// This is Pegasus bug, pegasus will return a error response that http version is 1.0 and ContentLength is missing.
		// And tcp connection isn't disconnect by the pegasus server.
		cimError := httpres.Header.Get("CIMError")
		errorDetail := httpres.Header.Get("PGErrorDetail")
		if "" == cimError {
			if "" != errorDetail {
				return errors.New(errorDetail)
			}
			return errors.New(httpres.Status)
		} else {
			if "" != errorDetail {
				return errors.New(cimError + ": " + errorDetail)
			}
			return errors.New(cimError)
		}
	}

	c.cached.Reset()
	if _, err = io.Copy(c.cached, httpres.Body); nil != err {
		return err
	}

	if DebugEnabled() {
		b, _ := httputil.DumpResponse(httpres, false)
		dumpWriter.Write([]byte("\r\n"))
		dumpWriter.Write(b)
		dumpWriter.Write(c.cached.Bytes())
	}

	if 200 != httpres.StatusCode {
		if 0 == c.cached.Len() {
			return errors.New(httpres.Status)
		}
		return errors.New(httpres.Status + ":" + c.cached.String())
	}

	dec := xml.NewDecoder(bytes.NewReader(c.cached.Bytes()))
	err = dec.Decode(resBody)
	if err != nil {
		return &DecodeError{bytes: c.cached.Bytes(), err: err}
	}

	if fault := resBody.Fault(); fault != nil {
		return &FaultError{bytes: c.cached.Bytes(), err: fault}
	}

	return nil
}

func StringsWith(instance CIMInstance, key string, defaultVlaue []string) []string {
	prop := instance.GetPropertyByName(key)
	if prop == nil {
		return defaultVlaue
	}
	value := prop.GetValue()
	if value == nil {
		return defaultVlaue
	}
	switch vv := value.(type) {
	case []string:
		return vv
	case []interface{}:
		var ss = make([]string, 0, len(vv))
		for _, v := range vv {
			ss = append(ss, fmt.Sprint(v))
		}
		return ss
	default:
		panic(fmt.Sprintf("[wbem] value isn't a array - %T %#v.", value, value))
	}
}
