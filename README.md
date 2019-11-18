# xrm
Remove all comments from c, c++, go, js.

## Build
```
git clone https://github.com/XieXiaomei-ptr/xrm.git
cd xrm
go build
```

## Usage
```
./xrm
  -dir string
    	The directory in which your code base is located. Multiple paths are supported.
    	When multiple paths are specified at the same time, a separate #.
    	For example: -dir=/data/my_project/code1#/data/my_project/code2
  -lan string
    	Specify the programming language contained in dir, which supports c, c++, go, and js.
    	If you specify multiple languages at the same time, you need to separate with #.
    	For Example: -lan=c#c++#js#go
  -v	Version
```