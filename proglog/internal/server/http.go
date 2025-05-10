package server

import (
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
)

type ProduceRequest struct {
    Record Record `json:"record"`
}
type ProduceResponse struct {
    Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
    Offset uint64 `json:"offset"`
}
type ConsumeResponse struct {
    Record Record `json:"record"`
}   

type customHttpServer struct {
    Log *Log
}

func newCustomHttpServer(log *Log) *customHttpServer {
    return &customHttpServer{
        Log: log,
    }
}

func (s *customHttpServer) handleProduce(w http.ResponseWriter, r *http.Request){
    var req ProduceRequest
    if err:= json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    offset, err:=s.Log.Append(req.Record)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    resp := ProduceResponse{Offset: offset}
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func (s *customHttpServer) handleConsume(w http.ResponseWriter, r *http.Request){
    var req ConsumeRequest
    if err:= json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    record, err:=s.Log.Read(req.Offset)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    resp := ConsumeResponse{Record: record}
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func NewHttpServer(addr string) *http.Server{

    svr:=newCustomHttpServer(
        newLog(),
    )
    r:=mux.NewRouter()
    r.HandleFunc("/", svr.handleProduce).Methods("POST")
    r.HandleFunc("/", svr.handleConsume).Methods("GET")

    return &http.Server{
        Addr: addr,
        Handler: r,
    }
}