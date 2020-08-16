# spatula
works with pan
## usage
### windows
```shell script
.\stapula.exe "链接： https://pan.baidu.com/s/xxxxxxxxxx 提取码：xxxx"
```
### linux
not tested
```shell script
./stapula "链接： https://pan.baidu.com/s/xxxxxxxxxx 提取码：xxxx"
```
## common questions
Q: Why I can't change the location where the file is downloaded?  
A: this is a feature under consideration, and since this project has just begun, more features will be added in the future.

Q: Why I can't monitor download speed in the console?    
A: For now, the software focuses on stability and performance. For a software that holds hundreds or even thousands of goroutines, it's hard to calculate the network usage of every goroutine without lowering performance. If you really want to monitor the speed, please use the task manager.    

Q: I encountered a problem. What should I do?    
A: Open an issue at git hub issue. It's the best to have the error log attached, and describe the situation as detail as possible.    

Q: I have an awesome idea and want to contribute    
A: PR is welcome, just open a pull request.
