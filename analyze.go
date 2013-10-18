package hypecheck

import (
    "fmt"
    "github.com/hahnicity/go-stringit"
    "math"
)


type Analyzer struct {
    days      int
    response  *Response
    threshold float64
}

// Analyze the stock data
func (a *Analyzer) AnalyzeStock(resp *Response) {
    indices, swings := a.findLargePriceSwings()
    rets := a.findReturnsAfterSwing(indices)
    fmt.Println("Symbol: ", a.response.Symbol)
    for i, index := range indices {
        fmt.Println(stringit.Format(
            "\tDate: {}, Dif: {}, Later: {}", a.response.Stock[index].Date, swings[i], rets[i],
        ))    
    }
}

// Find the dates after which a large price swing (denoted by the threshold variable)
// has occurred
func (a *Analyzer) findLargePriceSwings() (indices []int, swings []float64) {
    for i := 1; i < len(a.response.Stock); i++ {
        swing := (a.response.Stock[i-1].Adj - a.response.Stock[i].Adj) / a.response.Stock[i].Adj
        if math.Abs(swing) > a.threshold {
            indices = append(indices, i)
            swings = append(swings, swing)
        }
    }
    return
}

// Find the returns on a stock in <a.days> trading days after a large price swing 
// has occurred
func (a *Analyzer) findReturnsAfterSwing(indices []int) (rets []float64) {
    defer func() {
        // 
        if r := recover(); r != nil {
            err, ok := r.(error)
            if !ok {
                fmt.Println("An Error occurred but the program recovered. Error: ", err)    
            }    
        }     
    }()
    for _, index := range indices {
        ret := (a.response.Stock[index+a.days].Adj - a.response.Stock[index].Adj) / a.response.Stock[index].Adj
        rets = append(rets, ret)
    }
    return  
}

