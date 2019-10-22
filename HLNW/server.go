/******************************************************************
* N.Kozak // Lviv'2019 // example server Go(HLNW)-Asm for pkt2 SP *
*                         file: server.go                         *
*                                                          files: *
*                                                          calc.s *
*                                                       server.go *
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
    "fmt"
    "net/http"
)

func main() {
    html_code_fmt := 
    `<html>
    <body>
    <h2>X = K + B2 - D2/C1 + E1*F2</h2>
    <h2>--------------------------</h2>
    <h2>K = %d</h2>
    <h2>B = %f</h2>
    <h2>C = %f</h2>
    <h2>D = %f</h2>
    <h2>E = %f</h2>
    <h2>F = %f</h2>
    <h2>-------</h2>
    <h2>X(Assembly) = %f</h2>
    <h2>X(Go) = %f</h2>
    <h2>--------------------------</h1>
    </body>
    </html>`;

	var b2 float64 = 10 
	var c1 float32 = 20 
	var d2 float64 = 30 
	var e1 float32 = 40 
	var f2 float64 = 50 
	x_AssemblyResult := C.calc(C.double(b2), C.float(c1), C.double(d2), C.float(e1), C.double(f2))	
	x_GoResult := float64(0x00025630/*K*/) + b2 - d2/float64(c1) + float64(e1)*f2;
	println(x_AssemblyResult);
	println(x_GoResult);
	println("           Send response:");	
	fmt.Printf(html_code_fmt, 0x00025630/*K*/, b2, c1, d2, e1, f2, x_AssemblyResult, x_GoResult);
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, html_code_fmt, 0x00025630/*K*/, b2, c1, d2, e1, f2, x_AssemblyResult, x_GoResult)
	})
	http.ListenAndServe(":80", nil)	
	
}