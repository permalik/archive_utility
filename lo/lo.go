package lo

import "log"

func G(c int, m string, v interface{}) {
	if c == 0 {
		log.Printf("info::\n%s\n%v", m, v)
	} else if c == 1 {
		log.Fatalf("fatal::\n%s\n%v", m, v)
	} else if c == 2 {
		log.Panicf("panic::\n%s\n%v", m, v)
	}
}
