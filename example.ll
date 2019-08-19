; ModuleID = 'expectoroot'
source_filename = "expectoroot"

@strtmp = private unnamed_addr constant [14 x i8] c"IterThroughMe\00", align 1
@0 = private unnamed_addr constant [49 x i8] c"Panic: right side of the equation is equal to 0\0A\00", align 1
@strtmp.1 = private unnamed_addr constant [5 x i8] c"OWO\0A\00", align 1
@strtmp.2 = private unnamed_addr constant [3 x i8] c"Hi\00", align 1
@str = private unnamed_addr constant [4 x i8] c"OWO\00", align 1

declare i32 @printf(i8* %x)

declare void @exit(i32 %x)

define void @main() {
entry:
  br label %loop

loop:                                             ; preds = %loop, %entry
  %ind = phi i32 [ 0, %entry ], [ %nextind, %loop ]
  %puts = call i32 @puts(i8* getelementptr inbounds ([4 x i8], [4 x i8]* @str, i64 0, i64 0))
  %nextind = add i32 %ind, 1
  %loopcond = icmp eq i32 %nextind, 14
  br i1 %loopcond, label %exitloop, label %loop

exitloop:                                         ; preds = %loop
  %0 = call i32 @printf(i8* getelementptr inbounds ([3 x i8], [3 x i8]* @strtmp.2, i64 0, i64 0))
  ret void
}

; Function Attrs: nounwind
declare i32 @puts(i8* nocapture readonly) #0

attributes #0 = { nounwind }
