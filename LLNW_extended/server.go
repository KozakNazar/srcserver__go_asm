/******************************************************************
* N.Kozak // Lviv'2019 // example server Go(LLNW)-Asm for pkt2 SP *
*                         file: server.go                         *
*                                                          files: *
*                                                          calc.s *
*                                                       server.go *
*                                                                 *
*                                               *extended version *
*******************************************************************/
package main

// #include <stdio.h>
// #include <stdlib.h>
// #cgo LDFLAGS: calc.o
// float calc(double b2, float c1, double d2, float e1, double f2);
import (
     "C"
     //"unsafe" // for unsafe.Pointer()
)
import (
    //"io"
    //"io/ioutil"    
    "fmt"
    "net"
    "os"
    "bufio"	
    "strings"		
)

var usePostSubmit bool = true; // defaault value

var b2 float64 = 10            // defaault value
var c1 float32 = 20            // defaault value 
var d2 float64 = 30            // defaault value 
var e1 float32 = 40            // defaault value 
var f2 float64 = 50            // defaault value 


func buildResponse() []byte {
    const K int32 = 0x00025630     // const value 

    http_response_fmt := 
`HTTP/1.1 200 OK
Date: Mon, 27 Jul 2009 12:28:53 GMT
Server: Apache/2.2.14 (Win32)
Last-Modified: Wed, 22 Jul 2009 19:15:56 GMT
Content-Length: %d
Status: 200
Content-Type: text/html
Connection: Closed

%s`;

    html_code_fmt__withGetSubmit := 
`<html>
<head>
<link rel="icon" href="data:;base64,=">
</head>
<body>

<form action="/setSettings" method="get">
<h1>Settings:</h1>
<p>Select mode:</p>  
  <input type="radio" name="mode" value="1"> mode 1<br>
  <input type="radio" name="mode" value="2"> mode 2<br>
  <input type="radio" name="mode" value="3"> mode 3<br> 
<p>Change used http-method:</p>
  <input type="radio" name="http_method" value="0" checked="checked"> GET<br>
  <input type="radio" name="http_method" value="1"> POST<br> 

  <input type="submit" value="Submit parameters and reload page" text="">
</form>
<h1>Compute board:</h1>
<button type="submit" form="calcData">Send values by GET and compute result</button>

<form id="calcData"  method="get" action="callCalc">

<h2>X = K + B2 - D2/C1 + E1*F2</h2>
<h2>--------------------------</h2>
<h2>K = %d</h2>
<h2>B = <input name="B" type="text" value="%f"></h2>
<h2>C = <input name="C" type="text" value="%f"></h2>
<h2>D = <input name="D" type="text" value="%f"></h2>
<h2>E = <input name="E" type="text" value="%f"></h2>
<h2>F = <input name="F" type="text" value="%f"></h2>
<h2>-------</h2>
<h2>X(Assembly) = %f</h2>
<h2>X(Go) = %f</h2>
<h2>--------------------------</h2>
<input type="submit" value="Send values by GET and compute result">

</form>
</body>
</html>`;	

    html_code_fmt__withPostSubmit := 
`<html>
<head>
<link rel="icon" href="data:;base64,=">
</head>
<body>

<form action="/setSettings" method="post">
<h1>Settings:</h1>
<p>Select mode:</p>  
  <input type="radio" name="mode" value="1"> mode 1<br>
  <input type="radio" name="mode" value="2"> mode 2<br>
  <input type="radio" name="mode" value="3"> mode 3<br> 
<p>Change used http-method:</p>
  <input type="radio" name="http_method" value="0"> GET<br>
  <input type="radio" name="http_method" value="1" checked="checked"> POST<br> 

  <input type="submit" value="Submit parameters and reload page" text="">
</form>
<h1>Compute board:</h1>
<button type="submit" form="calcData">Send values by POST and compute result</button>

<form id="calcData"  method="post" action="callCalc">

<h2>X = K + B2 - D2/C1 + E1*F2</h2>
<h2>--------------------------</h2>
<h2>K = %d</h2>
<h2>B = <input name="B" type="text" value="%f"></h2>
<h2>C = <input name="C" type="text" value="%f"></h2>
<h2>D = <input name="D" type="text" value="%f"></h2>
<h2>E = <input name="E" type="text" value="%f"></h2>
<h2>F = <input name="F" type="text" value="%f"></h2>
<h2>-------</h2>
<h2>X(Assembly) = %f</h2>
<h2>X(Go) = %f</h2>
<h2>--------------------------</h2>
<input type="submit" value="Send values by POST and compute result">

</form>
</body>
</html>`;

    html_code_fmt := html_code_fmt__withGetSubmit;

    if(usePostSubmit){
        html_code_fmt = html_code_fmt__withPostSubmit;
    }

    x_AssemblyResult := C.calc(C.double(b2), C.float(c1), C.double(d2), C.float(e1), C.double(f2))  
    x_GoResult := float64(K) + b2 - d2/float64(c1) + float64(e1)*f2;
    println(x_AssemblyResult);
    println(x_GoResult);
    println("           Build response:");
    
	var html_code, http_response string;	
	html_code = fmt.Sprintf(html_code_fmt, K, b2, c1, d2, e1, f2, x_AssemblyResult, x_GoResult);	
	http_response = fmt.Sprintf(http_response_fmt, len(html_code), html_code);
	
	fmt.Print(http_response)
	return []byte(http_response);
}

