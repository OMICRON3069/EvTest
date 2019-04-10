package httpServer

import (
	"EvTest/evBus"
	"EvTest/jankyError"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type request struct {
	proto, method string
	path, query   string
	head, body    string
	remoteAddr    string
}

//connectionElement represent the base element of connection pool
//intend to achieve tcp connection reuse
type connectionElement struct {
	connID int
	conn   net.Conn
}

type connectionPool struct {
	sync.RWMutex
	pool []connectionElement
}

type Config struct {
	BindAddr, BindPort string
}

func ServeHttp(bus *evBus.Bus, config Config) error {
	bus.Lock()
	defer bus.Unlock()
	if  {

	}
	ln, err := net.Listen("tcp", config.BindAddr+":"+config.BindPort)
	fmt.Println("Listening on: " + ln.Addr().String())
	if err != nil {
		return &jankyError.TheError{
			//TODO
		}
	}
	go func() {
		for {
			conn, err := ln.Accept()
			fmt.Println("incoming " + conn.RemoteAddr().String() + " to " + conn.LocalAddr().String())
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("All green, start processing")
			go FirstContact(conn)
		}
	}()
	return nil
}


func FirstContact(c net.Conn) {
	var out []byte

	buf := make([]byte, 0, 4096) // big buffer

	tmp := make([]byte, 1000) // using small tmo buffer for demonstrating
	n, err := c.Read(tmp)
	if err != nil {
		if err != io.EOF {
			fmt.Println("read error:", err)
		}
	}
	buf = append(buf, tmp[:n]...)

	fmt.Println("Start printing message******")
	for _, v := range buf {
		fmt.Print(string(v))
	}
	fmt.Printf("**Data finish, %v bytes total \n", n)
	var req request
	out = appendHandle(out, &req)
	_, _ = c.Write(out)
	return
	//data := iss.Begin(in)
	//var res = "Hello World!\r\n"
	/*
		if bytes.Contains(buf, []byte("\r\n\r\n")) {
			// for testing minimal single packet request -> response.
			out = appendResponse(nil, "200 OK", "", res)
			_, _ = c.Write(out)
			return
		}
	*/
	// process the pipeline
	/*
		for {
			leftover, err := parseRequest(buf, &req)
			if err != nil {
				// bad thing happened
				out = appendResponse(out, StatusBadRequest, "", err.Error()+"\n")
				_, _ = c.Write(out)
				_ = c.Close()
				break
			} else if len(leftover) == len(buf) {
				// request not ready, yet
				break
			}
			// handle the request
			req.remoteAddr = c.RemoteAddr().String()
			out = appendHandle(out, &req)
			buf = leftover
			_, _ = c.Write(out)
		}

	*/
}

func appendHandle(b []byte, req *request) []byte {
	return appendResponse(b, StatusOK, "", "Hello World!\r\n")
}

// appendResponse will append a valid http response to the provide bytes.
// The status param should be the code plus text such as "200 OK".
// The head parameter should be a series of lines ending with "\r\n" or empty.
func appendResponse(b []byte, status, head, body string) []byte {
	b = append(b, "HTTP/1.1"...)
	b = append(b, ' ')
	b = append(b, status...)
	b = append(b, '\r', '\n')
	b = append(b, "Server: Goose\r\n"...)
	b = append(b, "Date: "...)
	b = time.Now().AppendFormat(b, "Mon, 02 Jan 2006 15:04:05 GMT")
	b = append(b, '\r', '\n')
	if len(body) > 0 {
		b = append(b, "Content-Length: "...)
		b = strconv.AppendInt(b, int64(len(body)), 10)
		b = append(b, '\r', '\n')
	}
	b = append(b, head...)
	b = append(b, '\r', '\n')
	if len(body) > 0 {
		b = append(b, body...)
	}
	return b
}

// parseRequest is a very simple http request parser. This operation
// waits for the entire payload to be buffered before returning a
// valid request.
func parseRequest(data []byte, req *request) (leftover []byte, err error) {
	sdata := string(data)

	//int type initiate with value zero
	var i, s int
	var top string
	var clen int
	var q = -1
	// method, path, proto line
	for ; i < len(sdata); i++ {
		if sdata[i] == ' ' {
			req.method = sdata[s:i]
			for i, s = i+1, i+1; i < len(sdata); i++ {
				if sdata[i] == '?' && q == -1 {
					q = i - s
				} else if sdata[i] == ' ' {
					if q != -1 {
						req.path = sdata[s:q]
						req.query = req.path[q+1 : i]
					} else {
						req.path = sdata[s:i]
					}
					for i, s = i+1, i+1; i < len(sdata); i++ {
						if sdata[i] == '\n' && sdata[i-1] == '\r' {
							req.proto = sdata[s:i]
							i, s = i+1, i+1
							break
						}
					}
					break
				}
			}
			break
		}
	}
	if req.proto == "" {
		return data, fmt.Errorf("malformed request")
	}
	top = sdata[:s]
	for ; i < len(sdata); i++ {
		if i > 1 && sdata[i] == '\n' && sdata[i-1] == '\r' {
			line := sdata[s : i-1]
			s = i + 1
			if line == "" {
				req.head = sdata[len(top)+2 : i+1]
				i++
				if clen > 0 {
					if len(sdata[i:]) < clen {
						break
					}
					req.body = sdata[i : i+clen]
					i += clen
				}
				return data[i:], nil
			}
			if strings.HasPrefix(line, "Content-Length:") {
				n, err := strconv.ParseInt(strings.TrimSpace(line[len("Content-Length:"):]), 10, 64)
				if err == nil {
					clen = int(n)
				}
			}
		}
	}
	// not enough data
	return data, nil
}

//HTTP status codes were stolen from net/http.
//and there's some minor modifications
const (
	StatusContinue           = "100 Continue"            // RFC 7231, 6.2.1
	StatusSwitchingProtocols = "101 Switching Protocols" // RFC 7231, 6.2.2
	StatusProcessing         = "102 Processing"          // RFC 2518, 10.1

	StatusOK                   = "200 OK"                            // RFC 7231, 6.3.1
	StatusCreated              = "201 Created"                       // RFC 7231, 6.3.2
	StatusAccepted             = "202 Accepted"                      // RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo = "203 Non-Authoritative Information" // RFC 7231, 6.3.4
	StatusNoContent            = "204 No Content"                    // RFC 7231, 6.3.5
	StatusResetContent         = "205 Reset Content"                 // RFC 7231, 6.3.6
	StatusPartialContent       = "206 Partial Content"               // RFC 7233, 4.1
	StatusMultiStatus          = "207 Multi-Status"                  // RFC 4918, 11.1
	StatusAlreadyReported      = "208 Already Reported"              // RFC 5842, 7.1
	StatusIMUsed               = "226 IM Used"                       // RFC 3229, 10.4.1

	StatusMultipleChoices  = "300 Multiple Choices"  // RFC 7231, 6.4.1
	StatusMovedPermanently = "301 Moved Permanently" // RFC 7231, 6.4.2
	StatusFound            = "302 Found"             // RFC 7231, 6.4.3
	StatusSeeOther         = "303 See Other"         // RFC 7231, 6.4.4
	StatusNotModified      = "304 Not Modified"      // RFC 7232, 4.1
	StatusUseProxy         = "305 Use Proxy"         // RFC 7231, 6.4.5

	StatusTemporaryRedirect = "307 Temporary Redirect" // RFC 7231, 6.4.7
	StatusPermanentRedirect = "308 Permanent Redirect" // RFC 7538, 3

	StatusBadRequest                   = "400 Bad Request"                     // RFC 7231, 6.5.1
	StatusUnauthorized                 = "401 Unauthorized"                    // RFC 7235, 3.1
	StatusPaymentRequired              = "402 Payment Required"                // RFC 7231, 6.5.2
	StatusForbidden                    = "403 Forbidden"                       // RFC 7231, 6.5.3
	StatusNotFound                     = "404 Not Found"                       // RFC 7231, 6.5.4
	StatusMethodNotAllowed             = "405 Method Not Allowed"              // RFC 7231, 6.5.5
	StatusNotAcceptable                = "406 Not Acceptable"                  // RFC 7231, 6.5.6
	StatusProxyAuthRequired            = "407 Proxy Authentication Required"   // RFC 7235, 3.2
	StatusRequestTimeout               = "408 Request Timeout"                 // RFC 7231, 6.5.7
	StatusConflict                     = "409 Conflict"                        // RFC 7231, 6.5.8
	StatusGone                         = "410 Gone"                            // RFC 7231, 6.5.9
	StatusLengthRequired               = "411 Length Required"                 // RFC 7231, 6.5.10
	StatusPreconditionFailed           = "412 Precondition Failed"             // RFC 7232, 4.2
	StatusRequestEntityTooLarge        = "413 Payload Too Large"               // RFC 7231, 6.5.11
	StatusRequestURITooLong            = "414 URI Too Long"                    // RFC 7231, 6.5.12
	StatusUnsupportedMediaType         = "415 Unsupported Media Type"          // RFC 7231, 6.5.13
	StatusRequestedRangeNotSatisfiable = "416 Range Not Satisfiable"           // RFC 7233, 4.4
	StatusExpectationFailed            = "417 Expectation Failed"              // RFC 7231, 6.5.14
	StatusTeapot                       = "418 I'm a teapot"                    // RFC 7168, 2.3.3
	StatusUnprocessableEntity          = "422 Unprocessable Entity"            // RFC 4918, 11.2
	StatusLocked                       = "423 Locked"                          // RFC 4918, 11.3
	StatusFailedDependency             = "424 Failed Dependency"               // RFC 4918, 11.4
	StatusUpgradeRequired              = "426 Upgrade Required"                // RFC 7231, 6.5.15
	StatusPreconditionRequired         = "428 Precondition Required"           // RFC 6585, 3
	StatusTooManyRequests              = "429 Too Many Requests"               // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  = "431 Request Header Fields Too Large" // RFC 6585, 5
	StatusUnavailableForLegalReasons   = "451 Unavailable For Legal Reasons"   // RFC 7725, 3

	StatusInternalServerError           = "500 Error"                           // RFC 7231, 6.6.1
	StatusNotImplemented                = "501 Not Implemented"                 // RFC 7231, 6.6.2
	StatusBadGateway                    = "502 Bad Gateway"                     // RFC 7231, 6.6.3
	StatusServiceUnavailable            = "503 Service Unavailable"             // RFC 7231, 6.6.4
	StatusGatewayTimeout                = "504 Gateway Timeout"                 // RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported       = "505 HTTP Version Not Supported"      // RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates         = "506 Variant Also Negotiates"         // RFC 2295, 8.1
	StatusInsufficientStorage           = "507 Insufficient Storage"            // RFC 4918, 11.5
	StatusLoopDetected                  = "508 Loop Detected"                   // RFC 5842, 7.2
	StatusNotExtended                   = "510 Not Extended"                    // RFC 2774, 7
	StatusNetworkAuthenticationRequired = "511 Network Authentication Required" // RFC 6585, 6
)
