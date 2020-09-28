## 1、打包成可在Windows运行的exe文件
#### 进入到主文件（main.go）目录下，执行以下命令：
``go build main.go``
#### 就会在当前目录下生成打包好的可直接运行的EXE文件

## 2、打包成可在Linux运行的文件
> 打包成二进制文件，可在Linux平台运行
#### 首先进到主文件（mian.go)目录下，执行以下命令：
``set GOARCH=amd64``  
``set GOOS=linux``
> GOOS 指的是目标操作系统，支持以下操作系统：
> darwin freebsd linux windows android dragonfly netbsd openbsd plan9 solaris

> GOARCH 指的目标处理器架构，支持以下处理器结构
> arm arm64 386 amd64 ppc64 ppc64le mips64 mips64le s390x

#### 设置好目标操作系统与目标处理器架构后，针对主文件（main.go）执行go build
``go build contrast.go json_compare.go settings.go``
#### 之后就会在当前目录生成打包好的Go项目文件了，是Linux平台可执行的二进制文件

#### 将该文件放入linux系统某个文件夹下，chmod 773 [文件名] 赋予文件权限，./xx 命令即可执行文件，不需要go的任何依赖，就可以直接运行了。

## 运行
#### ide中运行
``go test --filepath replay_middle_rep.log --diffresult result.log``
