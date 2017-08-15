# gofaster
A Micro-Web-Framework Designed to Support Headers

Many of the existing web-frameworks serve as an extension of net/http or utilize fasthttp
both of which seem to **choke** when dealing with a lot of additional headers

* **FastHttp looses x3 it's speed after adding 4 headers to every request**
* **Net/Http looses almost x7 it's speed after adding the same headers to every request**


### Benchmarks:
___

**Last Test Updated:** 2017-08-14

*test environment*

* CPU:      Intel(R) Core(TM) i3-5020U CPU @ 2.20GHz
* Memory:   12G
* Go:       1.8.3
* OS:       Ubuntu 14.04 (Trusty Tahr)

Net/Http **Without** Headers:
```
Running 30s test @ http://localhost:8080
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    18.52ms   45.66ms 773.89ms   91.98%
    Req/Sec     5.84k     2.89k   24.20k    77.60%
  2096219 requests in 30.09s, 257.89MB read
Requests/sec:  69668.78
Transfer/sec:      8.57MB
```
Net/Http **With** Headers:
```
Running 30s test @ http://localhost:8080
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    22.01ms   35.92ms 441.88ms   86.90%
    Req/Sec     1.06k     1.11k    6.93k    81.11%
  308726 requests in 30.10s, 50.05MB read
  Socket errors: connect 58, read 0, write 0, timeout 0
Requests/sec:  10257.50
Transfer/sec:      1.66MB
```
Fasthttp **Without** Headers:
```
Running 30s test @ http://localhost:8080
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.63ms    1.93ms  59.17ms   87.01%
    Req/Sec    13.16k     1.71k   47.80k    74.80%
  4720619 requests in 30.08s, 661.78MB read
Requests/sec: 156943.19
Transfer/sec:     22.00MB
```
Fasthttp **With** Headers:
```
Running 30s test @ http://localhost:8080
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     3.34ms    9.90ms 401.52ms   98.55%
    Req/Sec     3.87k     1.36k   10.45k    73.18%
  1356347 requests in 30.09s, 219.90MB read
Requests/sec:  45072.57
Transfer/sec:      7.31MB
```
GoFaster **Without** Headers:
```
Running 30s test @ http://localhost:8080
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     8.17ms    8.63ms 224.50ms   92.78%
    Req/Sec     3.53k   413.15     8.45k    71.90%
  1267217 requests in 30.10s, 93.06MB read
  Socket errors: connect 0, read 1267015, write 0, timeout 0
Requests/sec:  42100.76
Transfer/sec:      3.09MB
```
GoFaster **With** Headers:
```
Running 30s test @ http://localhost:8080
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     9.08ms    8.57ms 219.27ms   88.48%
    Req/Sec     3.37k   385.27     6.11k    72.85%
  1210974 requests in 30.10s, 196.33MB read
Requests/sec:  40237.18
Transfer/sec:      6.52MB
```
