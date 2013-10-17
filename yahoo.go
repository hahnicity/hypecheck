//Copyright (c) 2013 Michael Dvorkin. All Rights Reserved.
//"mike" + "@dvorkin" + ".net" || "twitter.com/mid"

// Took substantitial portions of code from https://github.com/michaeldv/mop for this
package hypecheck

import (
    "bytes"
    "fmt"
    "github.com/hahnicity/go-stringit"
    "github.com/hahnicity/hypecheck/config"
    "io/ioutil"
    "net/http"
    "net/url"
    "reflect"
)

type Options struct {
    Values url.Values    
}

func NewOptions() *Options {
    return &Options{make(url.Values)}    
}

// Really just a convenience method. All possible optional values we are able to add
// can be found at http://www.gummy-stuff.org/Yahoo-data.htm
func (o *Options) Add(name string, value interface{}) {
    o.Values.Add(name, stringit.Format("{}", value))    
}

type Response struct {
    Stock Stock
}

type Stock struct {
    Ticker      string  // Stock ticker.
    LastTrade   string  // l1: last trade.
    Change      string  // c6: change real time.
    ChangePct   string  // k2: percent change real time.
    Open        string  // o: market open price.
    Low         string  // g: day's low.
    High        string  // h: day's high.
    Low52       string  // j: 52-weeks low.
    High52      string  // k: 52-weeks high.
    Volume      string  // v: volume.
    AvgVolume   string  // a2: average volume.
    PeRatio     string  // r2: P/E ration real time.
    PeRatioX    string  // r: P/E ration (fallback when real time is N/A).
    Dividend    string  // d: dividend.
    Yield       string  // y: dividend yield.
    MarketCap   string  // j3: market cap real time.
    MarketCapX  string  // j1: market cap (fallback when real time is N/A).
    Advancing   bool    // True when change is >= $0.
}

type Request struct {
    Options  *Options
    Response chan *Response
    Symbol   string
}

func NewRequest(symbol string, values map[string]interface{}) (r *Request) {
    r = new(Request)
    r.Response = make(chan *Response)
    r.Symbol = symbol
    r.Options = NewOptions()
    r.Options.Add("s", symbol)
    for k, v := range values {
        r.Options.Add(k, v)    
    }
    return
}

// Execute the request. Save the csv with the same name as the symbol we are looking up
func (r *Request) Execute() (resp *Response) {
    resp = new(Response)
    url := stringit.Format(config.BaseUrl, r.Options.Values.Encode())
    response, err := http.Get(url)
    if err != nil {
        panic(err)
    }
    defer response.Body.Close()
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        panic(err)
    }
    resp.Stock = r.parse(body)
    return
}

// Use reflection to parse and assign the quotes data fetched using the Yahoo
// market API.
func (r *Request) parse(body []byte) (stock Stock) {
    lines := bytes.Split(body, []byte{'\n'})[0:1]
    stocks := make([]Stock, 1)
    //
    // Get the total number of fields in the Stock struct. Skip the last
    // Advanicing field which is not fetched.
    //
    fieldsCount := reflect.ValueOf(stocks[0]).NumField() - 1
    //
    // Split each line into columns, then iterate over the Stock struct
    // fields to assign column values.
    //
    for i, line := range lines {
        columns := bytes.Split(bytes.TrimSpace(line), []byte{','})
        fmt.Println("LINES ", lines)
        fmt.Println("STOCKS ", stocks)
        fmt.Println("COLUMNS ", columns, " ", len(columns))
        fmt.Println("FIELDS ", fieldsCount)
        for j := 0; j < fieldsCount; j++ {
            // ex. quotes.stocks[i].Ticker = string(columns[0])
            reflect.ValueOf(&stocks[i]).Elem().Field(j).SetString(string(columns[j]))
        }
        //
        // Try realtime value and revert to the last known if the
        // realtime is not available.
        //
        if stocks[i].PeRatio == `N/A` && stocks[i].PeRatioX != `N/A` {
            stocks[i].PeRatio = stocks[i].PeRatioX
        }
        if stocks[i].MarketCap == `N/A` && stocks[i].MarketCapX != `N/A` {
            stocks[i].MarketCap = stocks[i].MarketCapX
        }
        //
        // Stock is advancing if the change is not negative (i.e. $0.00
        // is also "advancing").
        //
        stocks[i].Advancing = (stocks[i].Change[0:1] != `-`)
    }
    stock = stocks[0]
    return 
}

//-----------------------------------------------------------------------------
func sanitize(body []byte) []byte {
    return bytes.Replace(bytes.TrimSpace(body), []byte{'"'}, []byte{}, -1)
}
