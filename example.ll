; ModuleID = 'expectoroot'
source_filename = "expectoroot"

@strtmp = private unnamed_addr constant [14 x i8] c"IterThroughMe\00", align 1
@strtmp.1 = private unnamed_addr constant [2 x i8] c"\0A\00", align 1
@strtmp.2 = private unnamed_addr constant [4 x i8] c"OwO\00", align 1

declare float @printf(i8* %x)

define void @main() {
entry:
  %0 = alloca i8
  %load = load i8, i8* getelementptr inbounds ([14 x i8], [14 x i8]* @strtmp, i32 0, i32 0)
  store i8 %load, i8* %0
  br label %loop

loop:                                             ; preds = %loop, %entry
  %ind = phi i32 [ 0, %entry ], [ %nextind, %loop ]
  %1 = call float @printf(i8* %0)
  %2 = call float @printf(i8* getelementptr inbounds ([2 x i8], [2 x i8]* @strtmp.1, i32 0, i32 0))
  %nextind = add i32 %ind, 1
  %3 = getelementptr inbounds i8, i8* getelementptr inbounds ([14 x i8], [14 x i8]* @strtmp, i32 0, i32 0), i32 %nextind
  %load1 = load i8, i8* %3
  store i8 %load1, i8* %0
  %loopcond = icmp ne i32 14, %nextind
  br i1 %loopcond, label %loop, label %exitloop

exitloop:                                         ; preds = %loop
  br i1 false, label %loop2, label %exitloop3

loop2:                                            ; preds = %loop2, %exitloop
  %ind4 = phi i32 [ 0, %exitloop ], [ %nextind5, %loop2 ]
  %4 = call float @printf(i8* getelementptr inbounds ([4 x i8], [4 x i8]* @strtmp.2, i32 0, i32 0))
  %nextind5 = add i32 %ind4, 1
  br i1 false, label %loop2, label %exitloop3

exitloop3:                                        ; preds = %loop2, %exitloop
  ret void
}
