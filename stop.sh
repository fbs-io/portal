###
###
 # @Author: reel
 # @Date: 2023-09-06 06:22:26
 # @LastEditors: reel
 # @LastEditTime: 2023-09-06 06:22:28
 # @Description: 请填写简介
### 
 # @Description: 
 # @Params: 
 # @Author: LenLee
 # @Date: 2022-12-14 23:23:58
 # @LastEditTime: 2022-12-14 23:25:30
 # @LastEditors: LenLee
 # @FilePath: /fbs/stop.sh
### 
kill_pid=`ps -ef | grep ./fbs_portal | awk '{print $2}'`
if [ -n "${kill_pid}" ]
then
    kill -9 ${kill_pid}
    echo "进程已停止"
else
    echo "进程不存在"
fi