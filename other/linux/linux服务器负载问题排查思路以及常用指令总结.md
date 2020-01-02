## linux服务器负载问题排查思路以及常用指令总结

服务器的性能主要看四大件：cpu、内存、磁盘、网络。下面是常用的排查思路。

## CPU和内存问题

### top命令
top命令能够实时监控系统的运行状态，并且按照CPU、内存和执行时间进行排序，
同时top命令还可以通过交互式命令进行设定显示，通过top命令可以查看即时活
跃的进程。  
````
top - 17:24:42 up  9:10,  1 user,  load average: 1.50, 1.50, 1.56
Tasks: 274 total,   1 running, 219 sleeping,   0 stopped,   1 zombie
%Cpu(s):  3.2 us,  3.1 sy,  0.0 ni, 93.4 id,  0.1 wa,  0.0 hi,  0.2 si,  0.0 st
KiB Mem : 16294156 total,  2769376 free,  5753276 used,  7771504 buff/cache
KiB Swap:  4194300 total,  4194300 free,        0 used.  9227048 avail Mem 

  PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND                                         
  623 root     -51   0       0      0      0 D  11.6  0.0  25:27.38 irq/136-SYNA308                                 
 1130 root      20   0  682972 238292 178304 S   5.0  1.5  34:02.22 Xorg                                            
 4163 liz       20   0 3451648 155092  94756 S   4.0  1.0  26:33.41 kwin_x11                                        
 8084 liz       20   0 3402536 320632  53192 S   3.7  2.0  32:08.25 WeChat.exe                                      
 8090 liz       20   0    9992   7192   1828 S   3.0  0.0  17:08.10 wineserver.real                                 
18503 liz       20   0  648956  71024  35452 S   3.0  0.4   0:24.04 deepin-terminal                                 
 8912 liz       20   0 7502660 1.452g 146036 S   2.7  9.3  36:00.50 java                                            
 8074 liz       20   0 1241784 471808 273892 S   2.0  2.9  37:12.16 dingtalk                                        
18542 liz       20   0  511728  38300  10644 S   0.7  0.2   1:19.92 docker-compose                                  
23773 liz       20   0 1060704 286388 110568 S   0.7  1.8   1:25.11 chrome                                          
25793 liz       20   0   46812   3752   3024 R   0.7  0.0   0:02.45 top                                             
 3171 liz       20   0 2101792  23552  16524 S   0.3  0.1   0:50.57 startdde                                        
 4200 liz       20   0 2244084  52888  32340 S   0.3  0.3   1:01.10 dde-session-dae                                 
18441 liz       20   0  750956 108936  78260 S   0.3  0.7   2:17.41 dde-file-manage                                 
18813 deepin-+  20   0   25252   2916   1812 S   0.3  0.0   0:27.21 redis-server                                    
19712 liz       20   0 3639568 332556  41980 S   0.3  2.0   1:01.00 Navicat.exe                                     
23525 liz       20   0  929804 270044 153808 S   0.3  1.7   0:44.78 chrome                                          
25179 liz       20   0 1626084 217176  86820 S   0.3  1.3   5:04.57 _Postman                                        
25391 root      20   0       0      0      0 I   0.3  0.0   0:00.83 kworker/u8:2                                    
    1 root      20   0  204880   7132   5184 S   0.0  0.0   0:04.24 systemd                                         
    2 root      20   0       0      0      0 S   0.0  0.0   0:00.02 kthreadd                                        
    4 root       0 -20       0      0      0 I   0.0  0.0   0:00.00 kworker/0:0H                                    
    6 root       0 -20       0      0      0 I   0.0  0.0   0:00.00 mm_percpu_wq                                    
    7 root      20   0       0      0      0 S   0.0  0.0   0:02.24 ksoftirqd/0                                     
    8 root      20   0       0      0      0 I   0.0  0.0   0:32.35 rcu_sched                                       
    9 root      20   0       0      0      0 I   0.0  0.0   0:00.00 rcu_bh                                          
   10 root      rt   0       0      0      0 S   0.0  0.0   0:00.03 migration/0                                     
   11 root      rt   0       0      0      0 S   0.0  0.0   0:00.06 watchdog/0                                      
   12 root      20   0       0      0      0 S   0.0  0.0   0:00.00 cpuhp/0                                         
````
那么我们来一一分析里面每个值的含义  
 -  17:24:42        系统运行的时间





常用参数：-H打印具体的线程，-p打印某个进程 进入后 按数字1 可以切换cpu的
图像看有几个核  



























## 参考  
【引用】https://cloud.tencent.com/developer/article/1378739  

