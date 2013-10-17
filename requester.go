package hypecheck

type Requester struct {
    // The number of active requests that are being run with wikipedia
    activeRequests int

    // All response received from wikipedia
    allResponses   []*Response

    // Close all wikipedia requests after they are made
    closeRequests  bool

    // The maximum number of requests that can be run concurrently
    maxRequests    int

    // Channel of active work
    Work           chan Request

    // The year to analyze requests
    year           string
}

// Make a new Requester object. 
func NewRequester(closeRequests bool, maxRequests int, year string) (r *Requester){
    r = new(Requester)    
    r.activeRequests = 0
    r.closeRequests = closeRequests
    r.maxRequests = maxRequests
    r.Work = make(chan Request)
    r.year = year
    return
}

// Given a map of companies and their corresponding wikipedia pages, make
// requests to stats.grok.se so that we can get statistics as to how frequently
// people are viewing their pages
func (r *Requester) MakeRequests(companies map[string]string) {
    c := make(chan *Response)
    for symbol, _ := range companies {
        r.activeRequests++
        r.Work <- *NewRequest(symbol, nil)
        r.manageActiveProc(c)
    }
    r.waitToFinish(c, companies)
}

// Throttle number of active requests stats.grok.se is the problem here
func (r *Requester) manageActiveProc(c chan *Response) {
    if r.activeRequests == r.maxRequests {
        resp := <- c
        r.allResponses = append(r.allResponses, resp)
    }
    r.activeRequests--
}

// Wait for all requests to finish
func (r *Requester) waitToFinish(c chan *Response, companies map[string]string) {
    for len(r.allResponses) < len(companies) {
        r.allResponses = append(r.allResponses, <-c)
    }
}
