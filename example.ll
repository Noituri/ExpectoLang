; ModuleID = 'expectoroot'
source_filename = "expectoroot"

@strtmp = private unnamed_addr constant [11 x i8] c"that's it\0A\00", align 1
@str = private unnamed_addr constant [10 x i8] c"that's it\00", align 1

declare float @printf(i8* %x)

define void @call() {
entry:
  %puts = call i32 @puts(i8* getelementptr inbounds ([10 x i8], [10 x i8]* @str, i64 0, i64 0))
  ret void
}

; Function Attrs: nounwind
declare i32 @puts(i8* nocapture readonly) #0

define void @main() {
entry:
  call void @call()
  ret void
}
