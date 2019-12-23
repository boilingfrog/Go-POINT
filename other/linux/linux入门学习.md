## linux中的基本命令


### 内部命令和外部命令
内部命令：集成在shell中的
外部命令：和shell是独立的个体

### 命令的执行过程

1、别名
2、内部命令
3、hash表（外部命令）
4、$path（外部命令）

### 命令的别名

编辑配置给的新配置不会立马生效
bash进程重新读取配置文件
````
source /path/to/config_file
. /path/to/config_file
````
如果别名同原命令同名，如果要执行原命令，可使用
````
\echo
'echo'
"echo"
command echo
````