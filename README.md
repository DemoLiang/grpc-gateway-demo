# grpc-gateway-demo
test grpc gateway for grpc &amp; http mutil






### option


```code

-I 或者 --proto_path：用于指定所编译的源码，就是我们所导入的proto文件，支持多次指定，按照顺序搜索，如果未指定，则使用当前工作目录。

--go_out：同样的也有其他语言的，例如--java_out、--csharp_out,用来指定语言的生成位置，用于生成*.pb.go 文件

--go_opt：paths=source_relative 指定--go_out生成文件是基于相对路径的

--go-grpc_out：用于生成 *_grpc.pb.go 文件

--go-grpc_opt：

paths=source_relative 指定--go_grpc_out生成文件是基于相对路径的

require_unimplemented_servers=false 默认是true，会在server类多生成一个接口

--grpc-gateway_out：是使用到了 protoc-gen-grpc-gateway.exe 插件，用于生成pb.gw.go文件

--grpc-gateway_opt：

logtostderr=true 记录log

paths=source_relative 指定--grpc-gateway_out生成文件是基于相对路径的

generate_unbound_methods=true 如果proto文件没有写api接口信息，也会默认生成

--openapiv2_out：使用到了protoc-gen-openapiv2.exe 插件，用于生成swagger.json 文件
当然，还有其他很多命令参数，可以使用protoc -help 查看，也提供了很详细的英文提示。
```