package node

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Node struct {
	chain  *Chain
	router *mux.Router
}

func (n *Node) Router() http.Handler {
	return n.router
}
