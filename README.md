# logrepeat
read access log then repeat requests to another target host:port

## Usage

- Fiast download your AWS ALB(Application Load Balancer) log from S3.

- exec logrepeat with your alb log file.

```
$ logrepeat -h newtarget -p 8080 -f youralblogfile.log
```

- you can check repeat plan 

```
--- Requests Source ---
REQUEST TIME	: 2017-09-09 09:55:00JST ~ 2017-09-09 10:55:00JST
REQUEST COUNT	: 98111      reqs
IGNORED COUNT	: 0          reqs
PARSE ERROR	: 0          reqs
DryRun     	: true

--- Repeat Plan ---
REPEAT TIME	: 2017-09-21 22:39:49JST ~ 2017-09-21 23:39:50JST
REPEAT TARGET	: newtarget:8080
Repeat Samples	:
 1: http://newtarget:8080/?test=1
 2: http://newtarget:8080/?test=2
 3: http://newtarget:8080/?test=3
 4: http://newtarget:8080/?test=4
 5: http://newtarget:8080/?test=5
...and more
```

- if ok, enter "start" at prompt and press Enter. then repeat will start.

```
Enter [start] and press Enter key>start
Start at: 2017-09-21 22:39:49JST
2017-09-21T22:39:49JST - 2017-09-21T22:39:49JST  /2xx:   0 /3xx:   0 /4xx:   0 /5xx:   0 /Oth:   1  (0.0%)
2017-09-21T22:39:50JST - 2017-09-21T22:39:52JST  /2xx:   0 /3xx:   0 /4xx:   0 /5xx:   0 /Oth: 116  (0.1%)
2017-09-21T22:39:53JST - 2017-09-21T22:39:56JST  /2xx:   0 /3xx:   0 /4xx:   0 /5xx:   0 /Oth: 129  (0.3%)
2017-09-21T22:39:57JST - 2017-09-21T22:40:00JST  /2xx:   0 /3xx:   0 /4xx:   0 /5xx:   0 /Oth: 131  (0.4%)
2017-09-21T22:40:01JST - 2017-09-21T22:40:04JST  /2xx:   0 /3xx:   0 /4xx:   0 /5xx:   0 /Oth: 127  (0.5%)
```




