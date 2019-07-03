; ModuleID = 'expectoroot'
source_filename = "expectoroot"

@strtmp = private unnamed_addr constant [11 x i8] c"that's it\0A\00", align 1
@str = private unnamed_addr constant [10 x i8] c"that's it\00", align 1
@strtmp.1 = private unnamed_addr constant [10 x i8] c"I'm done\0A\00", align 1
@strtmp.2 = private unnamed_addr constant [3 x i8] c"Hi\00", align 1
@strtmp.3 = private unnamed_addr constant [6 x i8] c"hello\00", align 1
@strtmp.4 = private unnamed_addr constant [5 x i8] c"test\00", align 1
@strtmp.5 = private unnamed_addr constant [2 x i8] c"\0A\00", align 1
@str.6 = private unnamed_addr constant [9 x i8] c"I'm done\00", align 1

declare float @printf(i8* %x)

define float @call() {
entry:
  %puts = call i32 @puts(i8* getelementptr inbounds ([10 x i8], [10 x i8]* @str, i64 0, i64 0))
  ret float 0.000000e+00
}

; Function Attrs: nounwind
declare i32 @puts(i8* nocapture readonly) #0

define float @main() {
entry:
  %calltmp = call float @call()
  %puts = call i32 @puts(i8* getelementptr inbounds ([9 x i8], [9 x i8]* @str.6, i64 0, i64 0))
  %calltmp6 = call float @printf(i8* getelementptr inbounds ([6 x i8], [6 x i8]* @strtmp.3, i64 0, i64 0))
  %putchar = call i32 @putchar(i32 10)
  ret float 1.200000e+01
}

; Function Attrs: nounwind
declare i32 @putchar(i32) #0

attributes #0 = { nounwind }
