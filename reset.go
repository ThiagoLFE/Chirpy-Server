package main

import (
	"net/http"
)

func (c *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	c.fileserverHits.Store(0)

	if c.PLATFORM != "dev" {
		respondWithError(w, http.StatusForbidden, "Permission denied, you haven't permission to do this")
		return
	}

	if err := c.db.ClearUsers(r.Context()); err != nil {
		respondWithError(w, http.StatusBadGateway, "Fail to delete all users data: "+err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, "Users deleted")
}
