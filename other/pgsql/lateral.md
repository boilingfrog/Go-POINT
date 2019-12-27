## pgsql中的lateral

我们先来看官方对lateral的定义

> 可以在出现于FROM中的子查询前放置关键词LATERAL。这允许它们引用前面的FROM项提供的列（如果没有
> LATERAL，每一个子查询将被独立计算，并且因此不能被其他FROM项交叉引用）。