## 记一次pgsql的查询优化

### 前言

这是一个子查询的场景，对于这个查询我们不能避免子查询，下面是我一次具体的优化过程。

### 优化策略

1、拆分子查询，将需要的数据提前在cte中查询出来
2、连表查询，直接去查询对应cte里面的内容

````sql
WITH RECURSIVE nodes AS (
  SELECT t.*
  FROM document_directories t
  INNER JOIN document_permissions p
    ON
      p.resource_id = t.id
      AND p.enterprise_id = t.enterprise_id
      AND p.resource ='directory'
      AND p.target = 'role'
  INNER JOIN role_user_rel r
    ON r.role_id = p.target_id AND r.enterprise_id = p.enterprise_id
  WHERE 1 = 1
        AND t.enterprise_id = 11
    AND 1 = ANY(p.perm)
    AND r.uid = 11
    UNION
    SELECT pn.*
    FROM document_directories pn
      INNER JOIN nodes n
        ON pn.id = n.parent_id
    WHERE 1 = 1
  ),
  resJoin AS(
  SELECT t.*, dd.perm, e.* FROM nodes t
  LEFT JOIN LATERAL (
  SELECT array_agg(perm order by perm) perm
    FROM (SELECT DISTINCT unnest(perm) as perm
          FROM document_permissions d
          INNER JOIN roles r ON r.id = d.target_id
        WHERE d.resource = 'directory'
          AND d.target = 'role'
          AND d.enterprise_id = 11
         ) ss
  ) AS dd ON TRUE
    LEFT JOIN  LATERAL(
      SELECT NOT EXISTS(
        SELECT 1 FROM document_directories as dd
          LEFT JOIN documents dm ON dm.directory_id = dd.id
        WHERE (dd.parent_id = t.id OR (dm.id IS NOT NULL AND dd.id = t.id AND dm.deleted = false))
      ) as is_empty_directory
    ) AS e ON TRUE
     ORDER BY t.position DESC
  )
  SELECT
   *
  FROM resJoin r
````
一个RECURSIVE查询出所有的节点信息，后面的resJoin，查询出返回数据需要的信息，里面用到了两个LATERAL，并且里面也用到了子查询。  

