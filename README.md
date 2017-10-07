# logrepeat
read access log then repeat requests to another target host:port

## Usage

- Download your AWS ALB(Application Load Balancer) log from S3.

- Run logrepeat with your alb log file.

```
$ logrepeat -h newtarget -p 8080 -f youralblogfile.log
```

- You will see and check repeat plan 

```
--- Requests Source ---
REQUEST TIME	: 2017-09-09 09:55:00JST ~ 2017-09-09 10:55:00JST
REQUESTS    	: 98111      reqs
IGNORED    	  : 0          reqs
NON SUPPORTED	: 0          reqs
PARSE ERROR	  : 0          reqs
DryRun     	  : true

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

- If ok, enter "start" at prompt and press Enter. then repeat requests will start.

```
Enter [start] and press Enter key>start
Start at: 2017-09-21 22:39:49JST
2017-09-22T23:11:07JST - 2017-09-22T23:11:07JST  /2xx:   1 /3xx:   0 /4xx:   0 /5xx:   0 /Oth:   0 /Err:   0 Avg 113 msec  (0.0%Done)
2017-09-22T23:11:09JST - 2017-09-22T23:11:10JST  /2xx:  43 /3xx:   0 /4xx:   0 /5xx:   0 /Oth:   0 /Err:   0 Avg 202 msec  (1.3%Done)
2017-09-22T23:11:10JST - 2017-09-22T23:11:14JST  /2xx: 127 /3xx:   0 /4xx:   0 /5xx:   0 /Oth:   0 /Err:   0 Avg 235 msec  (4.9%Done)
2017-09-22T23:11:15JST - 2017-09-22T23:11:18JST  /2xx:  99 /3xx:   0 /4xx:   0 /5xx:   0 /Oth:   0 /Err:   0 Avg 216 msec  (7.8%Done)
2017-09-22T23:11:18JST - 2017-09-22T23:11:22JST  /2xx: 113 /3xx:   0 /4xx:   0 /5xx:   0 /Oth:   0 /Err:   0 Avg 210 msec  (11.0%Done)
...
...
```

## Advanced Usage

Download ALB log with selecting time range, and repeat them to another target by one liner with aloget(https://github.com/yomon8/aloget).

```
$ logrepeat -h newtarget -p 8080 -f <(aloget -b <S3Bucket> -p <ALBAccessLogPrefix> -stdout -s yyyy-MM-ddTHH:mm:ss -e yyyy-MM-ddTHH:mm:ss)
```



