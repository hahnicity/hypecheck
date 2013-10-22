package hypecheck

import (
    "encoding/csv"
    "fmt"
    "github.com/grd/histogram"
    "os"
    "strconv"
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

// Analyze all significant dates with a histogram. But also write any significant data
// to a csv file so we can do further processing in R
func AnalyzeAllResponses(a *Analyzer, ar []*Response) {
    f, err := os.Create("swing-data.csv")
    if err != nil { panic(err) }
    defer f.Close()
    w := csv.NewWriter(f)
    defer w.Flush()
    Range := histogram.Range(-1.0, 200, .01)
    h, err := histogram.NewHistogram(Range)
    if err != nil {
        panic(err)
    }
    for _, resp := range ar {
        for _, oi := range a.AnalyzeStock(resp) {
            var toWrite = []string{
                strconv.FormatFloat(oi.Swing, 'f', 4, 64),
                strconv.FormatFloat(oi.Ret, 'f', 4, 64),
            }
            w.Write(toWrite)
            h.Add(oi.Ret)
        }
    }
    fmt.Println("MEAN: ", h.Mean())
    fmt.Println("SIGMA ", h.Sigma())
}
 
func NewAnalyzer(days int, threshold float64) *Analyzer {
    return &Analyzer{days, threshold}
}

// Analyze the stock data
func (a *Analyzer) AnalyzeStock(resp *Response) (ois []OfInterest) {
    ois = a.findLargePriceSwings(resp)
    a.findReturnsAfterSwing(&ois, resp)
    filterNullRets(&ois)
    return
}

// Find the dates after which a large price swing (denoted by the threshold variable)
// has occurred
func (a *Analyzer) findLargePriceSwings(resp *Response) (ois []OfInterest) {
    for i := 1; i < len(resp.Stock); i++ {
        swing := (resp.Stock[i].Adj - resp.Stock[i-1].Adj) / resp.Stock[i-1].Adj
        if a.threshold > 0.0 && swing > a.threshold {
            // XXX Abstract this
            oi := new(OfInterest)
            oi.Index = i
            oi.Stock = resp.Stock[i]
            oi.Swing = swing
            ois = append(ois, *oi)
        } else if a.threshold < 0.0 && swing < a.threshold {
            // XXX Abstract this
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
            return    
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
