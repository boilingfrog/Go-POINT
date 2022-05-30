<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 如何应对并发访问](#redis-%E5%A6%82%E4%BD%95%E5%BA%94%E5%AF%B9%E5%B9%B6%E5%8F%91%E8%AE%BF%E9%97%AE)
  - [使用 LUA 脚本](#%E4%BD%BF%E7%94%A8-lua-%E8%84%9A%E6%9C%AC)
  - [Redis 中如何使用 LUA 脚本](#redis-%E4%B8%AD%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8-lua-%E8%84%9A%E6%9C%AC)
    - [EVAL](#eval)
    - [EVALSHA](#evalsha)
    - [SCRIPT 命令](#script-%E5%91%BD%E4%BB%A4)
      - [SCRIPT LOAD](#script-load)
      - [SCRIPT EXISTS](#script-exists)
      - [SCRIPT FLUSH](#script-flush)
      - [SCRIPT KILL](#script-kill)
    - [SCRIPT DEBUG](#script-debug)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Redis 如何应对并发访问

上个文章中，我们分析了Redis 中命令的执行是单线程的，虽然 Redis6.0 版本之后，引入了 I/O 多线程，但是对于 Redis 命令的还是单线程去执行的。所以如果业务中，我们只用 Redis 中的单命令去处理业务的话，命令的原子性是可以得到保障的。  

但是很多业务场景中，需要多个命令组合的使用，例如前面介绍的 `读取-修改-写回` 场景，这时候就不能保证组合命令的原子性了。所以这时候 LUA 就登场了。

### 使用 LUA 脚本

Redis 在 2.6 版本推出了 lua 脚本功能。  

引入 lua 脚本的优点：  

1、减少网络开销。可以将多个请求通过脚本的形式一次发送，减少网络时延。  

2、原子操作。Redis会将整个脚本作为一个整体执行，中间不会被其他请求插入。因此在脚本运行过程中无需担心会出现竞态条件，无需使用事务。  

3、复用。客户端发送的脚本会永久存在redis中，这样其他客户端可以复用这一脚本，而不需要使用代码完成相同的逻辑。  

关于 lua 的语法和 lua 是一门什么样的语言，可以自行 google。  

### Redis 中如何使用 LUA 脚本

redis 中支持 LUA 脚本的几个命令  

redis 自 2.6.0 加入了 lua 脚本相关的命令，在 3.2.0 加入了 lua 脚本的调试功能和命令 `SCRIPT DEBUG`。这里对命令做下简单的介绍。

EVAL：使用改命令来直接执行指定的Lua脚本；  

SCRIPT LOAD：将脚本 script 添加到脚本缓存中，以达到重复使用，避免多次加载浪费带宽，该命令不会执行脚本。仅加载脚本缓存中；  

EVALSHA：执行由 `SCRIPT LOAD` 加载到缓存的命令；  

SCRIPT EXISTS：以 SHA1 标识为参数,检查脚本是否存在脚本缓存里面

SCRIPT FLUSH：清空 Lua 脚本缓存，这里是清理掉所有的脚本缓存；  

SCRIPT KILL：杀死当前正在运行的 Lua 脚本，当且仅当这个脚本没有执行过任何写操作时，这个命令才生效；  

SCRIPT DEBUG：设置调试模式，可设置同步、异步、关闭，同步会阻塞所有请求。

#### EVAL

通过这个命令来直接执行执行的 LUA 脚本，也是 Redis 中执行 LUA 脚本最常用的命令。  

```
EVAL script numkeys key [key ...] arg [arg ...]
```

来看下具体的参数  

- script: 需要执行的 LUA 脚本；  

- numkeys: 指定的 Lua 脚本需要处理键的数量，其实就是 key 数组的长度；  

- key: 传递给 Lua 脚本零到多个键，空格隔开，在 Lua 脚本中通过 `KEYS[INDEX]` 来获取对应的值，其中`1 <= INDEX <= numkeys`；  

- arg: 自定义的参数,在 Lua 脚本中通过 `ARGV[INDEX]` 来获取对应的值，其中 INDEX 的值从1开始。   

看了这些还是好迷糊，举个栗子  

```go
127.0.0.1:6379> eval "return {KEYS[1],KEYS[2],ARGV[1],ARGV[2],ARGV[3]}" 2 key1 key2 arg1 arg2 arg3
1) "key1"
2) "key2"
3) "arg1"
4) "arg2"
5) "arg3"
```

可以看到上面指定了 numkeys 的长度是2，然后后面 key 中放了两个键值 key1 和 key2，通过 `KEYS[1],KEYS[2]` 就能获取到传入的两个键值对。`arg1 arg2 arg3` 即为传入的自定义参数，通过 `ARGV[index]` 就能获取到对应的参数。   

一般情况下，会将 lua 放在一个单独的 Lua 文件中，然后去执行这个 Lua 脚本。  

<img src="/img/redis/redis-lua.jpg"  alt="redis" />

执行语法 `--eval script  key1 key2 , arg1 age2`  

举个栗子  

```
# cat test.lua
return {KEYS[1],KEYS[2],ARGV[1],ARGV[2],ARGV[3]}

# redis-cli --eval ./test.lua  key1 key2 ,  arg1 arg2 arg3
1) "key1"
2) "key2"
3) "arg1"
4) "arg2"
5) "arg3"
```

需要注意的是，使用文件去执行，key 和 value 用一个逗号隔开，并且也不需要指定 numkeys。  

Lua 脚本中一般会使用下面两个函数来调用 Redis 命令  

```
redis.call()
redis.pcall()
```

redis.call() 与 redis.pcall() 很类似, 他们唯一的区别是当redis命令执行结果返回错误时， redis.call() 将返回给调用者一个错误，而 redis.pcall() 会将捕获的错误以 Lua 表的形式返回。  

```
127.0.0.1:6379> EVAL "return redis.call('SET','test')" 0
(error) ERR Error running script (call to f_77810fca9b2b8e2d8a68f8a90cf8fbf14592cf54): @user_script:1: @user_script: 1: Wrong number of args calling Redis command From Lua script
127.0.0.1:6379> EVAL "return redis.pcall('SET','test')" 0
(error) @user_script: 1: Wrong number of args calling Redis command From Lua script
```

同样需要注意的是，脚本里使用的所有键都应该由 KEYS 数组来传递，就像这样：  

```
127.0.0.1:6379>  eval "return redis.call('set',KEYS[1],'bar')" 1 foo
OK
```

下面这种就是不推荐的  

```
127.0.0.1:6379> eval "return redis.call('set','foo','bar')" 0
OK
```

原因有下面两个  

1、Redis 中所有的命令，在执行之前都会被分析，来确定会对那些键值对进行操作，对于 EVAL 命令来说，必须使用正确的形式来传递键，才能确保分析工作正确地执行；  

2、使用正确的形式来传递键还有很多其他好处，它的一个特别重要的用途就是确保 Redis 集群可以将你的请求发送到正确的集群节点。  

#### EVALSHA

用来执行被 `SCRIPT LOAD` 加载到缓存的命令，具体看下文的 `SCRIPT LOAD` 命令介绍。  

#### SCRIPT 命令

Redis 提供了以下几个 SCRIPT 命令，用于对脚本子系统(scripting subsystem)进行控制。   

##### SCRIPT LOAD

将脚本 script 添加到脚本缓存中，以达到重复使用，避免多次加载浪费带宽，该命令不会执行脚本。仅加载脚本缓存中。    

在脚本被加入到缓存之后，会返回一个通过SHA校验返回唯一字符串标识，使用 EVALSHA 命令来执行缓存后的脚本。  

```
127.0.0.1:6379> SCRIPT LOAD "return {KEYS[1]}"
"8e5266f6a4373624739bd44187744618bc810de3"
127.0.0.1:6379> EVALSHA 8e5266f6a4373624739bd44187744618bc810de3 1 hello
1) "hello"
```

##### SCRIPT EXISTS

以 SHA1 标识为参数,检查脚本是否存在脚本缓存里面。

这个命令可以接受一个或者多个脚本 SHA1 信息，返回一个1或者0的列表。    

```
127.0.0.1:6379> SCRIPT EXISTS 8e5266f6a4373624739bd44187744618bc810de3 2323211
1) (integer) 1
2) (integer) 0
```
 
1 表示存在，0 表示不存在  

##### SCRIPT FLUSH

清空 Lua 脚本缓存 `Flush the Lua scripts cache`，这个是清掉所有的脚本缓存。要慎重使用。   

##### SCRIPT KILL

杀死当前正在运行的 Lua 脚本，当且仅当这个脚本没有执行过任何写操作时，这个命令才生效。  

这个命令主要用于终止运行时间过长的脚本，比如一个因为 BUG 而发生无限 loop 的脚本。  

```
# 没有脚本在执行时
127.0.0.1:6379> SCRIPT KILL
(error) ERR No scripts in execution right now.

# 成功杀死脚本时
127.0.0.1:6379> SCRIPT KILL
OK
(1.10s)

# 尝试杀死一个已经执行过写操作的脚本，失败
127.0.0.1:6379> SCRIPT KILL
(error) ERR Sorry the script already executed write commands against the dataset. You can either wait the script termination or kill the server in an hard way using the SHUTDOWN NOSAVE command.
(1.19s)
```

假如当前正在运行的脚本已经执行过写操作，那么即使执行 `SCRIPT KILL` ，也无法将它杀死，因为这是违反 Lua 脚本的原子性执行原则的。在这种情况下，唯一可行的办法是使用 `SHUTDOWN NOSAVE` 命令，通过停止整个 Redis 进程来停止脚本的运行，并防止不完整(half-written)的信息被写入数据库中。  

#### SCRIPT DEBUG

redis 从 v3.2.0 开始支持 lua debugger，可以加断点、print 变量信息、调试正在执行的代码......  

如何进入调试模式？  

在原本执行的命令中增加 `--ldb` 即可进入调试模式。  

栗子  

```
# redis-cli --ldb  --eval ./test.lua  key1 key2 ,  arg1 arg2 arg3
Lua debugging session started, please use:
quit    -- End the session.
restart -- Restart the script in debug mode again.
help    -- Show Lua script debugging commands.

* Stopped at 1, stop reason = step over
-> 1   local key1   = tostring(KEYS[1])
```

调试模式有两种，同步模式和调试模式：  

1、调试模式：使用 `--ldb` 开启，调试模式下 Redis 会 fork 一个进程进去到隔离环境中，不会影响到 Redis 中的正常执行，同样 Redis 中正常命令的执行也不会影响到调试模式，两者相互隔离，同时调试模式下，调试脚本结束时，回滚脚本操作的所有数据更改。  

2、同步模式：使用 `--ldb-sync-mode` 开启，同步模式下，会阻塞 Redis 中的命令，完全模拟了正常模式下的命令执行，调试命令的执行结果也会被记录。在此模式下调试会话期间，Redis 服务器将无法访问，因此需要谨慎使用。   

这里简单下看下，Redis 中如何进行调试  

看下 debugger 模式支持的命令  

```
lua debugger> h
Redis Lua debugger help:
[h]elp               Show this help.
[s]tep               Run current line and stop again.
[n]ext               Alias for step.
[c]continue          Run till next breakpoint.
[l]list              List source code around current line.
[l]list [line]       List source code around [line].
                     line = 0 means: current position.
[l]list [line] [ctx] In this form [ctx] specifies how many lines
                     to show before/after [line].
[w]hole              List all source code. Alias for 'list 1 1000000'.
[p]rint              Show all the local variables.
[p]rint <var>        Show the value of the specified variable.
                     Can also show global vars KEYS and ARGV.
[b]reak              Show all breakpoints.
[b]reak <line>       Add a breakpoint to the specified line.
[b]reak -<line>      Remove breakpoint from the specified line.
[b]reak 0            Remove all breakpoints.
[t]race              Show a backtrace.
[e]eval <code>       Execute some Lua code (in a different callframe).
[r]edis <cmd>        Execute a Redis command.
[m]axlen [len]       Trim logged Redis replies and Lua var dumps to len.
                     Specifying zero as <len> means unlimited.
[a]bort              Stop the execution of the script. In sync
                     mode dataset changes will be retained.

Debugger functions you can call from Lua scripts:
redis.debug()        Produce logs in the debugger console.
redis.breakpoint()   Stop execution like if there was a breakpoing.
                     in the next line of code.
```

这里来个简单的分析  

```
# cat test.lua
local key1   = tostring(KEYS[1])
local key2   = tostring(KEYS[2])
local arg1   = tostring(ARGV[1])

if key1 == 'test1' then
   return 1
end

if key2 == 'test2' then
   return 2
end

return arg1

# 进入 debuge 模式
# redis-cli --ldb  --eval ./test.lua  key1 key2 ,  arg1 arg2 arg3
Lua debugging session started, please use:
quit    -- End the session.
restart -- Restart the script in debug mode again.
help    -- Show Lua script debugging commands.

* Stopped at 1, stop reason = step over
-> 1   local key1   = tostring(KEYS[1])

# 添加断点 
lua debugger> b 3
   2   local key2   = tostring(KEYS[2])
  #3   local arg1   = tostring(ARGV[1])
   4
   
# 打印输入的参数 key
lua debugger> p KEYS
<value> {"key1"; "key2"}
```

### 为什么 Redis 中的 Lua 脚本的执行是原子性的

我们知道 Redis 中的单命令的执行是原子性的，因为命令的执行都是单线程去处理的。  

那么对于 Redis 中执行 Lua 脚本也是原子性的，是如何实现的呢？这里来探讨下。  

Redis 使用单个 Lua 解释器去运行所有脚本，并且， Redis 也保证脚本会以原子性(atomic)的方式执行： 当某个脚本正在运行的时候，不会有其他脚本或 Redis 命令被执行。 这和使用 MULTI / EXEC 包围的事务很类似。 在其他别的客户端看来，脚本的效果(effect)要么是不可见的(not visible)，要么就是已完成的(already completed)。  

这里看下里面核心 EVAL 的实现  

```
// https://github.com/redis/redis/blob/6.2/src/scripting.c#L1490
void evalCommand(client *c) {
    replicationFeedMonitors(c,server.monitors,c->db->id,c->argv,c->argc);
    if (!(c->flags & CLIENT_LUA_DEBUG))
        evalGenericCommand(c,0);
    else
        evalGenericCommandWithDebugging(c,0);
}

// https://github.com/redis/redis/blob/6.2/src/scripting.c#L1677  
void evalGenericCommand(client *c, int evalsha) {
    lua_State *lua = server.lua;
    char funcname[43];
    long long numkeys;
    long long initial_server_dirty = server.dirty;
    int delhook = 0, err;

    /* When we replicate whole scripts, we want the same PRNG sequence at
     * every call so that our PRNG is not affected by external state. */
    redisSrand48(0);

    /* We set this flag to zero to remember that so far no random command
     * was called. This way we can allow the user to call commands like
     * SRANDMEMBER or RANDOMKEY from Lua scripts as far as no write command
     * is called (otherwise the replication and AOF would end with non
     * deterministic sequences).
     *
     *  - lua_random_dirty 记录脚本是否执行了随机命令
     *
     *  - lua_write_dirty 记录脚本是否进行了写入命令
     *
     * Thanks to this flag we'll raise an error every time a write command
     * is called after a random command was used. */
     // 通过这两个变量，程序可以在脚本试图在执行随机命令之后执行写入时报错。
    server.lua_random_dirty = 0;
    server.lua_write_dirty = 0;
    server.lua_replicate_commands = server.lua_always_replicate_commands;
    server.lua_multi_emitted = 0;
    server.lua_repl = PROPAGATE_AOF|PROPAGATE_REPL;

    /* Get the number of arguments that are keys */
    // 获取输入键的数量
    if (getLongLongFromObjectOrReply(c,c->argv[2],&numkeys,NULL) != C_OK)
        return;
    // 对输入键的正确性做个快速检查
    if (numkeys > (c->argc - 3)) {
        addReplyError(c,"Number of keys can't be greater than number of args");
        return;
    } else if (numkeys < 0) {
        addReplyError(c,"Number of keys can't be negative");
        return;
    }

    // 我们获得脚本SHA1，然后检查这个函数是否已经定义为Lua状态
    funcname[0] = 'f';
    funcname[1] = '_';
    if (!evalsha) {
        // 如果执行的是 EVAL 命令，那么计算脚本的 SHA1 校验和
        sha1hex(funcname+2,c->argv[1]->ptr,sdslen(c->argv[1]->ptr));
    } else {
        // 如果执行的是 EVALSHA 命令，直接使用传入的 SHA1 值
        int j;
        char *sha = c->argv[1]->ptr;

        /* Convert to lowercase. We don't use tolower since the function
         * managed to always show up in the profiler output consuming
         * a non trivial amount of time. */
        for (j = 0; j < 40; j++)
            funcname[j+2] = (sha[j] >= 'A' && sha[j] <= 'Z') ?
                sha[j]+('a'-'A') : sha[j];
        funcname[42] = '\0';
    }

    /* Push the pcall error handler function on the stack. */
    lua_getglobal(lua, "__redis__err__handler");

    // 根据函数名，在 Lua 环境中检查函数是否已经定义
    lua_getfield(lua, LUA_REGISTRYINDEX, funcname);
    if (lua_isnil(lua,-1)) {
        lua_pop(lua,1); /* remove the nil from the stack */
         // 如果执行的函数不存在
         // 如果执行的是 EVALSHA ，返回脚本未找到错误
        if (evalsha) {
            lua_pop(lua,1); /* remove the error handler from the stack. */
            addReplyErrorObject(c, shared.noscripterr);
            return;
        }
        
        // 如果执行的是 EVAL ，那么创建新函数，然后将代码添加到脚本字典中
        if (luaCreateFunction(c,lua,c->argv[1]) == NULL) {
            lua_pop(lua,1); /* remove the error handler from the stack. */
            /* The error is sent to the client by luaCreateFunction()
             * itself when it returns NULL. */
            return;
        }
        /* Now the following is guaranteed to return non nil */
        lua_getfield(lua, LUA_REGISTRYINDEX, funcname);
        serverAssert(!lua_isnil(lua,-1));
    }

     // 将用户传入的键数组和参数数组设为 Lua 环境中的 KEYS 全局变量和 ARGV 全局变量
    luaSetGlobalArray(lua,"KEYS",c->argv+3,numkeys);
    luaSetGlobalArray(lua,"ARGV",c->argv+3+numkeys,c->argc-3-numkeys);

    /* Set a hook in order to be able to stop the script execution if it
     * is running for too much time.
     * We set the hook only if the time limit is enabled as the hook will
     * make the Lua script execution slower.
     *
     * If we are debugging, we set instead a "line" hook so that the
     * debugger is call-back at every line executed by the script. */
    server.in_eval = 1;
    server.lua_caller = c;
    server.lua_cur_script = funcname + 2;
    server.lua_time_start = getMonotonicUs();
    server.lua_time_snapshot = mstime();
    server.lua_kill = 0;
    if (server.lua_time_limit > 0 && ldb.active == 0) {
        lua_sethook(lua,luaMaskCountHook,LUA_MASKCOUNT,100000);
        delhook = 1;
    } else if (ldb.active) {
        lua_sethook(server.lua,luaLdbLineHook,LUA_MASKLINE|LUA_MASKCOUNT,100000);
        delhook = 1;
    }

    prepareLuaClient();

    /* At this point whether this script was never seen before or if it was
     * already defined, we can call it. We have zero arguments and expect
     * a single return value. */
    err = lua_pcall(lua,0,1,-2);

    resetLuaClient();

    /* Perform some cleanup that we need to do both on error and success. */
    if (delhook) lua_sethook(lua,NULL,0,0); /* Disable hook */
    if (server.lua_timedout) {
        server.lua_timedout = 0;
        blockingOperationEnds();
        /* Restore the client that was protected when the script timeout
         * was detected. */
        unprotectClient(c);
        if (server.masterhost && server.master)
            queueClientForReprocessing(server.master);
    }
    server.in_eval = 0;
    server.lua_caller = NULL;
    server.lua_cur_script = NULL;

    /* Call the Lua garbage collector from time to time to avoid a
     * full cycle performed by Lua, which adds too latency.
     *
     * The call is performed every LUA_GC_CYCLE_PERIOD executed commands
     * (and for LUA_GC_CYCLE_PERIOD collection steps) because calling it
     * for every command uses too much CPU. */
    #define LUA_GC_CYCLE_PERIOD 50
    {
        static long gc_count = 0;

        gc_count++;
        if (gc_count == LUA_GC_CYCLE_PERIOD) {
            lua_gc(lua,LUA_GCSTEP,LUA_GC_CYCLE_PERIOD);
            gc_count = 0;
        }
    }

    if (err) {
        addReplyErrorFormat(c,"Error running script (call to %s): %s\n",
            funcname, lua_tostring(lua,-1));
        lua_pop(lua,2); /* Consume the Lua reply and remove error handler. */
    } else {
        /* On success convert the Lua return value into Redis protocol, and
         * send it to * the client. */
        luaReplyToRedisReply(c,lua); /* Convert and consume the reply. */
        lua_pop(lua,1); /* Remove the error handler. */
    }

    /* If we are using single commands replication, emit EXEC if there
     * was at least a write. */
    if (server.lua_replicate_commands) {
        preventCommandPropagation(c);
        if (server.lua_multi_emitted) {
            execCommandPropagateExec(c->db->id);
        }
    }

    /* EVALSHA should be propagated to Slave and AOF file as full EVAL, unless
     * we are sure that the script was already in the context of all the
     * attached slaves *and* the current AOF file if enabled.
     *
     * To do so we use a cache of SHA1s of scripts that we already propagated
     * as full EVAL, that's called the Replication Script Cache.
     *
     * For replication, everytime a new slave attaches to the master, we need to
     * flush our cache of scripts that can be replicated as EVALSHA, while
     * for AOF we need to do so every time we rewrite the AOF file. */
    if (evalsha && !server.lua_replicate_commands) {
        if (!replicationScriptCacheExists(c->argv[1]->ptr)) {
            /* This script is not in our script cache, replicate it as
             * EVAL, then add it into the script cache, as from now on
             * slaves and AOF know about it. */
            robj *script = dictFetchValue(server.lua_scripts,c->argv[1]->ptr);

            replicationScriptCacheAdd(c->argv[1]->ptr);
            serverAssertWithInfo(c,NULL,script != NULL);

            /* If the script did not produce any changes in the dataset we want
             * just to replicate it as SCRIPT LOAD, otherwise we risk running
             * an aborted script on slaves (that may then produce results there)
             * or just running a CPU costly read-only script on the slaves. */
            if (server.dirty == initial_server_dirty) {
                rewriteClientCommandVector(c,3,
                    shared.script,
                    shared.load,
                    script);
            } else {
                rewriteClientCommandArgument(c,0,shared.eval);
                rewriteClientCommandArgument(c,1,script);
            }
            forceCommandPropagation(c,PROPAGATE_REPL|PROPAGATE_AOF);
        }
    }
}
```


### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/   
【EVAL简介】http://www.redis.cn/commands/eval.html   
【Redis Lua 脚本调试器】http://www.redis.cn/topics/ldb.html    


