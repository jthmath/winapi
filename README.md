# winapi #
这是用Go语言编写的调用Windows API的包。将Windows API转换成Go语言的风格（如返回`BOOL`的函数改为返回`error`），使其易用。

## 使用方法 ##
本包不依赖任何第三方包、框架或库。<br>
任何安装了Go环境的Windows系统都可以使用。

## 已实现的或有所改变的 API ##
我没有照搬 Windows API，而是在做了一些妥善的改动。例如，不封装`CreateWindow`而封装 `CreateWindowEx`，封装之后的Go的函数名为`CreateWindow`。

### gdi ###
BitBlt<br>
DeleteObject<br>
GetObject<br>
CreateCompatibleDC<br>
SelectObject<br>
DeleteDC<br>

### kernel ###
GetLastError<br>
ExitProcess<br>
CreateFile<br>
ReadFile<br>
WriteFile<br>
SetFilePointer<br>
GetModuleHandle<br>
CloseHandle<br>
FormatMessage<br>

### user ###
DefWindowProc<br>
GetMessage<br>
RegisterClass<br>
MessageBox<br>
CreateWindow<br>
ShowWindow<br>
UpdateWindow<br>
TranslateMessage<br>
DispatchMessage<br>
PostQuitMessage<br>
DestroyWindow<br>
LoadString<br>
LoadIcon - 分解为 LoadIconById 和 LoadIconByName<br>
LoadCursor - 分解为 LoadCursorById 和 LoadCursorByName<br>
LoadBitmap - 分解为 LoadBitmapById 和 LoadBitmapByName (todo)<br>
LoadImage - 分解为 LoadImageById 和 LoadImageByName<br>
BeginPaint<br>
EndPaint<br>

### comdlg ###
1

## 增加的函数 ##
ErrorBox<br>

## 不实现的函数 ##
1. 与线程相关的函数，如CreateThread，因为Go已经有了很好的并发特性。<br>
2. 同步有关的函数，如CreateMutex。

## TODO ##
实现COM接口的封装，例如任务栏的进度条、DirectX等。


## 协议 ##
本项目采用与`golang`相同的协议。