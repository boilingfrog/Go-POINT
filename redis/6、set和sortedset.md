<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [set 和 sorted set](#set-%E5%92%8C-sorted-set)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [set](#set)
    - [常见命令](#%E5%B8%B8%E8%A7%81%E5%91%BD%E4%BB%A4)
    - [set 的使用场景](#set-%E7%9A%84%E4%BD%BF%E7%94%A8%E5%9C%BA%E6%99%AF)
    - [看下源码实现](#%E7%9C%8B%E4%B8%8B%E6%BA%90%E7%A0%81%E5%AE%9E%E7%8E%B0)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## set 和 sorted set

### 前言

前面在几个文章聊到了 `list，string，hash` 等结构的实现，这次来聊一下 set 和 `sorted set` 的细节。  

### set

Redis 的 Set 是 String 类型的无序集合，集合成员是唯一的。  

底层实现主要用到了两种数据结构 hashtable 和 inset(整数集合)。    

集合中最大的成员数为2的32次方-1 (4294967295, 每个集合可存储40多亿个成员)。   

#### 常见命令

来看下几个常用的命令  

```
# 向集合添加一个或多个成员
SADD key member1 [member2]

# 获取集合的成员数
SCARD key

# 返回第一个集合与其他集合之间的差异。
SDIFF key1 [key2]

# 返回给定所有集合的差集并存储在 destination 中
SDIFFSTORE destination key1 [key2]

# 返回给定所有集合的交集
SINTER key1 [key2]

# 返回给定所有集合的交集并存储在 destination 中
SINTERSTORE destination key1 [key2]

# 判断 member 元素是否是集合 key 的成员
SISMEMBER key member

# 返回集合中的所有成员
SMEMBERS key

# 将 member 元素从 source 集合移动到 destination 集合
SMOVE source destination member

# 移除并返回集合中的一个随机元素
SPOP key

# 返回集合中一个或多个随机数
SRANDMEMBER key [count]

# 移除集合中一个或多个成员
SREM key member1 [member2]

# 返回所有给定集合的并集
SUNION key1 [key2]

# 所有给定集合的并集存储在 destination 集合中
SUNIONSTORE destination key1 [key2]

# 迭代集合中的元素
SSCAN key cursor [MATCH pattern] [COUNT count]
```

来个栗子  

```
127.0.0.1:6379>  SADD set-test xiaoming
(integer) 1
127.0.0.1:6379>  SADD set-test xiaoming
(integer) 0
127.0.0.1:6379>  SADD set-test xiaoming1
(integer) 1
127.0.0.1:6379>  SADD set-test xiaoming2


127.0.0.1:6379> SMEMBERS set-test
1) "xiaoming2"
2) "xiaoming"
3) "xiaoming1"
```

上面重复值的插入，只有第一次可以插入成功  

#### set 的使用场景

比较适用于聚合分类  

1、标签：比如我们博客网站常常使用到的兴趣标签，把一个个有着相同爱好，关注类似内容的用户利用一个标签把他们进行归并。  

2、共同好友功能，共同喜好，或者可以引申到二度好友之类的扩展应用。  

3、统计网站的独立IP。利用set集合当中元素不唯一性，可以快速实时统计访问网站的独立IP。  

不过对于 set 中的命令要合理的应用，不然很容易造成慢查询  

1、使用高效的命令，比如说，如果你需要返回一个 SET 中的所有成员时，不要使用 SMEMBERS 命令，而是要使用 SSCAN 多次迭代返回，避免一次返回大量数据，造成线程阻塞。  

2、当你需要执行排序、交集、并集操作时，可以在客户端完成，而不要用`SORT、SUNION、SINTER`这些命令，以免拖慢 Redis 实例。  

#### 看下源码实现  

这里来看下 set 中主要用到的数据类型  

代码路径`https://github.com/redis/redis/blob/6.2/src/t_set.c`  

```
void saddCommand(client *c) {
    robj *set;
    int j, added = 0;

    set = lookupKeyWrite(c->db,c->argv[1]);
    if (checkType(c,set,OBJ_SET)) return;
    
    if (set == NULL) {
        set = setTypeCreate(c->argv[2]->ptr);
        dbAdd(c->db,c->argv[1],set);
    }

    for (j = 2; j < c->argc; j++) {
        if (setTypeAdd(set,c->argv[j]->ptr)) added++;
    }
    if (added) {
        signalModifiedKey(c,c->db,c->argv[1]);
        notifyKeyspaceEvent(NOTIFY_SET,"sadd",c->argv[1],c->db->id);
    }
    server.dirty += added;
    addReplyLongLong(c,added);
}

// 当 value 是整数类型使用 intset
// 否则使用哈希表
robj *setTypeCreate(sds value) {
    if (isSdsRepresentableAsLongLong(value,NULL) == C_OK)
        return createIntsetObject();
    return createSetObject();
}

/* Add the specified value into a set.
 *
 * If the value was already member of the set, nothing is done and 0 is
 * returned, otherwise the new element is added and 1 is returned. */
int setTypeAdd(robj *subject, sds value) {
    long long llval;
    // 如果是 OBJ_ENCODING_HT 说明是哈希类型
    if (subject->encoding == OBJ_ENCODING_HT) {
        dict *ht = subject->ptr;
        dictEntry *de = dictAddRaw(ht,value,NULL);
        if (de) {
            // 设置key ,value 设置成null
            dictSetKey(ht,de,sdsdup(value));
            dictSetVal(ht,de,NULL);
            return 1;
        }
    // OBJ_ENCODING_INTSET 代表是 inset
    } else if (subject->encoding == OBJ_ENCODING_INTSET) {
        if (isSdsRepresentableAsLongLong(value,&llval) == C_OK) {
            uint8_t success = 0;
            subject->ptr = intsetAdd(subject->ptr,llval,&success);
            if (success) {
                /* Convert to regular set when the intset contains
                 * too many entries. */
                // 如果条目过多将会转换成集合
                size_t max_entries = server.set_max_intset_entries;
                /* limit to 1G entries due to intset internals. */
                if (max_entries >= 1<<30) max_entries = 1<<30;
                if (intsetLen(subject->ptr) > max_entries)
                    setTypeConvert(subject,OBJ_ENCODING_HT);
                return 1;
            }
        } else {
            /* Failed to get integer from object, convert to regular set. */
            setTypeConvert(subject,OBJ_ENCODING_HT);

            /* The set *was* an intset and this value is not integer
             * encodable, so dictAdd should always work. */
            serverAssert(dictAdd(subject->ptr,sdsdup(value),NULL) == DICT_OK);
            return 1;
        }
    } else {
        serverPanic("Unknown set encoding");
    }
    return 0;
}
```

通过上面的源码分析，可以看到  

1、set 中主要用到了 hashtable 和 inset；  

2、如果存储的类型是整数类型就会使用 inset，否则使用 hashtable；  

3、使用 inset 有一个最大的限制，达到了最大的限制，也是会使用 hashtable；  

再来看下 inset 数据结构  

代码地址`https://github.com/redis/redis/blob/6.2/src/intset.h`  

```
typedef struct intset {
    // 编码方法，指定当前存储的是 16 位，32 位，还是 64 位的整数
    uint32_t encoding;
    uint32_t length;
    // 实际保存元素的数组
    int8_t contents[];
} intset;
```

来看下 intset 的数据插入  

```
/* Insert an integer in the intset */
intset *intsetAdd(intset *is, int64_t value, uint8_t *success) {
    // 计算value的编码长度
    uint8_t valenc = _intsetValueEncoding(value);
    uint32_t pos;
    if (success) *success = 1;

    /* Upgrade encoding if necessary. If we need to upgrade, we know that
     * this value should be either appended (if > 0) or prepended (if < 0),
     * because it lies outside the range of existing values. */
    // 如果value的编码长度大于当前的编码位数，进行升级
    if (valenc > intrev32ifbe(is->encoding)) {
        /* This always succeeds, so we don't need to curry *success. */
        return intsetUpgradeAndAdd(is,value);
    } else {
        /* Abort if the value is already present in the set.
         * This call will populate "pos" with the right position to insert
         * the value when it cannot be found. */
        // 当前值不存在的时候才进行插入  
        if (intsetSearch(is,value,&pos)) {
            if (success) *success = 0;
            return is;
        }

        is = intsetResize(is,intrev32ifbe(is->length)+1);
        if (pos < intrev32ifbe(is->length)) intsetMoveTail(is,pos,pos+1);
    }

    // 数据插入
    _intsetSet(is,pos,value);
    is->length = intrev32ifbe(intrev32ifbe(is->length)+1);
    return is;
}
```

intset 的数据插入有一个数据升级的过程，当一个整数被添加到整数集合时，首先需要判断下 新元素的类型和集合中现有元素类型的长短，如果新元素是一个32位的数字，现有集合类型是16位的，那么就需要将整数集合进行升级，然后才能将新的元素加入进来。  

这样做的优点：  

1、提升整数集合的灵活性，可以随意的添加整数，而不用关心整数的类型；  

2、可以尽可能的节约内存。  

了解完数据的插入再来看下 intset 中是如何来快速的搜索里面的数据  

```
/* Search for the position of "value". Return 1 when the value was found and
 * sets "pos" to the position of the value within the intset. Return 0 when
 * the value is not present in the intset and sets "pos" to the position
 * where "value" can be inserted. */
// 如果找到了对应的数据，返回 1 将 pos 设置为对应的位置
// 如果找不到，返回0，设置 pos 为可以为数据可以插入的位置
// intset 中的数据是排好序的，所以使用二分查找来寻找对应的元素  
static uint8_t intsetSearch(intset *is, int64_t value, uint32_t *pos) {
    int min = 0, max = intrev32ifbe(is->length)-1, mid = -1;
    int64_t cur = -1;

    /* The value can never be found when the set is empty */
    if (intrev32ifbe(is->length) == 0) {
        if (pos) *pos = 0;
        return 0;
    } else {
        /* Check for the case where we know we cannot find the value,
         * but do know the insert position. */
        if (value > _intsetGet(is,max)) {
            if (pos) *pos = intrev32ifbe(is->length);
            return 0;
        } else if (value < _intsetGet(is,0)) {
            if (pos) *pos = 0;
            return 0;
        }
    }

    // 使用二分查找
    while(max >= min) {
        mid = ((unsigned int)min + (unsigned int)max) >> 1;
        cur = _intsetGet(is,mid);
        if (value > cur) {
            min = mid+1;
        } else if (value < cur) {
            max = mid-1;
        } else {
            break;
        }
    }

    if (value == cur) {
        if (pos) *pos = mid;
        return 1;
    } else {
        if (pos) *pos = min;
        return 0;
    }
}
```

再来看下 dict 数据的插入  

```
dictEntry *dictAddRaw(dict *d, void *key, dictEntry **existing)
{
    long index;
    dictEntry *entry;
    dictht *ht;

    if (dictIsRehashing(d)) _dictRehashStep(d);

    /* Get the index of the new element, or -1 if
     * the element already exists. */
    if ((index = _dictKeyIndex(d, key, dictHashKey(d,key), existing)) == -1)
        return NULL;

    /* Allocate the memory and store the new entry.
     * Insert the element in top, with the assumption that in a database
     * system it is more likely that recently added entries are accessed
     * more frequently. */
    // hash 中得物使用主要用到了
    ht = dictIsRehashing(d) ? &d->ht[1] : &d->ht[0];
    entry = zmalloc(sizeof(*entry));
    entry->next = ht->table[index];
    ht->table[index] = entry;
    ht->used++;

    /* Set the hash entry fields. */
    dictSetKey(d, entry, key);
    return entry;
}
```


### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/    
【redis 集合（set）类型的使用和应用场景】https://www.oraclejsq.com/redisjc/040101720.html    

