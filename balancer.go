// The ideas in this code and the load balancer are taken from
// http://concur.rspace.googlecode.com/hg/talk/concur.html#slide-51
package discern

import "container/heap"


type Worker struct {
    requests chan Request
    index    int
}

func (w *Worker) work(done chan *Worker) {
    req := <- w.requests
    req.Resp <- req.Execute()
    done <- w
    return
}

type Pool []*Worker

func (p Pool) Len() int { 
    return len(p) 
}

func (p Pool) Less(i, j int) bool {
    return p[i].index < p[j].index
}

func (p Pool) Swap(i, j int) { 
    p[i], p[j] = p[j], p[i] 
}

func (p *Pool) Push(x interface{}) {
    x.(*Worker).index = p.Len()
    *p = append(*p, x.(*Worker))    
}

func (p *Pool) Pop() interface{} {
    old := *p
    n := len(old)
    x := old[n-1]
    *p = old[0 : n-1]
    return x    
}

type Balancer struct {
    Pool *Pool    
    done chan *Worker
}

func (b *Balancer) Balance(work chan WikiRequest) {
    for {
        select {
        case req := <-work: // received a Request...
            b.dispatch(req) // ...so send it to a Worker
        case w := <- b.done:
            b.completed(w)
        }
    }       
}

// Job is complete; update heap
func (b *Balancer) completed(w *Worker) {
    // Put it into its place on the heap.
    heap.Push(b.Pool, w)
}

// Pull a worker off the heap, and send it to work
func (b *Balancer) dispatch(req WikiRequest) {
    w := heap.Pop(b.Pool).(*Worker)
    go w.work(b.done) // tell the task to get to work
    w.requests <- req  // send it to the task
}

// Constructor method for the Load Balancer
func MakeBalancer(n int) *Balancer {
    b := &Balancer{makePool(n), make(chan *Worker)}
    heap.Init(b.Pool) //initialize the pool
    return b
}

// Make the pool of workers
func makePool(n int) *Pool {
    p := new(Pool)
    for i := 0; i < n; i++ {
        requests := make(chan WikiRequest)
        p.Push(&Worker{requests, i})
    }
    return p
}
