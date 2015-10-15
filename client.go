package gowbem

import (
	"bytes"
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
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/vmware/govmomi/vim25/progress"
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

func (c *Client) RoundTrip(action string, headers map[string]string, reqBody interface{}, resBody HasFault) error {
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
	//if pwd, ok := c.u.User.Password(); ok {
	//	httpreq.SetBasicAuth(c.u.User.Username(), pwd)
	//}

	httpreq.Header.Set(`Content-Type`, `text/xml; charset="utf-8"`)
	if 0 != len(headers) {
		for k, v := range headers {
			httpreq.Header.Set(k, v)
		}
	}
	//httpreq.Header.Set(`SOAPAction`, `urn:vim25/5.5`)

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
		return err
	}

	//var rawresbody io.Reader = httpres.Body
	defer httpres.Body.Close()

	if httpres.ContentLength <= 0 && httpres.StatusCode < http.StatusOK && httpres.StatusCode >= http.StatusMultipleChoices {
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

// ParseURL wraps url.Parse to rewrite the URL.Host field
// In the case of VM guest uploads or NFC lease URLs, a Host
// field with a value of "*" is rewritten to the Client's URL.Host.
func (c *Client) ParseURL(urlStr string) (*url.URL, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	host := strings.Split(u.Host, ":")
	if host[0] == "*" {
		// Also use Client's port, to support port forwarding
		u.Host = c.URL().Host
	}

	return u, nil
}

type Upload struct {
	Type          string
	Method        string
	ContentLength int64
	Progress      progress.Sinker
}

var DefaultUpload = Upload{
	Type:   "application/octet-stream",
	Method: "PUT",
}

// Upload PUTs the local file to the given URL
func (c *Client) Upload(f io.Reader, u *url.URL, param *Upload) error {
	var err error

	if param.Progress != nil {
		pr := progress.NewReader(param.Progress, f, param.ContentLength)
		f = pr

		// Mark progress reader as done when returning from this function.
		defer func() {
			pr.Done(err)
		}()
	}

	req, err := http.NewRequest(param.Method, u.String(), f)
	if err != nil {
		return err
	}

	req.ContentLength = param.ContentLength
	req.Header.Set("Content-Type", param.Type)

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusCreated:
	default:
		err = errors.New(res.Status)
	}

	return err
}

// UploadFile PUTs the local file to the given URL
func (c *Client) UploadFile(file string, u *url.URL, param *Upload) error {
	if param == nil {
		p := DefaultUpload // Copy since we set ContentLength
		param = &p
	}

	s, err := os.Stat(file)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	param.ContentLength = s.Size()

	return c.Upload(f, u, param)
}

type Download struct {
	Method   string
	Progress progress.Sinker
}

var DefaultDownload = Download{
	Method: "GET",
}

// DownloadFile GETs the given URL to a local file
func (c *Client) DownloadFile(file string, u *url.URL, param *Download) error {
	var err error

	if param == nil {
		param = &DefaultDownload
	}

	fh, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fh.Close()

	req, err := http.NewRequest(param.Method, u.String(), nil)
	if err != nil {
		return err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
	default:
		err = errors.New(res.Status)
	}

	if err != nil {
		return err
	}

	var r io.Reader = res.Body
	if param.Progress != nil {
		pr := progress.NewReader(param.Progress, res.Body, res.ContentLength)
		r = pr

		// Mark progress reader as done when returning from this function.
		defer func() {
			pr.Done(err)
		}()
	}

	_, err = io.Copy(fh, r)
	if err != nil {
		return err
	}

	// Assign error before returning so that it gets picked up by the deferred
	// function marking the progress reader as done.
	err = fh.Close()
	if err != nil {
		return err
	}

	return nil
}
