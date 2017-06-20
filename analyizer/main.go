package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/runner-mei/gowbem"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println(len(args))
		log.Fatalln("参数不正确！\r\n", args[0], " logs目录")
		return
	}
	output := args[1] + "_output"

	files, err := filepath.Glob(filepath.Join(args[1], "*.log"))
	if err != nil {
		log.Fatalln(err)
		return
	}

	for _, filename := range files {
		bs, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println(err)
			continue
		}

		bsArray := bytes.Split(bs, []byte("HTTP/1.1 200 OK"))
		if len(bsArray) != 2 {
			fmt.Println("====================")
			fmt.Println(string(bs))
			fmt.Println()
			continue
		}
		req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(bsArray[0])))
		if err != nil {
			fmt.Println("====================", err)
			fmt.Println(string(bs))
			fmt.Println()
			fmt.Println(err)
			continue
		}

		reqArray := bytes.Split(bsArray[0], []byte("\r\n\r\n"))
		if len(reqArray) < 2 {
			fmt.Println("====================")
			fmt.Println(string(bsArray[0]))
			fmt.Println()
			continue
		}

		var cimReq gowbem.CIM
		dec := xml.NewDecoder(bytes.NewReader(reqArray[1]))
		err = dec.Decode(&cimReq)
		if err != nil {
			fmt.Println("====================", err)
			fmt.Println(string(bsArray[0]))
			fmt.Println(string(reqArray[1]))
			fmt.Println()
			continue
		}

		respBytes := bs[len(bsArray[0]):]
		respBytes = bytes.Replace(respBytes, []byte("Transfer-Encoding: chunked\r\n"), []byte(""), -1)
		resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(respBytes)), req)
		if err != nil {
			fmt.Println("====================", err)
			fmt.Println(string(bs))
			fmt.Println()
			fmt.Println(err)
			continue
		}

		var cim gowbem.CIM

		dec = xml.NewDecoder(resp.Body)
		err = dec.Decode(&cim)
		if err != nil {
			fmt.Println("====================", err)
			fmt.Println(string(bs))
			fmt.Println()
			continue
		}

		if nil == cim.Message {
			continue
		}
		if nil == cim.Message.SimpleRsp {
			continue
		}
		if nil == cim.Message.SimpleRsp.IMethodResponse {
			continue
		}
		if nil != cim.Message.SimpleRsp.IMethodResponse.Error {
			continue
		}

		namespace := req.Header.Get("Cimobject")

		if "GetInstance" == cim.Message.SimpleRsp.IMethodResponse.Name {
			instances := cim.Message.SimpleRsp.IMethodResponse.ReturnValue.Instances

			for _, instance := range instances {
				classPath := filepath.Join(output, namespace, instance.GetClassName())
				if err := os.MkdirAll(classPath, 666); err != nil && !os.IsExist(err) {
					log.Fatalln(err)
				}
				for _, paramValue := range cimReq.Message.SimpleReq.IMethodCall.ParamValues {
					if paramValue.Name == "InstanceName" {
						filename := filepath.Join(classPath, "instances.txt")
						bs := []byte(paramValue.InstanceName.String())

						if old, err := ioutil.ReadFile(filename); err == nil {
							bs = append(bs, []byte("\r\n")...)
							bs = append(bs, old...)
						} else if !os.IsNotExist(err) {
							fmt.Println(err)
							continue
						}

						if err := ioutil.WriteFile(filename, bs, 666); err != nil {
							fmt.Println(err)
							continue
						}
					}
				}

				bs, err := xml.MarshalIndent(instance, "", "  ")
				if err != nil {
					fmt.Println(err)
					continue
				}

				if err := ioutil.WriteFile(filepath.Join(classPath, filepath.Base(filename)+".xml"), bs, 666); err != nil {
					fmt.Println(err)
					continue
				}
			}
		} else if "EnumerateInstanceNames" == cim.Message.SimpleRsp.IMethodResponse.Name {

			var className string
			for _, paramValue := range cimReq.Message.SimpleReq.IMethodCall.ParamValues {
				if paramValue.Name == "ClassName" {
					className = paramValue.ClassName.Name
					break
				}
			}
			if className == "" {
				fmt.Println("==================== class name is empty")
				fmt.Println(string(bs))
				fmt.Println()
				continue
			}

			if cim.Message.SimpleRsp.IMethodResponse.ReturnValue == nil {
				continue
			}

			if len(cim.Message.SimpleRsp.IMethodResponse.ReturnValue.InstanceNames) == 0 {
				continue
			}

			var buf bytes.Buffer
			for _, instanceName := range cim.Message.SimpleRsp.IMethodResponse.ReturnValue.InstanceNames {
				buf.WriteString(instanceName.String())
				buf.WriteString("\r\n")
			}

			/// @begin 将类定义写到文件
			classPath := filepath.Join(output, namespace, className)
			if err := os.MkdirAll(classPath, 666); err != nil && !os.IsExist(err) {
				log.Fatalln(err)
			}
			if err := ioutil.WriteFile(filepath.Join(classPath, "instances.txt"), buf.Bytes(), 666); err != nil {
				log.Fatalln(err)
			}
			/// @end

		}
	}
}
