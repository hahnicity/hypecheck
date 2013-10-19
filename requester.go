package hypecheck

import "time"


type Requester struct {
    // The analyzer struct we will use to analyze data on the requests
    a              *Analyzer

    // The number of active requests that are being run with wikipedia
    activeRequests int

    // All response received from wikipedia
    allResponses   []*Response

    // The maximum number of requests that can be run concurrently
    maxRequests    int
    
    // The delay in milliseconds between each request
    requestDelay   time.Duration

    // Channel of active work
    Work           chan Request
}

// Make a new Requester object. 
func NewRequester(a *Analyzer, maxRequests, requestDelay int) (r *Requester){
    r = new(Requester)
    r.a = a
    r.activeRequests = 0
    r.maxRequests = maxRequests
    r.requestDelay = time.Duration(requestDelay) * time.Millisecond
    r.Work = make(chan Request)
    return
}

// Given a map of companies and their corresponding wikipedia pages, make
// requests to stats.grok.se so that we can get statistics as to how frequently
// people are viewing their pages
func (r *Requester) MakeRequests(companies []string) {
    c := make(chan *Response)
    for _, symbol := range companies {
        r.activeRequests++
        r.Work <- *NewRequest(c, symbol, nil)
        // If we need to wait for a request to finish do not implement an additional delay
        if r.manageActiveProc(c) {
            continue    
        } else {
            time.Sleep(r.requestDelay)       
        }
    }
    r.waitToFinish(c, companies)
    for _, resp := range r.allResponses {
        r.a.AnalyzeStock(resp)    
    }
}

// Throttle number of active requests if we are at the number of requests
// that we have allowed to run concurrently
func (r *Requester) manageActiveProc(c chan *Response) bool {
    if r.activeRequests == r.maxRequests {
        resp := <- c
        r.allResponses = append(r.allResponses, resp)
        r.activeRequests--
        return true
    }
    return false
}

// Wait for all requests to finish
func (r *Requester) waitToFinish(c chan *Response, companies []string) {
    for len(r.allResponses) < len(companies) {
        r.allResponses = append(r.allResponses, <-c)
    }
    close(c)
}
