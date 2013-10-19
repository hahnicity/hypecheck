package main

import (
    "flag"
    "github.com/hahnicity/hypecheck"
    "github.com/hahnicity/hypecheck/data"
)

var (
    days         int
    maxRequests  int
    requestDelay int
    threshold    float64
)

func parseArgs() {
    flag.IntVar(
        &days,
        "days",
        10,
        "The number of trading days after a swing has occurred to look for a return",
    )
    flag.IntVar(
        &maxRequests,
        "maxRequests",
        100,
        "The maximum number of requests we can make to yahoo at once",
    )
    flag.IntVar(
        &requestDelay,
        "d",
        100,
        "The time to wait between for a new request to be executed after one is made",
    )
    flag.Float64Var(
        &threshold,
        "threshold",
        .05,
        "The threshold at which we determine a price swing should be analyzed",
    )
    flag.Parse()
}

func main() {
    parseArgs()
    a := hypecheck.NewAnalyzer(days, threshold)
    r := hypecheck.NewRequester(a, maxRequests, requestDelay) 
    go hypecheck.NewBalancer(maxRequests).Balance(r.Work)
    r.MakeRequests(data.TESTLIST)
}
