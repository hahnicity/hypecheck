package hypecheck

// XXX DEBUG
import "fmt"


type Requester struct {
    // The number of active requests that are being run with wikipedia
    activeRequests int

    // All response received from wikipedia
    allResponses   []*Response

    // The maximum number of requests that can be run concurrently
    maxRequests    int

    // Channel of active work
    Work           chan Request
}

// Make a new Requester object. 
func NewRequester(maxRequests int) (r *Requester){
    r = new(Requester)    
    r.activeRequests = 0
    r.maxRequests = maxRequests
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
        fmt.Println("MAKE NEW REQUEST")
        fmt.Println("CHANNEL", c)
        r.Work <- *NewRequest(c, symbol, nil)
        r.manageActiveProc(c)
    }
    r.waitToFinish(c, companies)
}

// Throttle number of active requests stats.grok.se is the problem here
func (r *Requester) manageActiveProc(c chan *Response) {
    if r.activeRequests == r.maxRequests {
        resp := <- c
        r.allResponses = append(r.allResponses, resp)
        r.activeRequests--
    }
}

// Wait for all requests to finish
func (r *Requester) waitToFinish(c chan *Response, companies []string) {
    fmt.Println("WAITING TO FINISH")
    for len(r.allResponses) < len(companies) {
        fmt.Println("ALL RESPONSES ", r.allResponses)
        r.allResponses = append(r.allResponses, <-c)
    }
}
