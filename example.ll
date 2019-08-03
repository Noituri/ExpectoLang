; ModuleID = 'expectoroot'
source_filename = "expectoroot"

@strtmp = private unnamed_addr constant [11 x i8] c"Here I am\0A\00", align 1
@str = private unnamed_addr constant [10 x i8] c"Here I am\00", align 1
@strtmp.1 = private unnamed_addr constant [14 x i8] c"IterThroughMe\00", align 1
@strtmp.2 = private unnamed_addr constant [2 x i8] c"\0A\00", align 1
@strtmp.3 = private unnamed_addr constant [4 x i8] c"OwO\00", align 1

declare float @printf(i8* %x)

define void @test_loop(i1 %use_loop) {
entry:
  br i1 %use_loop, label %loop.preheader, label %exitloop

loop.preheader:                                   ; preds = %entry
  br label %loop

loop:                                             ; preds = %loop.preheader, %loop
  %puts = call i32 @puts(i8* getelementptr inbounds ([10 x i8], [10 x i8]* @str, i64 0, i64 0))
  br label %loop

exitloop:                                         ; preds = %entry
  ret void
}

; Function Attrs: nounwind
declare i32 @puts(i8* nocapture readonly) #0

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
  %3 = getelementptr inbounds [14 x i8], [14 x i8]* @strtmp.1, i64 0, i64 %2
  %load1 = load i8, i8* %3, align 1
  store i8 %load1, i8* %0, align 1
  %loopcond = icmp eq i32 %nextind, 14
  br i1 %loopcond, label %exitloop3, label %loop

exitloop3:                                        ; preds = %loop
  call void @test_loop(i1 true)
  ret void
}

; Function Attrs: nounwind
declare i32 @putchar(i32) #0

attributes #0 = { nounwind }
