# Procss CPU stat

Linux tool show show thread CPU usage under a process.

## The python version

Usage: `./proc_cpu_stat.py [pid]`

For example:

```
$ ./proc_cpu_stat.py 21283

21283 chromium            56.28 ########################################################
-------------------------------------
21283 chromium             1.97 ##
21284 ThreadPoolServi      0.00 
21286 Chrome_ChildIOT      1.97 ##
21288 GpuMemoryThread      0.00 
21289 Compositor           0.99 #
21290 ThreadPoolSingl      1.97 ##
21291 CompositorTileW      4.94 #####
21292 CompositorTileW      5.92 ######
21293 CompositorTileW      0.00 
21335 MemoryInfra          0.00 
21337 Media                3.95 ####
21653 ThreadPoolForeg      0.99 #
21968 ThreadPoolForeg      0.00 
21977 ThreadPoolForeg      0.00 
21978 Media               10.86 ###########
21979 Media                8.89 #########
21980 Media               11.85 ############
21982 AudioOutputDevi      1.97 ##

```

## The golang version

To get the golang version:

```
go get github.com/trentzhou/proc-cpu-stat
```

Then you'll see the executable `proc-cpu-stat` under $GOPATH/bin.
