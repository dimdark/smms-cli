使用 `Go` 语言编写 [sm.ms](https://sm.ms) 的命令行客户端

主要参考项目 `sm_ms_api`(https://github.com/sndnvaps/sm_ms_api) , 基本按照(抄袭)其思路(代码很类似),感谢该项目的作者!

- 安装

  ```shell
  go get github.com/dimdark/smms-cli
  cd ${GOPATH}/src/github.com/dimdark/smms-cli
  go build -o smms main.go
  ```

  即可得到可执行文件 `smms`

  

- 使用

   - 上传图片到 [sm.ms](https://sm.ms)

     ```shell
     smms upload xxx.jpg
     or 
     smms u xxx.jpg
     ```

     > xxx.jpg必须在命令执行的当前目录, 否则必须使用图片的绝对路径

     返回结果类似:

     ```shell
     # markdown 格式
     md: ![02.jpg](https://i.loli.net/2019/05/18/5cdffe2d0abf314528.jpg)
     # 使用该url来删除 这张上传到[sm.ms](https://sm.ms)的图片 
     del: https://sm.ms/delete/Ig6yNDdz32e4OWS
     ```

     

  - 删除上传到 [sm.ms](https://sm.ms) 上的图片

     ```shell
    # 后接特定的删除url
    smms delete https://sm.ms/delete/Ig6yNDdz32e4OWS
    or 
    smms d https://sm.ms/delete/Ig6yNDdz32e4OWS
     ```

    

  - 列出上传到 [sm.ms](https://sm.ms) 的图片的历史记录

    ```shell
    smms list
    or 
    smms l
    ```

    

  - 清除上传到 [sm.ms](https://sm.ms) 的图片的历史记录

    ```shell
    smms clear
    or 
    smms c
    ```

    

- 文档

  - [sm.ms](https://sm.ms) 的`API`文档

    ```
    https://sm.ms/doc
    ```

    

​			