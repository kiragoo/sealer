# cmd_proxy

cmd_proxy的作用是需要远程进行一些操作而这些操作又是操作系统的shell不支持的。

如需要在node下发一个ipvs规则，大部分系统是没有ipvsadm的，所以需要seautil这个二进制工具去创建，

因为每个节点都会安装seautil工具。

再如需要在master1和master2上生成证书，也需要seautil这个命令代理去执行

cmd_proxy会去远程调用 seautil 命令，这个命令里包含了如ipvs下发，证书生成等能力, seautil也会放到

cloudrootfs的bin目录中分发到所有节点上，所有系统不自带的一些能力都通过seautil去执行.