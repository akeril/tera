package main

import "fmt"

func injectLink(html []byte) []byte {
	script := "<script src='/tera'></script>"
	inject := fmt.Sprintf("<head>%v</head>", script)
	return []byte(inject + string(html))
}
