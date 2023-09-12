# rtservice

#更新镜像流程
1、停止并删除原有镜像实例
    docker stop rtservice
    docker rm rtservice
2、删除原有镜像
    docker rmi rtservice:0.1.0
3、加载新的镜像
    docker load -i rtserviceXXXXX.tar
4、启动实例
    docker run 。。。

#导出镜像包命令
docker save -o rtservice.tar wangzhsh/rtservice:0.1.0
#导入镜像包命令
docker load -i rtservice.tar

#run smartlock in docker
docker run -d --name rtservice -p8300:80 -v /root/rtservice/localcache:/services/rtservice/localcache -v /root/rtservice/conf:/services/rtservice/conf wangzhsh/rtservice:0.1.0

