## linux服务器负载问题排查思路以及常用指令总结

服务器的性能主要看四大件：cpu、内存、磁盘、网络。下面是常用的排查思路。
- [CPU和内存问题](#CPU%e5%92%8c%e5%86%85%e5%ad%98%e9%97%ae%e9%a2%98)
   - [top命令](#top%e5%91%bd%e4%bb%a4)
   - [vmstat命令](#vmstat%e5%91%bd%e4%bb%a4)
   - [free命令](#free%e5%91%bd%e4%bb%a4)
- [磁盘问题](#%e7%a3%81%e7%9b%98%e9%97%ae%e9%a2%98)
   - [iostat命令](#iostat%e5%91%bd%e4%bb%a4)
   - [iotop命令](#iotop%e5%91%bd%e4%bb%a4)
   - [du和df命令](#du%e5%92%8cdf%e5%91%bd%e4%bb%a4)
- [网络问题](#%e7%bd%91%e7%bb%9c%e9%97%ae%e9%a2%98)
   - [netstat命令](#netstat%e5%91%bd%e4%bb%a4)
   - [nload命令](#nload%e5%91%bd%e4%bb%a4)
   - [nethogs命令](#nethogs%e5%91%bd%e4%bb%a4)
   - [tcpdump命令](#tcpdump%e5%91%bd%e4%bb%a4)

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


|           参数                  |                含义                |
| :----------------------------: | :--------------------------------: |
| 17:24:42                       |系统时间                             |
| up  9:10                       |系统运行的时间                        |
| user                           |系统当前登录的用户数                   |
| load average: 1.50, 1.50, 1.56 |过去一分钟五分钟十五分钟系统负载         |
| ----                           |----                                |
| Tasks                          |任务进程数                            |
| total                          |总进程数                              |
| running                        |正在运行的进程数                       |
| sleeping                       |休眠状态进程数                         |
| stopped                        |停止进程数                            |
| ----                           |----                                 |
| %us                            |用户进程消耗的CPU时间百分比               |
| %sy                            |系统进程消耗的CPU时间百分比               |
| %ni	                         |运行正常进程消耗的CPU时间百分比            |
| %id                            |CPU空闲状态的时间百分比                  |
| %wa                            |I/O 等待所占 CPU 时间百分比              |
| %hi	                         |硬中断（Hardware IRQ）占用CPU的百分比     |
| %si	                         |软中断（Software Interrupts）占用CPU的百分比|
| %st	                         |在内存紧张环境下，pagein 强制对不同的页面进行的 steal 操作|
| ----                           |----                                 |
| Mem total                      |物理内存总量                            |
| Mem used                       |使用中的内存总量                         |
| Mem free                       |空闲内存总量                            |
| Mem buffers                    |缓存的内存量                            |
| ----                           |----                                 |
| PID                            |进程id                                |
| USER                           |进程所有者                             |
| PR                             |进程优先级                             |
| NI                             |nice值。负值表示高优先级，正值表示低优先级   |
| VIRT                           |进程使用的虚拟内存总量，单位kb。VIRT=SWAP+RES|
| RES                            |进程使用的、未被换出的物理内存大小，单位kb。RES=CODE+DATA|
| SHR                            |共享内存大小，单位kb                    |
| S                              |进程状态。D=不可中断的睡眠状态 R=运行 S=睡眠 T=跟踪/停止 Z=僵尸进程|
| %CPU                           |次更新到现在的CPU时间占用百分比          |
| %MEM                           |进程使用的物理内存百分比                 |
| TIME+                          |进程使用的CPU时间总计，单位1/100秒       |
| COMMAND                        |进程名称（命令名/命令行）               |

下面是几个常用的交互操作：

- 展开CPU  
在top交互界面直接按键盘的数字1。  
这里还是要强调一下，%cpu的值是跟内核数成正比的，如8核cpu的%cpu最大可以800%。    
- 显示线程    
在top交互界面按ctrl+h显示线程，再按一次关闭显示。  
- 增加或删除显示列  
在top交互界面按h进入，输入想显示的列的首字母如n，退出直接回车。  
- 排序  
Cpu ： 在top交互界面按shift+p。  
Mem ：在top交互界面按shift+m。  
Time ：在top交互界面按shift+t。  
- 显示程序名  
在top交互界面按c。  
- 监控线程下的进程  
在命令行输入top -H -p pid，其中pid为进程id，进入界面后显示的PID为线程ID；或者使用命
令top -H -p pid进入界面之后在按shift+h来显示线程。     


下面是几个重点观察的指标  
**%Cpu(s): 5.1 us, 3.4 sy, 0.0 wa**  
>这里可以非常直观的看到当前cpu的负载情况，us用户cpu占用时间，sy是系统调用cpu占用时间，wa是
>cpu等待io的时间，前面两个比较直观，但是第三个其实也很重要，如果wa很高，那么你就该重点关注下磁
>盘的负载了，尤其是像mysql这种服务器。

**load average: 0.08, 0.26, 0.19**  
>cpu任务队列的负载，这个队列包括正在运行的任务和等待运行的任务，三个数字分别是1分钟、5分钟和15分
>钟的平均值。这个和cpu占用率一般是正相关的，反应的是用户代码，如果超过了内核数，表示系统已经过载。
>也就是说如果你是8核，那么这个数字小于等于8的负载都是没问题的，我看网上的建议一般这个值不要超过
>ncpu*2-2为好。

**KiB Mem : 985856 total, 81736 free, 646360 used, 257760 buff/cache**
>内存占用情况，total总内存，free空余内存， used已经分配内存，buff/cache块设备和缓冲区占用的
>内存，因为Linux的内存分配，如果有剩余内存，他就会将内存用于cache，这样可以较少磁盘的读写提高
>效率，如果有应用申请内存，buff/cache这部分内存也是可用的，所以正真的剩余内存应该是
>free+buff/cache

### vmstat命令  

vmstat 可以对操作系统的内存信息、进程状态、CPU 活动、磁盘等信息进行监控，不足之处是无法对某个进程
进行深入分析。
````
$  vmstat 2 3 -S M
procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 0  0      0   8234    747   3930    0    0   194   102  353 1307  5  1 93  0  0
 0  0      0   8254    747   3910    0    0     0     0 1165 4908  1  1 98  0  0
 1  0      0   8255    747   3910    0    0     0     0 1053 3867  1  1 98  0  0

````
2表示每2秒取样一次，3表示取数3次，-S表示单位，可选有 k 、K 、m 、M。

- procs  
R表示等待运行和等待CPU时间片的进程数，这个值如果长期大于系统CPU个数，说明CPU不足，需要增加CPU
B表示在等待资源的进程数，比如正在等待I/O或者交换内存等。  
- memory  
swpd 表示切换到内存交换区的内存大小（单位KB）,通俗讲就是虚拟内存的大小。如果swap值不为0或者比较
大，只要si,so的值长期为0，这种情况一般属于正常的情况  
free 表示当前空闲的物理内存 （单位KB）  
buff 列表示buffers cached内存大小，也就是缓存区大小，一般对块设备的读写才需要缓冲。  
cache 列表示page cached的内存大小，也就是缓存大小，一般作为文件系统进行缓冲，频繁访问的文件
都会被缓存，如果cache值非常大说明缓存文件较多，如果此时io中的bi比较小，说明文件系统效率比较好。  
- swap  
si 列表示由磁盘调入内存，也就是内存进入内存交换区的内存大小。  
so 列表示由内存进入磁盘，也就是内存交换区进入内存的内存大小。  
一般情况，si、so的值都为0，如果si、so的值长期不为0，则说明系统内存不足，则需要增加系统内存。  
- io   
bi 列表表示由块设备读入数据的总量，既读磁盘，单位kb/s。  
bo 列表示写到快设备数据的总量，既写磁盘，单位kb/s。  
如果 bi+bo 值过大，且 wa 值较大，则表示系统磁盘 IO 瓶颈。  
- system  
in 列表示某一时刻时间间隔内观测到的每秒设备中断数。  
cs 列表示每秒产生的上下文切换次数。  
这两个值越大，则由内核消耗的CPU就越多。  
- cpu  
us 列表示用户进程消耗的CPU时间百分比，us值越高，说明用户进程消耗cpu时间越多，如果长期大于50%，
则需要考虑优化程序或算法。  
sy 列表示系统内核进程消耗的cpu时间百分比，一般来说us+sy应该小于80%，如果大于80%，说明可能处
cpu瓶颈。  
id 列表示cpu处在空闲状态的百分比。  
wa 列表示IP等待说占的cpu时间百分比，wa值越高，说明I/O等待越严重，根据经验wa的参考值为20%，如果
超过20%，说明I/O等待严重，引起I/O等待的原因可能是磁盘大量随机读写造成的，也可能是磁盘造成的。

### free命令

free命令是监控linux内存使用最常用的命令，参数[-m]表示m为单位查看内存的使用情况（默认kb）  
````
$ free -m
              total        used        free      shared  buff/cache   available
Mem:          15912        5389        3914        1042        6608        9163
Swap:          4095           0        4095
````
- Mem:物理内存大小  
- total:总计物理内存的大小  
- used:已使用的大小  
- free:可用的大小  
- shared:多个进程共享的内存的总额  
- buffers:缓冲区内存总量  
- cached:交换缓存区内存总量   
- Swap:交换缓冲区内存总量  

## 磁盘问题

磁盘问题在mysql服务器中非常常见，很多时候mysql服务器的CPU不高但是却出现慢查询日志飙升，就是因为
磁盘出现了瓶颈。还有mysql的备份策略，如果没有监控磁盘空间，可能出现磁盘满了服务不可用的现象。  

### iostat命令
deepin上面的安装
````
apt-get install sysstat
````
常用参数： -k 用kb为单位 -d 监控磁盘 -x显示详情 num count 每个几秒刷新 显示次数  

使用iostat -kdx 2 10 跑一下
````
$ iostat -kdx 2 10
Linux 4.15.0-30deepin-generic (liz-PC) 	2020年01月14日 	_x86_64_	(4 CPU)

Device:         rrqm/s   wrqm/s     r/s     w/s    rkB/s    wkB/s avgrq-sz avgqu-sz   await r_await w_await  svctm  %util
nvme0n1           0.00    10.62   45.76    8.17   702.73   348.46    38.99     0.10    2.04    2.08    1.82   0.04   0.23

Device:         rrqm/s   wrqm/s     r/s     w/s    rkB/s    wkB/s avgrq-sz avgqu-sz   await r_await w_await  svctm  %util
nvme0n1           0.00     0.00    0.00    4.00     0.00   332.00   166.00     0.00    0.00    0.00    0.00   0.00   0.00

````

- rkB/s和wkB/s  
分别对应读写速度

- avgqu-sz  
读写队列的平均请求长度，可以类比top命令的load average

- await r_await w_await  
io请求的平均时间（毫秒），分别是读写，读和写三个平均值。这个时间都包括在队列中等待的时间和实际
处理读写请求的时间，还有svctm这个参数，他说的是实际处理读写请求的时间，照理来讲wawait肯定是
大于svctm的，但是我在线上看到有wawait小于svctm的情况，不知道是什么原因。我看iostat的man手动
中说svctm已经废弃，所以一般我看的是这三个。

- %util  
这个参数直观的看磁盘的负载情况，我首先看的就是这个参数。和top的wa命令有关联。

### iotop命令

主要是用于直观的看那些进程占用io较高，是否有异常的进程。
````
Total DISK READ :       0.00 B/s | Total DISK WRITE :       0.00 B/s
Actual DISK READ:       0.00 B/s | Actual DISK WRITE:       0.00 B/s
TID  PRIO  USER     DISK READ  DISK WRITE  SWAPIN     IO>    COMMAND                                                                                                                                                    
1 be/4 root        0.00 B/s    0.00 B/s  0.00 %  0.00 % init auto noprompt
2 be/4 root        0.00 B/s    0.00 B/s  0.00 %  0.00 % [kthreadd]
4 be/0 root        0.00 B/s    0.00 B/s  0.00 %  0.00 % [kworker/0:0H]
6 be/0 root        0.00 B/s    0.00 B/s  0.00 %  0.00 % [mm_percpu_wq]
7 be/4 root        0.00 B/s    0.00 B/s  0.00 %  0.00 % [ksoftirqd/0]
8 be/4 root        0.00 B/s    0.00 B/s  0.00 %  0.00 % [rcu_sched]
9 be/4 root        0.00 B/s    0.00 B/s  0.00 %  0.00 % [rcu_bh]
````
iotop具有与top相似的UI，其中包括PID、用户、I/O、进程等相关信息。如果你想知道每个进程是如何使用IO
的就比较麻烦，使用iotop命令可以很方便的查看。  

####　输出的参数的意义： 

第一行：Read和Write速率总计  

第二：实际的Read和Write速率  

第三行：参数如下：  

线程ID（按p切换为进程ID）  
- 优先级
- 用户
- 磁盘读速率
- 磁盘写速率
- swap交换百分比
- IO等待所占的百分比
- 线程/进程命令

#### 常用的参数
````
-o, --only只显示正在产生I/O的进程或线程。除了传参，可以在运行过程中按o生效。
-b, --batch非交互模式，一般用来记录日志。
-n NUM, --iter=NUM设置监测的次数，默认无限。在非交互模式下很有用。
-d SEC, --delay=SEC设置每次监测的间隔，默认1秒，接受非整形数据例如1.1。
-p PID, --pid=PID指定监测的进程/线程。
-u USER, --user=USER指定监测某个用户产生的I/O。
-P, --processes仅显示进程，默认iotop显示所有线程。
-a, --accumulated显示累积的I/O，而不是带宽。
-k, --kilobytes使用kB单位，而不是对人友好的单位。在非交互模式下，脚本编程有用。
-t, --time 加上时间戳，非交互非模式。
-q, --quiet 禁止头几行，非交互模式。有三种指定方式。
-q 只在第一次监测时显示列名
-qq 永远不显示列名。
-qqq 永远不显示I/O汇总。
````

#### 交互按钮　　
和top命令类似，iotop也支持以下几个交互按键。　　
````
left和right方向键：改变排序。
r：反向排序。
o：切换至选项–only。
p：切换至–processes选项。
a：切换至–accumulated选项。
q：退出。
i：改变线程的优先级。
````
### du和df命令

主要通过这两个命令看系统的磁盘占用率和文件夹大小，有时候日志文件不清理会导致磁盘用满等情况。  

- du  
````
用法： df -h 查看磁盘占用情况
````
- df  
````
用法： du -sh 查看当前目录 容量
````

**典型问题**  
mysql负载大时，很多时候磁盘先到了瓶颈，大量个请求超时，cpu负载却不高，如果mysql服务器异常，建
议重点看下磁盘。

### 网络问题  

在线上服务器，大部分服务器都市只能内网访问，放在公网的的服务器也就那几台nginx和ftp的，另外公网
的那些服务器都有流量监控，所以网络问题一般并不大。

### netstat命令  
netstat命令用于显示本机网络连接、运行端口、路由表信息，接口状态 (Interface Statistics)，
masquerade 连接，多播成员 (Multicast Memberships) 等等。

下面是netstat的命令详解
https://www.cnblogs.com/ricklz/p/11796319.html  

### nload命令

**在Debian上面安装nload**  
````
sudo apt-get install nload
````
用于监控整体的带宽，即时监看网路状态和各IP所使用的频宽。  
**nload默认分为上下两块：**  
上半部分是：Incoming也就是进入网卡的流量，  
下半部分是：Outgoing，也就是从这块网卡出去的流量，  
每部分都有当前流量（Curr），  

平均流量（Avg），  

最小流量（Min），  

最大流量（Max），  

总和流量（Ttl）这几个部分，看起来还是蛮直观的。  

nload默认的是eth0网卡，如果你想监测eth1网卡的流量  

#nload eth1  



-a：这个好像是全部数据的刷新时间周期，单位是秒，默认是300.  

-i：进入网卡的流量图的显示比例最大值设置，默认10240 kBit/s.  

-m：不显示流量图，只显示统计数据。  

-o：出去网卡的流量图的显示比例最大值设置，默认10240 kBit/s.  

-t：显示数据的刷新时间间隔，单位是毫秒，默认500。  

-u：设置右边Curr、Avg、Min、Max的数据单位，默认是自动变的.注意大小写单位不同！  

h|b|k|m|g h: auto, b: Bit/s, k: kBit/s, m: MBit/s etc.  

H|B|K|M|G H: auto, B: Byte/s, K: kByte/s, M: MByte/s etc.  

-U：设置右边Ttl的数据单位，默认是自动变的.注意大小写单位不同（与-u相同）！  
Devices：自定义监控的网卡，默认是全部监控的，使用左右键切换。  
如只监控eth0命令：# nload eth0  
使用 $ nload eth0 ，可以查看第一网卡的流量情况，显示的是实时的流量图， $ nload -m 可以同时查看多个网卡的流量情况。  

### nethogs命令  

用于监控进程的带宽使用情况  

### tcpdump命令

## 参考  
【引用】https://cloud.tencent.com/developer/article/1378739  
【Linux中的nload命令】https://www.jianshu.com/p/08b60e90a909  
【linux 服务器性能监控（一）】https://www.jianshu.com/p/9e571b2b4971  

  