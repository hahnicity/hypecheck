package hypecheck

import (
    "fmt"
    "github.com/hahnicity/go-stringit"
    "math"
)


type Analyzer struct {
    response  *Response
    threshold float64
}

func (a *Analyzer) AnalyzeStock(resp *Response) {
    indices, difs := a.findLargePriceSwings()
    rets := a.findReturnsInTen(indices)
    fmt.Println("Symbol: ", a.response.Symbol)
    for i, index := range indices {
        fmt.Println(stringit.Format(
            "\tDate: {}, Dif: {}, Later: {}", a.response.Stock[index].Date, difs[i], rets[i],
        ))    
    }
}

// Find the dates after which a large price swing (denoted by the threshold variable)
// has occurred
func (a *Analyzer) findLargePriceSwings() (indices []int, difs []float64) {
    for i := 1; i < len(a.response.Stock); i++ {
        dif := math.Abs(a.response.Stock[i].Adj - a.response.Stock[i-1].Adj) / a.response.Stock[i].Adj
        if dif > a.threshold {
            indices = append(indices, i)
            difs = append(difs, dif)
        }
    }
    return
}

// Find the returns on a stock in ten trading days after a large price swing 
// has occurred
func (a *Analyzer) findReturnsInTen(indices []int) (rets []float64) {
    return   
}

