; ModuleID = 'expectoroot'
source_filename = "expectoroot"

@strtmp = private unnamed_addr constant [11 x i8] c"that's it\0A\00", align 1
@strtmp.1 = private unnamed_addr constant [14 x i8] c"IterThroughMe\00", align 1
@strtmp.2 = private unnamed_addr constant [2 x i8] c"\0A\00", align 1

declare float @printf(i8* %x)

define void @call() {
entry:
  %0 = call float @printf(i8* getelementptr inbounds ([11 x i8], [11 x i8]* @strtmp, i32 0, i32 0))
  ret void
}

define float @sum(float %a, float %b) {
entry:
  %addtmp = fadd float %a, %b
  ret float %addtmp
}

define void @initer(i8* %x) {
entry:
  ret void
}

define void @main() {
entry:
  br label %loop

loop:                                             ; preds = %loop, %entry
  %ind = phi i32 [ 0, %entry ], [ %nextind, %loop ]
  %0 = alloca i8*
  %1 = alloca i8
  store i8* getelementptr inbounds ([14 x i8], [14 x i8]* @strtmp.1, i32 0, i32 0), i8** %0
  %load = load i8*, i8** %0
  %2 = getelementptr inbounds i8, i8* %load, i32 0
  %load1 = load i8, i8* %2
  store i8 %load1, i8* %1
  call void @initer(i8* %1)
  %3 = call float @printf(i8* %1)
  %4 = call float @printf(i8* getelementptr inbounds ([2 x i8], [2 x i8]* @strtmp.2, i32 0, i32 0))
  %nextind = add i32 %ind, 1
  %loopcond = icmp ne i32 14, %nextind
  br i1 %loopcond, label %loop, label %exitloop

exitloop:                                         ; preds = %loop
  call void @call()
  ret void
}
