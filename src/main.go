package main

import "log"

func main() {
	test()
	test2()
}
func test() bool {
	log.Printf("yes")
	return true
}
func test2() bool {
	log.Printf("no")
	return false
}