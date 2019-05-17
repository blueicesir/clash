REM 通过在线的png转换为ico格式，给命令行clash增加一个图标
rsrc -manifest .\rsc\clash.manifest -ico .\rsc\clash.ico -o .\rsrc.syso
REM go build  -ldflags="-H windowsgui"
go build
@PAUSE