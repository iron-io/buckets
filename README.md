
## What is this?

Buckets is basically way to count events over time. 

## Usage

```go
bm := buckets.New()
bm.AddSet("set1")
bm.AddSet("set2")

# Then whenever an event happens that you want to record:
bm.Inc(&DefaultOccurrence{"set2", time.Now()})
// Occurence can be custom, just needs to adhere to [Occurrence interface](https://github.com/iron-io/buckets/blob/master/buckets.go#L155) 
```

Then you can access the buckets directly:

```go
set1 := bm.Get("set1")

// This will be an array of int64's:
set1Buckets := set1.Buckets
// You can use this to generate graphs/charts or whatever you need

// Or ask for total
fmt.Println("Total:", set1.Total()
```

## Reporters

Reporters can post the data to log or data collection service:

```
bm.AddReporter(NewStdoutReporter())
```

You can make your own custom reporters easily too, they just need to follow the [Reporter interface](https://github.com/iron-io/buckets/blob/master/reporters.go#L15). 

## Testing

Make a `test_config.json` file with:

```json
{
    "reporters": [
        {
            "service": "stdout"
        }
    ]
}
```