分析下这个sql  
````sql
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('CTE Scan on resjoin r  (cost=116084.38..116086.40 rows=101 width=117) (actual time=2423.410..2423.656 rows=432 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('  CTE nodes');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('    ->  Recursive Union  (cost=22.13..1042.04 rows=101 width=71) (actual time=0.309..5.270 rows=432 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('          ->  Nested Loop  (cost=22.13..207.26 rows=1 width=71) (actual time=0.301..1.841 rows=108 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                ->  Nested Loop  (cost=21.84..200.27 rows=1 width=16) (actual time=0.080..0.741 rows=110 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      Join Filter: (p.target_id = r_1.role_id)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      Rows Removed by Join Filter: 110');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      ->  Index Only Scan using roles_users_rel_enterprise_id_uid_role_id_uindex on role_user_rel r_1  (cost=0.15..8.17 rows=1 width=16) (actual time=0.018..0.021 rows=2 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            Index Cond: ((enterprise_id = 11) AND (uid = 11))');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            Heap Fetches: 2');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      ->  Bitmap Heap Scan on document_permissions p  (cost=21.69..190.72 rows=110 width=24) (actual time=0.056..0.295 rows=110 loops=2)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            Recheck Cond: ((enterprise_id = 11) AND (target = ''role''::text) AND (resource = ''directory''::text))');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            Filter: (1 = ANY (perm))');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            Heap Blocks: exact=162');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            ->  Bitmap Index Scan on document_permissions_etr_uindex  (cost=0.00..21.66 rows=110 width=0) (actual time=0.037..0.037 rows=110 loops=2)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                  Index Cond: ((enterprise_id = 11) AND (target = ''role''::text) AND (resource = ''directory''::text))');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                ->  Index Scan using document_directories_pk on document_directories t  (cost=0.29..6.96 rows=1 width=71) (actual time=0.008..0.008 rows=1 loops=110)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      Index Cond: (id = p.resource_id)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      Filter: (enterprise_id = 11)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      Rows Removed by Filter: 0');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('          ->  Nested Loop  (cost=0.29..83.28 rows=10 width=71) (actual time=0.057..0.616 rows=81 loops=4)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                ->  WorkTable Scan on nodes n  (cost=0.00..0.20 rows=10 width=8) (actual time=0.001..0.034 rows=108 loops=4)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                ->  Index Scan using document_directories_pk on document_directories pn  (cost=0.29..8.31 rows=1 width=71) (actual time=0.004..0.004 rows=1 loops=432)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      Index Cond: (id = n.parent_id)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('  CTE resjoin');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('    ->  Sort  (cost=115042.09..115042.34 rows=101 width=117) (actual time=2423.408..2423.451 rows=432 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('          Sort Key: t_1."position" DESC');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('          Sort Method: quicksort  Memory: 139kB');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('          ->  Nested Loop Left Join  (cost=1437.14..115038.72 rows=101 width=117) (actual time=11.306..2422.614 rows=432 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                ->  Nested Loop Left Join  (cost=301.19..304.49 rows=101 width=116) (actual time=0.897..7.130 rows=432 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      ->  CTE Scan on nodes t_1  (cost=0.00..2.02 rows=101 width=84) (actual time=0.310..5.849 rows=432 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      ->  Materialize  (cost=301.19..301.21 rows=1 width=32) (actual time=0.002..0.002 rows=1 loops=432)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            ->  Aggregate  (cost=301.19..301.20 rows=1 width=32) (actual time=0.581..0.581 rows=1 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                  ->  HashAggregate  (cost=298.68..299.93 rows=100 width=4) (actual time=0.558..0.559 rows=6 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                        Group Key: unnest(d.perm)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                        ->  ProjectSet  (cost=46.31..271.18 rows=11000 width=4) (actual time=0.081..0.401 rows=660 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                              ->  Hash Join  (cost=46.31..215.36 rows=110 width=45) (actual time=0.078..0.229 rows=110 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                                    Hash Cond: (d.target_id = r_2.id)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                                    ->  Bitmap Heap Scan on document_permissions d  (cost=21.69..189.35 rows=110 width=53) (actual time=0.053..0.163 rows=110 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                                          Recheck Cond: ((enterprise_id = 11) AND (target = ''role''::text) AND (resource = ''directory''::text))');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                                          Heap Blocks: exact=81');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                                          ->  Bitmap Index Scan on document_permissions_etr_uindex  (cost=0.00..21.66 rows=110 width=0) (actual time=0.039..0.039 rows=110 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                                                Index Cond: ((enterprise_id = 11) AND (target = ''role''::text) AND (resource = ''directory''::text))');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                                    ->  Hash  (cost=16.50..16.50 rows=650 width=8) (actual time=0.011..0.011 rows=2 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                                          Buckets: 1024  Batches: 1  Memory Usage: 9kB');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                                          ->  Seq Scan on roles r_2  (cost=0.00..16.50 rows=650 width=8) (actual time=0.006..0.007 rows=2 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                ->  Result  (cost=1135.95..1135.96 rows=1 width=1) (actual time=5.589..5.590 rows=1 loops=432)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      InitPlan 2 (returns $5)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                        ->  Nested Loop Left Join  (cost=0.00..1135.95 rows=1 width=0) (actual time=5.588..5.588 rows=1 loops=432)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                              Join Filter: (dm.directory_id = dd.id)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                              Rows Removed by Join Filter: 42');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                              Filter: ((dd.parent_id = $4) OR ((dm.id IS NOT NULL) AND (dd.id = $4) AND (NOT dm.deleted)))');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                              Rows Removed by Filter: 1');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                              ->  Seq Scan on document_directories dd  (cost=0.00..1134.06 rows=2 width=16) (actual time=4.287..5.575 rows=2 loops=432)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                    Filter: ((parent_id = $4) OR (id = $4))');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                    Rows Removed by Filter: 23789');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                              ->  Materialize  (cost=0.00..1.25 rows=17 width=17) (actual time=0.000..0.002 rows=24 loops=756)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                    ->  Seq Scan on documents dm  (cost=0.00..1.17 rows=17 width=17) (actual time=0.005..0.011 rows=24 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('Planning time: 1.665 ms');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('Execution time: 2424.040 ms');
````

可以看到性能瓶颈就在resjoin，也就是下面的两个LATERAL查询。这两个查询是不可避免的，所以也不能拿掉，但是pgsql里面的cte支持并行的查询，
我们可以先把这些资源查询出来，然后在连表查询  

````sql
WITH RECURSIVE nodes AS (
  SELECT t.*
  FROM document_directories t
  INNER JOIN document_permissions p
    ON
      p.resource_id = t.id
      AND p.enterprise_id = t.enterprise_id
      AND p.resource ='directory'
      AND p.target = 'role'
  INNER JOIN role_user_rel r
    ON r.role_id = p.target_id AND r.enterprise_id = p.enterprise_id
  WHERE 1 = 1
        AND t.enterprise_id = 11
    AND 1 = ANY(p.perm)
    AND r.uid = 11
    UNION
    SELECT pn.*
    FROM document_directories pn
      INNER JOIN nodes n
        ON pn.id = n.parent_id
    WHERE 1 = 1
  ),
  perms AS (
    SELECT d.perm, d.resource_id
    FROM document_permissions d
      INNER JOIN roles r ON r.id = d.target_id
    WHERE d.resource = 'directory'
      AND d.target = 'role'
      AND d.enterprise_id = 11
    GROUP BY d.resource_id, d.perm
  ),
  directory_exists AS (
    SELECT distinct dd.id
    FROM document_directories as dd
      LEFT JOIN document_directories t on t.parent_id = dd.id
      LEFT JOIN documents dm ON dm.directory_id = dd.id
    WHERE (t.id IS NOT NULL OR (dm.id IS NOT NULL AND dm.deleted = false))
 ),
  resJoin AS(
  SELECT t.*, dd.perm,de.id IS NULL as is_empty_directory FROM nodes t
  LEFT JOIN LATERAL (
    SELECT array_agg(perm order by perm) perm
      FROM (SELECT DISTINCT unnest(perm) as perm
        FROM perms d
        WHERE
          d.resource_id = t.id
    ) AS t
  ) AS dd ON TRUE
  LEFT JOIN directory_exists de ON de.id = t.id
     ORDER BY t.position
  )
  SELECT
    *
  FROM resJoin r
````
分析下结果  

````sql
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('CTE Scan on resjoin r  (cost=8366.56..8770.60 rows=20202 width=117) (actual time=118.496..118.699 rows=432 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('  CTE nodes');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('    ->  Recursive Union  (cost=22.13..1042.04 rows=101 width=71) (actual time=0.113..2.833 rows=432 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('          ->  Nested Loop  (cost=22.13..207.26 rows=1 width=71) (actual time=0.109..0.913 rows=108 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                ->  Nested Loop  (cost=21.84..200.27 rows=1 width=16) (actual time=0.082..0.458 rows=110 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      Join Filter: (p.target_id = r_1.role_id)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      Rows Removed by Join Filter: 110');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      ->  Index Only Scan using roles_users_rel_enterprise_id_uid_role_id_uindex on role_user_rel r_1  (cost=0.15..8.17 rows=1 width=16) (actual time=0.019..0.021 rows=2 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            Index Cond: ((enterprise_id = 11) AND (uid = 11))');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            Heap Fetches: 2');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      ->  Bitmap Heap Scan on document_permissions p  (cost=21.69..190.72 rows=110 width=24) (actual time=0.052..0.188 rows=110 loops=2)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            Recheck Cond: ((enterprise_id = 11) AND (target = ''role''::text) AND (resource = ''directory''::text))');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            Filter: (1 = ANY (perm))');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            Heap Blocks: exact=162');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            ->  Bitmap Index Scan on document_permissions_etr_uindex  (cost=0.00..21.66 rows=110 width=0) (actual time=0.035..0.035 rows=110 loops=2)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                  Index Cond: ((enterprise_id = 11) AND (target = ''role''::text) AND (resource = ''directory''::text))');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                ->  Index Scan using document_directories_pk on document_directories t  (cost=0.29..6.96 rows=1 width=71) (actual time=0.003..0.003 rows=1 loops=110)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      Index Cond: (id = p.resource_id)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      Filter: (enterprise_id = 11)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      Rows Removed by Filter: 0');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('          ->  Nested Loop  (cost=0.29..83.28 rows=10 width=71) (actual time=0.050..0.348 rows=81 loops=4)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                ->  WorkTable Scan on nodes n  (cost=0.00..0.20 rows=10 width=8) (actual time=0.000..0.015 rows=108 loops=4)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                ->  Index Scan using document_directories_pk on document_directories pn  (cost=0.29..8.31 rows=1 width=71) (actual time=0.002..0.002 rows=1 loops=432)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      Index Cond: (id = n.parent_id)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('  CTE perms');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('    ->  HashAggregate  (cost=215.91..217.01 rows=110 width=53) (actual time=0.297..0.321 rows=110 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('          Group Key: d.resource_id, d.perm');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('          ->  Hash Join  (cost=46.31..215.36 rows=110 width=53) (actual time=0.065..0.216 rows=110 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                Hash Cond: (d.target_id = r_2.id)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                ->  Bitmap Heap Scan on document_permissions d  (cost=21.69..189.35 rows=110 width=61) (actual time=0.044..0.150 rows=110 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      Recheck Cond: ((enterprise_id = 11) AND (target = ''role''::text) AND (resource = ''directory''::text))');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      Heap Blocks: exact=81');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      ->  Bitmap Index Scan on document_permissions_etr_uindex  (cost=0.00..21.66 rows=110 width=0) (actual time=0.031..0.031 rows=110 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            Index Cond: ((enterprise_id = 11) AND (target = ''role''::text) AND (resource = ''directory''::text))');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                ->  Hash  (cost=16.50..16.50 rows=650 width=8) (actual time=0.010..0.010 rows=2 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      Buckets: 1024  Batches: 1  Memory Usage: 9kB');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      ->  Seq Scan on roles r_2  (cost=0.00..16.50 rows=650 width=8) (actual time=0.005..0.006 rows=2 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('  CTE directory_exists');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('    ->  HashAggregate  (cost=3474.64..3874.68 rows=40004 width=8) (actual time=80.055..88.016 rows=30003 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('          Group Key: dd.id');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('          ->  Hash Right Join  (cost=3373.22..3374.63 rows=40004 width=8) (actual time=58.189..68.144 rows=30003 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                Hash Cond: (dm.directory_id = dd.id)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                Filter: ((t_1.id IS NOT NULL) OR ((dm.id IS NOT NULL) AND (NOT dm.deleted)))');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                Rows Removed by Filter: 10001');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                ->  Seq Scan on documents dm  (cost=0.00..1.17 rows=17 width=17) (actual time=0.012..0.018 rows=24 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                ->  Hash  (cost=2873.17..2873.17 rows=40004 width=16) (actual time=57.854..57.854 rows=40004 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      Buckets: 65536  Batches: 1  Memory Usage: 2310kB');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      ->  Hash Right Join  (cost=1434.09..2873.17 rows=40004 width=16) (actual time=18.692..45.374 rows=40004 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            Hash Cond: (t_1.parent_id = dd.id)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            ->  Seq Scan on document_directories t_1  (cost=0.00..934.04 rows=40004 width=16) (actual time=0.010..5.319 rows=40004 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            ->  Hash  (cost=934.04..934.04 rows=40004 width=8) (actual time=18.373..18.373 rows=40004 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                  Buckets: 65536  Batches: 1  Memory Usage: 2075kB');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                  ->  Seq Scan on document_directories dd  (cost=0.00..934.04 rows=40004 width=8) (actual time=0.008..7.829 rows=40004 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('  CTE resjoin');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('    ->  Sort  (cost=3182.33..3232.83 rows=20202 width=117) (actual time=118.494..118.535 rows=432 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('          Sort Key: t_2."position"');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('          Sort Method: quicksort  Memory: 99kB');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('          ->  Hash Right Join  (cost=585.55..1737.66 rows=20202 width=117) (actual time=95.397..118.216 rows=432 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                Hash Cond: (de.id = t_2.id)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                ->  CTE Scan on directory_exists de  (cost=0.00..800.08 rows=40004 width=8) (actual time=80.058..97.650 rows=30003 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                ->  Hash  (cost=584.28..584.28 rows=101 width=116) (actual time=15.277..15.277 rows=432 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      Buckets: 1024  Batches: 1  Memory Usage: 63kB');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                      ->  Nested Loop Left Join  (cost=5.74..584.28 rows=101 width=116) (actual time=0.519..15.044 rows=432 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            ->  CTE Scan on nodes t_2  (cost=0.00..2.02 rows=101 width=84) (actual time=0.115..3.140 rows=432 loops=1)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                            ->  Aggregate  (cost=5.74..5.75 rows=1 width=32) (actual time=0.027..0.027 rows=1 loops=432)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                  ->  HashAggregate  (cost=3.23..4.48 rows=100 width=4) (actual time=0.023..0.024 rows=2 loops=432)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                        Group Key: unnest(d_1.perm)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                        ->  ProjectSet  (cost=0.00..2.98 rows=100 width=4) (actual time=0.019..0.022 rows=2 loops=432)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                              ->  CTE Scan on perms d_1  (cost=0.00..2.48 rows=1 width=32) (actual time=0.019..0.021 rows=0 loops=432)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                                    Filter: (resource_id = t_2.id)');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('                                                    Rows Removed by Filter: 110');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('Planning time: 1.518 ms');
INSERT INTO "MY_TABLE"("QUERY PLAN") VALUES ('Execution time: 121.197 ms');
````

发现有了质的飞跃