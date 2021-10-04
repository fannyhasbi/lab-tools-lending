# Benchmark Result
This is the result of HTTP benchmarking of the staging/production server, see the `benchmark` folder to get a more detail of the pipeline process.

## Register
```
Running 30s test @ https://example.com
  1 threads and 1 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   557.64ms   33.41ms 712.60ms   90.38%
    Req/Sec     1.02      0.24     2.00     94.23%
  52 requests in 30.10s, 6.65KB read
Requests/sec:      1.73
Transfer/sec:     226.32B
```

## Help
```
Running 30s test @ https://example.com
  1 threads and 1 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   567.09ms   50.97ms 840.79ms   90.20%
    Req/Sec     1.12      0.38     2.00     84.31%
  51 requests in 30.08s, 6.52KB read
Requests/sec:      1.70
Transfer/sec:     222.09B
```

## Tool
```
Running 30s test @ https://example.com
  1 threads and 1 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   532.30ms  426.09ms   1.29s     0.00%
    Req/Sec     2.33      0.95     4.00     69.33%
  75 requests in 30.09s, 9.59KB read
Requests/sec:      2.49
Transfer/sec:     326.50B
```

## Borrow
```
Running 30s test @ https://example.com
  1 threads and 1 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   470.77ms  714.99ms   1.98s     0.00%
    Req/Sec     2.48      1.13     5.00     65.33%
  75 requests in 30.01s, 9.59KB read
  Socket errors: connect 0, read 0, write 0, timeout 4
Requests/sec:      2.50
Transfer/sec:     327.38B
```

## Returning
```
Running 30s test @ https://example.com
  1 threads and 1 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   513.21ms  610.86ms   1.67s     0.00%
    Req/Sec     2.41      1.20     5.00     58.57%
  72 requests in 30.02s, 9.21KB read
  Socket errors: connect 0, read 0, write 0, timeout 1
Requests/sec:      2.40
Transfer/sec:     314.22B
```

## Response
```
Running 30s test @ https://example.com
  1 threads and 1 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   634.24ms   62.37ms 743.63ms   53.57%
    Req/Sec     0.87      0.34     1.00     87.10%
  31 requests in 30.02s, 5.07KB read
  Socket errors: connect 0, read 0, write 0, timeout 3
  Non-2xx or 3xx responses: 11
Requests/sec:      1.03
Transfer/sec:     173.04B
```

## Manage
```
Running 30s test @ https://example.com
  1 threads and 1 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   628.30ms   59.09ms 770.16ms   57.14%
    Req/Sec     0.87      0.34     1.00     87.10%
  31 requests in 30.00s, 5.07KB read
  Socket errors: connect 0, read 0, write 0, timeout 3
  Non-2xx or 3xx responses: 11
Requests/sec:      1.03
Transfer/sec:     173.10B
```

## Report
```
Running 30s test @ https://example.com
  1 threads and 1 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   621.20ms   52.27ms 757.58ms   71.43%
    Req/Sec     0.87      0.34     1.00     87.10%
  31 requests in 30.02s, 5.07KB read
  Socket errors: connect 0, read 0, write 0, timeout 3
  Non-2xx or 3xx responses: 11
Requests/sec:      1.03
Transfer/sec:     173.03B
```