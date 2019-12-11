## json数组查询结果变成了字符串


### 场景复原

最近使用到了json的数组，用来存储多个文件的值，发现在连表查询的时候返回结果变成了字符串。
````
 {
            "id": "repl-placeholder-007",
            "sn": "63165580943163393",
            "name": "1212",
            "implementPlanID": "2632920653191188481",
            "auditLeaderID": "1000",
            "auditLeaderInfo": null,
            "status": "draft",
            "deleted": false,
            "userID": "10000",
            "userInfo": null,
            "appInstanceID": "88888888",
            "createdAt": "repl-placeholder-008",
            "updatedAt": "repl-placeholder-009",
            "auditImplementPlan": {
                "id": "repl-placeholder-010",
                "sn": "",
                "auditName": "1212",
                "createdAt": "repl-placeholder-011",
                "updatedAt": "repl-placeholder-012",
                "annualPlanID": "1",
                "auditType": "all",
                "auditLeaderID": "12",
                "auditDate": "repl-placeholder-031",
                "remark": "",
                "deleted": false,
                "createdBy": "1000",
                "appInstanceID": "1212",
                "attaches": "[\"jimeng.io\",\"baidu.com\"]"
                "attaches": [
                    "jimeng.io",
                    "baidu.com"
                ]
            },
            "attaches": "[\"jimeng.io\",\"baidu.com\"]"
            "attaches": [
                "打豆豆.io",
                "baidu.com"
            ]
        }
````

我们发现attaches被转换成了字符串，但是我attaches字段明明定义的是json类型的，但是返回
结果变成了字符串。

我们来看下数据库的字段 
````
 ADD COLUMN "attaches" text NOT NULL DEFAULT '[]'::jsonb;
````
可以看到用的是json类型

还有就是在查询的时候使用了to_json

这是to_json的函数的文档描述
****
把值返回为json或者jsonb。数组和组合被（递归地）转换成数组和对象；否则， 如果有从该类型到json的投影，将使用该投影函数来执行转换； 否则将产生一个标量值。对任何一个数值、布尔量或空值的标量类型， 将使用其文本表达，以这样一种方式使其成为有效的json或者jsonb值。~~~~
****
所以这是正常的情况，但是我们需要的是以数组的形式输出。这是发现 attaches字段用的是text字段，也就是文本字段，他可能就是导致问题出现的原因，
于是更改了字段的类型为jsonb，发现解决了，attaches的输出已经正常了。

当然这种操作还有一个改进的办法，就是使用数组，而不是json数组，这样也不会出现这些问题了。
````
ADD COLUMN "attachess" text[] DEFAULT '{}'::text[];
````