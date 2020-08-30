# spatula
works with pan
## usage
### parameters
-l the share link, must not be empty    
-o output file name, will acquire from link if not provided    
-p out put path, if not provided, the default path is the location where the binary is located    
-u custom useragent, default value is empty. Use this parameter at your own risk.     
-c size of each file block(in Bytes), default value is 4194304(4 MB),if download speed is limited, set a value below the limit.    
-g max goroutines opened, default value is 20, value over 1000 is not recommended

### examples
#### windows
```shell script
.\stapula.exe -l https://oxygenos.oneplus.net/OnePlus6Oxygen_22_OTA_047_all_2007191515_bd6f7476887846cb.zip -c 4194304 -g 20
```
#### linux
not well tested
```shell script
./stapula -l https://oxygenos.oneplus.net/OnePlus6Oxygen_22_OTA_047_all_2007191515_bd6f7476887846cb.zip
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
