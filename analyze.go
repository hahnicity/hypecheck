package hypecheck

import (
    "fmt"
    //"github.com/hahnicity/go-stringit"
    "github.com/VividCortex/gohistogram"
    "math"
)

type OfInterest struct {
    Index int
    Ret   float64
    Stock Stock
    Swing float64
}

type Analyzer struct {
    days      int
    threshold float64
}

func NewAnalyzer(days int, threshold float64) *Analyzer {
    return &Analyzer{days, threshold}
}

// Analyze the stock data
func (a *Analyzer) AnalyzeStock(resp *Response) {
    ois := a.findLargePriceSwings(resp)
    a.findReturnsAfterSwing(&ois, resp)
    filterNullRets(&ois)
    nh := gohistogram.NewHistogram(20)
    fmt.Println("Symbol: ", resp.Symbol)
    for _, oi := range ois {
        nh.Add(oi.Ret)
        //fmt.Println(stringit.Format(
        //    "\tDate: {}, Dif: {}, Later: {}", oi.Stock.Date, oi.Swing, oi.Ret,
        //))    
    }
    for i := 1; i < 10; i++ {
        fmt.Println("\t.", i, " Quantile:", nh.Quantile(float64(i) * 0.1))
    }
}

// Find the dates after which a large price swing (denoted by the threshold variable)
// has occurred
func (a *Analyzer) findLargePriceSwings(resp *Response) (ois []OfInterest) {
    for i := 1; i < len(resp.Stock); i++ {
        swing := (resp.Stock[i-1].Adj - resp.Stock[i].Adj) / resp.Stock[i].Adj
        if math.Abs(swing) > a.threshold {
            oi := new(OfInterest)
            oi.Index = i
            oi.Stock = resp.Stock[i]
            oi.Swing = swing
            ois = append(ois, *oi)
        }
    }
    return
}

// Find the returns on a stock in <a.days> trading days after a large price swing 
// has occurred
func (a *Analyzer) findReturnsAfterSwing(ois *[]OfInterest, resp *Response) {
    defer func() {
        if r := recover(); r != nil {
            err, ok := r.(error)
            if ok {
                fmt.Println("An Error occurred but the program recovered. Error: ", err)    
            }    
        }     
    }()
    for i, oi := range *ois {
        ret := (resp.Stock[oi.Index + a.days].Adj - oi.Stock.Adj) / oi.Stock.Adj
        (&oi).Ret = ret
        (*ois)[i] = oi
    }
    return 
}

/* ----------------------------------------------------------------------------- */
func filterNullRets (ois *[]OfInterest) {
    for i, oi := range *ois {
        if oi.Ret == 0.0 {
            *ois = (*ois)[:i]
            break
        }
    }    
}
