#!/bin/bash

## 通过scan命令，输出Redis中，所有永久缓存的keys

db_ip=127.0.0.1      # redis ip
db_port=6379         # redis 端口
password=''          # redis 密码
cursor=0             # 第一次游标
cnt=1000             # 每次迭代的数量
new_cursor=0         # 下一次游标

cli=/apps/svr/redis-2.8.19/bin/redis-cli   # redis-cli 工具路径

counter=1            # 统计总keys数量

${cli} -h ${db_ip} -p ${db_port} -a ${password} scan ${cursor} count ${cnt} > scan_tmp_result

new_cursor=`sed -n '1p' scan_tmp_result`               # 获取下一次游标
sed -n '2,$p' scan_tmp_result > scan_result            # 获取 keys
cat scan_result | while read line                      # 循环遍历所有 keys
do
    ttl_result=`${cli} -h ${db_ip} -p ${db_port} -a ${password} ttl ${line}`      # 获取 key 过期时间
    if [[ ${ttl_result} == -1 ]];then                  # 判断过期时间，-1是不过期
        counter=$(($counter+1))
        echo ${line} >> no_ttl.log                     # 追加到指定文件
        echo "`date "+%Y-%m-%d %H:%M:%S"`: ${line} [${counter}]"
    fi
done


while [[ ${cursor} -ne ${new_cursor} ]];               # 若游标不为0 ，则证明没有迭代完所有的key，继续执行
do
    ${cli} -h ${db_ip} -p ${db_port} -a ${password} scan ${new_cursor} count ${cnt} > scan_tmp_result
    new_cursor=`sed -n '1p' scan_tmp_result`
    sed -n '2,$p' scan_tmp_result > scan_result
    cat scan_result | while read line
    do
        ttl_result=`${cli} -h ${db_ip} -p ${db_port} -a ${password} ttl ${line}`
        if [[ ${ttl_result} == -1 ]];then
            counter=$(($counter+1))
            echo ${line} >> no_ttl.log
            echo "`date "+%Y-%m-%d %H:%M:%S"`: ${line} [${counter}]"
        fi
    done
done

rm -rf scan_tmp_result
rm -rf scan_result