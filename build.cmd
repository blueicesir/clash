REM ͨ�����ߵ�pngת��Ϊico��ʽ����������clash����һ��ͼ��
rsrc -manifest .\rsc\clash.manifest -ico .\rsc\clash.ico -o .\rsrc.syso
REM go build  -ldflags="-H windowsgui"
go build
@PAUSE