package httpServer

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type request struct {
	proto, method string
	path, query   string
	head, body    string
	remoteAddr    string
}

/*
func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:666")
	fmt.Println("Listening on: " + ln.Addr().String())
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := ln.Accept()
		fmt.Println("incoming " + conn.RemoteAddr().String() + " to " + conn.LocalAddr().String())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("All green, start processing")
		go Sucker(conn)
	}
}


 */
func Sucker(c net.Conn) {
	var out []byte
	//is := c.Context().(*InputStream)

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
	out = appendhandle(out, &req)
	_, _ = c.Write(out)
	return
	//data := iss.Begin(in)
	//var res = "Hello World!\r\n"
	/*
		if bytes.Contains(buf, []byte("\r\n\r\n")) {
			// for testing minimal single packet request -> response.
			out = appendresp(nil, "200 OK", "", res)
			_, _ = c.Write(out)
			return
		}
	*/
	// process the pipeline
	/*
		for {
			leftover, err := parsereq(buf, &req)
			if err != nil {
				// bad thing happened
				out = appendresp(out, "500 Error", "", err.Error()+"\n")
				_, _ = c.Write(out)
				_ = c.Close()
				break
			} else if len(leftover) == len(buf) {
				// request not ready, yet
				break
			}
			// handle the request
			req.remoteAddr = c.RemoteAddr().String()
			out = appendhandle(out, &req)
			buf = leftover
			_, _ = c.Write(out)
		}

	*/
}

func appendhandle(b []byte, req *request) []byte {
	return appendresp(b, "200 OK", "", "Hello World!\r\n")
}

// appendresp will append a valid http response to the provide bytes.
// The status param should be the code plus text such as "200 OK".
// The head parameter should be a series of lines ending with "\r\n" or empty.
func appendresp(b []byte, status, head, body string) []byte {
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

// parsereq is a very simple http request parser. This operation
// waits for the entire payload to be buffered before returning a
// valid request.
func parsereq(data []byte, req *request) (leftover []byte, err error) {
	sdata := string(data)
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
