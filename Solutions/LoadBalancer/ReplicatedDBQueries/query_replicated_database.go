package main

type Conn interface {
	DoQuery(query string) Result
}

type Result interface {

}

func Query(conns []Conn, query string) Result {
	ch := make(chan Result, len(conns))  // buffered
	for _, conn := range conns {
		go func(c Conn) {
			ch <- c.DoQuery(query)
		}(conn)
	}
	return <-ch
}
