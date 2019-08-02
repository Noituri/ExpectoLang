; ModuleID = 'expectoroot'
source_filename = "expectoroot"

@strtmp = private unnamed_addr constant [14 x i8] c"IterThroughMe\00", align 1
@strtmp.1 = private unnamed_addr constant [2 x i8] c"\0A\00", align 1

declare float @printf(i8* %x)

define void @main() {
entry:
  %0 = alloca i8, align 1
  store i8 73, i8* %0, align 1
  br label %loop

loop:                                             ; preds = %loop, %entry
  %ind = phi i32 [ 0, %entry ], [ %nextind, %loop ]
  %1 = call float @printf(i8* nonnull %0)
  %putchar = call i32 @putchar(i32 10)
  %nextind = add i32 %ind, 1
  %2 = sext i32 %nextind to i64
  %3 = getelementptr inbounds [14 x i8], [14 x i8]* @strtmp, i64 0, i64 %2
  %load1 = load i8, i8* %3, align 1
  store i8 %load1, i8* %0, align 1
  %loopcond = icmp eq i32 %nextind, 14
  br i1 %loopcond, label %exitloop, label %loop

exitloop:                                         ; preds = %loop
  ret void
}

; Function Attrs: nounwind
declare i32 @putchar(i32) #0

attributes #0 = { nounwind }
