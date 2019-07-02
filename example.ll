; ModuleID = 'expectoroot'
source_filename = "expectoroot"

@strtmp = private unnamed_addr constant [11 x i8] c"that's it\0A\00", align 1
@str = private unnamed_addr constant [10 x i8] c"that's it\00", align 1
@strtmp.1 = private unnamed_addr constant [10 x i8] c"I'm done\0A\00", align 1
@str.2 = private unnamed_addr constant [9 x i8] c"I'm done\00", align 1

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
  %puts = call i32 @puts(i8* getelementptr inbounds ([9 x i8], [9 x i8]* @str.2, i64 0, i64 0))
  ret float 9.000000e+00
}

attributes #0 = { nounwind }
