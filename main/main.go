package main

import (
    "flag"
    "github.com/hahnicity/go-stringit"
    "github.com/hahnicity/hypecheck"
    "github.com/hahnicity/hypecheck/data"
    "strconv"
    "strings"
    "time"
)

var (
    days         int
    endDate      string
    maxRequests  int
    requestDelay int
    startDate    string
    threshold    float64
)

func parseArgs() {
    flag.StringVar(
        &startDate,
        "s",
        "2013-01-01",
        "The date to begin searching",
    )
    flag.StringVar(
        &endDate,
        "e",
        stringit.Format("{}-{}-{}", time.Now().Year(), int(time.Now().Month()), time.Now().Day()),
        "The date to finish searching",
    )
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

func addURLOptions() map[string]interface{} {
    // XXX Maybe just use time.Parse instead? too tired to do anything sensible
    endMonth, _ := strconv.ParseInt(strings.Split(endDate, "-")[1], 10, 8)
    startMonth, _ := strconv.ParseInt(strings.Split(startDate, "-")[1], 10, 8)
    var values = map[string]interface{}{
        "a": startMonth - 1, // Yahoo makes us subtract 1 from the month number so March (3) becomes 2
        "b": strings.Split(startDate, "-")[2],
        "c": strings.Split(startDate, "-")[0], 
        "d": endMonth - 1,
        "e": strings.Split(endDate, "-")[2],
        "f": strings.Split(endDate, "-")[0],
    }
    return values
}

func main() {
    parseArgs()
    a := hypecheck.NewAnalyzer(days, threshold)
    r := hypecheck.NewRequester(a, maxRequests, requestDelay)
    go hypecheck.NewBalancer(600).Balance(r.Work) // XXX There is a bug with push
    r.MakeRequests(data.SP500, addURLOptions())
}
