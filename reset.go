package main

import (
	"fmt"
	"net/http"
)

func (c *apiConfig) handlerReset(w http.ResponseWriter, _ *http.Request) {
	c.fileserverHits.Store(0)
	w.Write(fmt.Appendf([]byte{}, "reseted hits"))

}
