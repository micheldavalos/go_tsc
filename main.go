package main

import (
    "fmt"
    "log"
    "syscall"
	"unsafe"
)

func UintPtrToString(r uintptr) string {
    p := (*byte)(unsafe.Pointer(r))
    if p == nil {
        return ""
    }

	n, end, add := 0, unsafe.Pointer(p), unsafe.Sizeof(*p)
	for *(*byte)(end) != 0 {
		end = unsafe.Add(end, add)
        n++

	}
    return string(unsafe.Slice(p, n))
}

func main() {
    h, e := syscall.LoadLibrary("TSCLIB.dll")   // Make sure this DLL follows Golang machine bit architecture (64-bit in my case)
    if e != nil {
        log.Fatal(e)
    }
    defer syscall.FreeLibrary(h)
    
    proc, e := syscall.GetProcAddress(h, "openport") // One of the functions in the DLL
    if e != nil {
        log.Fatal(e)
    }
    
	s := "USB" // port, connection or printer's name shown in Windows
	b := append([]byte(s), 0)


    n, _, _ := syscall.Syscall(proc, uintptr(1), uintptr(unsafe.Pointer(&b[0])), uintptr(0), uintptr(0) )  // Pay attention to the positioning of the parameter
    fmt.Printf("Connection %d\n", n)  // successful = 1, fail = -1

	proc, e = syscall.GetProcAddress(h, "usbprinterserial") // One of the functions in the DLL
    if e != nil {
        log.Fatal(e)
    }

	n, _, _ = syscall.Syscall(proc, uintptr(0), uintptr(0), uintptr(0), uintptr(0) )  
    serial_number := UintPtrToString(n) // convert []byte in string, successful = a serial number, fail = -1
	fmt.Printf(serial_number)  

	
}