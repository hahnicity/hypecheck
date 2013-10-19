//Copyright (c) 2013 Michael Dvorkin. All Rights Reserved.
//"mike" + "@dvorkin" + ".net" || "twitter.com/mid"

// Took substantitial portions of code from https://github.com/michaeldv/mop for this
package hypecheck

import (
    "bytes"
    "github.com/hahnicity/go-stringit"
    "github.com/hahnicity/hypecheck/config"
    "io/ioutil"
    "net/http"
    "net/url"
    "reflect"
    "strconv"
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
    Stock  []Stock
    Symbol string
}

type Stock struct {
    Date        string 
    Open        float64  // market open price.
    High        float64  // day's high.
    Low         float64  // day's low.
    Close       float64  // closing price
    Volume      int      // volume.
    Adj         float64  // closing price adjusted for inflation and other junk
}

type Request struct {
    Options  *Options
    Response chan *Response
    Symbol   string
}

func NewRequest(c chan *Response, symbol string, values map[string]interface{}) (r *Request) {
    r = new(Request)
    r.Response = c
    r.Symbol = symbol
    r.Options = NewOptions()
    for k, v := range values {
        r.Options.Add(k, v)    
    }
    return
}

// Execute the request to the Yahoo finance API. Get the necessary data and then return it
// to the requester object so it can be processed
func (r *Request) Execute() (resp *Response) {
    resp = new(Response)
    url := stringit.Format(config.BaseUrl, r.Symbol, r.Options.Values.Encode())
    response, err := http.Get(url)
    if err != nil {
        panic(err)
    }
    defer response.Body.Close()
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        panic(err)
    }
    resp.Stock = r.parse(sanitize(body))
    resp.Symbol = r.Symbol
    return
}

// Use reflection to parse and assign the quotes data fetched using the Yahoo
// market API.
func (r *Request) parse(body []byte) (stocks []Stock) {
    lines := bytes.Split(body, []byte{'\n'})[1:] // Cut off the header
    stocks = make([]Stock, len(lines))
    fieldsCount := reflect.ValueOf(stocks[0]).NumField()
    // Split each line into columns, then iterate over the Stock struct
    // fields to assign column values.
    for i, line := range lines {
        columns := bytes.Split(bytes.TrimSpace(line), []byte{','})
        for j := 0; j < fieldsCount; j++ {
            switch reflect.TypeOf(&stocks[i]).Elem().Field(j).Type.String() {
            case "string":
                reflect.ValueOf(&stocks[i]).Elem().Field(j).SetString(string(columns[j]))
            case "float64":
                f, err := strconv.ParseFloat(string(columns[j]), 64)
                if err != nil {
                    panic(err)    
                }
                reflect.ValueOf(&stocks[i]).Elem().Field(j).SetFloat(f)
            case "int":
                x, err := strconv.ParseInt(string(columns[j]), 10, 64)
                if err != nil {
                    panic(err)    
                }
                reflect.ValueOf(&stocks[i]).Elem().Field(j).SetInt(x)
            }
        }
    }
    return 
}

//-----------------------------------------------------------------------------
func sanitize(body []byte) []byte {
    return bytes.Replace(bytes.TrimSpace(body), []byte{'"'}, []byte{}, -1)
}
