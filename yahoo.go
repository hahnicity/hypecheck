package hypecheck

import (
    "github.com/hahnicity/go-stringit"
    "github.com/hahnicity/go-wget"
    "github.com/hahnicity/hypecheck/config"
    "net/url"
)

type Options struct {
    Options url.Values    
}

func NewOptions() *Options {
    return &Options{make(url.Values)}    
}

// Really just a convenience method. All possible optional values we are able to add
// can be found at http://www.gummy-stuff.org/Yahoo-data.htm
func (o *Options) Add(name string, value interface{}) {
    o.Options.Add(name, stringit.Format("{}", value))    
}

type Request struct {
    Symbol  string
    Options  *Options
}

func NewRequest(symbol string, values map[string]interface{}) (r *Request) {
    r = new(Request)
    r.Symbol = symbol
    r.Options = NewOptions()
    r.Options.Add("s", symbol)
    for k, v := range values {
        r.Options.Add(k, v)    
    }
    return
}

// Execute the request. Save the csv with the same name as the symbol we are looking up
func (r *Request) Execute() {
    url := stringit.Format("{}?{}", config.BaseUrl, r.Options.Values.Encode())
    filename := stringit.Format("{}.csv", r.symbol)
    wget.Wget(url, filename)
    return
}
