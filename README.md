# winapi #
这是用 Go 编写的调用 Windows API 的包。将 Windows API 转换成 Go 的风格，使其易用。

## 使用方法 ##

本包不依赖任何第三方包、框架或库。  
任何安装了 Go 环境的 Windows 系统都可以使用。

## 已实现的或有所改变的 API ##

没有照搬 Windows API，而是在做了一些妥善的改动。例如，不封装`CreateWindow`而封装 `CreateWindowEx`，封装之后的Go的函数名为`CreateWindow`。

### gdi ###
- BitBlt
- DeleteObject
- GetObject
- CreateCompatibleDC
- SelectObject
- DeleteDC

### kernel ###
- GetLastError
- ExitProcess
- CreateFile
- ReadFile
- WriteFile
- SetFilePointer
- GetModuleHandle
- CloseHandle
- FormatMessage

### user ###
- DefWindowProc
- GetMessage
- RegisterClass
- MessageBox
- CreateWindow
- ShowWindow
- UpdateWindow
- TranslateMessage
- DispatchMessage
- PostQuitMessage
- DestroyWindow
- LoadString
- LoadIcon
- LoadCursor
- LoadBitmap
- LoadImage
- BeginPaint
- EndPaint

## 不实现的函数 ##
与线程、同步原语相关的函数，如`CreateThread`、`CreateMutex`，因为Go已经有了很好的并发特性。

## TODO ##
COM，注册表，DirectX 等。

## 协议 ##
Apache 2