func handleClient(conn net.Conn){
    body := make([]byte, 2048)
    bufio.NewReader(conn).Read(body)
    var message string = string(body)

    http_method_key := "http_method="; 
    B_key := "B=";
    C_key := "C=";
    D_key := "D=";
    E_key := "E=";
    F_key := "F="; 	
	
	if isPOST, indexPOSTValues := strings.Index(message, "POST") > -1, strings.Index(message, "\r\n\r\n"); isPOST && indexPOSTValues > -1 {
        usePostSubmit = true

        if index := strings.Index(message[indexPOSTValues:], http_method_key); index > -1 {		
		    index += indexPOSTValues;		
            var usePostSubmitValue int32 = 0;
            if(usePostSubmit) {
                usePostSubmitValue = 1;
            }
            fmt.Sscanf(message[index + len(http_method_key):], "%d", &usePostSubmitValue)
            usePostSubmit = true;
            if(usePostSubmitValue == 0) {
                usePostSubmit = false;
            }       
        }                   

        if index := strings.Index(message[indexPOSTValues:], B_key); index > -1 {
		    index += indexPOSTValues;		
            fmt.Sscanf(message[index + len(B_key):], "%f", &b2)   
        }
        if index := strings.Index(message[indexPOSTValues:], C_key); index > -1 {
		    index += indexPOSTValues;		
            fmt.Sscanf(message[index + len(C_key):], "%f", &c1)   
        }
        if index := strings.Index(message[indexPOSTValues:], D_key); index > -1 {
		    index += indexPOSTValues;		
            fmt.Sscanf(message[index + len(D_key):], "%f", &d2)   
        }
        if index := strings.Index(message[indexPOSTValues:], E_key); index > -1 {
		    index += indexPOSTValues;		
            fmt.Sscanf(message[index + len(E_key):], "%f", &e1)   
        }
        if index := strings.Index(message[indexPOSTValues:], F_key); index > -1 {
		    index += indexPOSTValues;		
            fmt.Sscanf(message[index + len(F_key):], "%f", &f2)   
        }

    } else {
	
        if index := strings.Index(message, http_method_key); index > -1 {
            var usePostSubmitValue int32 = 0;
            if(usePostSubmit) {
                usePostSubmitValue = 1;
            }
            fmt.Sscanf(message[index + len(http_method_key):], "%d", &usePostSubmitValue)
            usePostSubmit = true;	
            if(usePostSubmitValue == 0) {		
                usePostSubmit = false;
            }       
        } else {
            if index := strings.Index(message, B_key); index > -1 {
                usePostSubmit = false
                fmt.Sscanf(message[index + len(B_key):], "%f", &b2)   
            }
            if index := strings.Index(message, C_key); index > -1 {
                usePostSubmit = false        
                fmt.Sscanf(message[index + len(C_key):], "%f", &c1)   
            }
            if index := strings.Index(message, D_key); index > -1 {
                usePostSubmit = false        
                fmt.Sscanf(message[index + len(D_key):], "%f", &d2)   
            }
            if index := strings.Index(message, E_key); index > -1 {
                usePostSubmit = false        
                fmt.Sscanf(message[index + len(E_key):], "%f", &e1)   
            }
            if index := strings.Index(message, F_key); index > -1 {
                usePostSubmit = false        
                fmt.Sscanf(message[index + len(F_key):], "%f", &f2)   
            }        
		}		
    }

    fmt.Println()
    fmt.Println("#########Message Received:##")
    fmt.Println(message)
    fmt.Println("############################")      
    fmt.Println()

    conn.Write(buildResponse())

    conn.Close()
}

func main() {
	tcp_port := ":80"
    tcpAddr, err := net.ResolveTCPAddr("tcp4", tcp_port)
    checkError(err)
    listener, err := net.ListenTCP("tcp", tcpAddr)
    checkError(err)
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        go handleClient(conn)
    }		
}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}