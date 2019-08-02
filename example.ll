; ModuleID = 'expectoroot'
source_filename = "expectoroot"

@strtmp = private unnamed_addr constant [11 x i8] c"that's it\0A\00", align 1
@strtmp.1 = private unnamed_addr constant [14 x i8] c"IterThroughMe\00", align 1
@strtmp.2 = private unnamed_addr constant [2 x i8] c"\0A\00", align 1
@strtmp.3 = private unnamed_addr constant [11 x i8] c"HERE I AM\0A\00", align 1

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

define void @main() {
entry:
  %0 = alloca i8
  %1 = alloca i8*
  store i8* getelementptr inbounds ([14 x i8], [14 x i8]* @strtmp.1, i32 0, i32 0), i8** %1
  %load = load i8*, i8** %1
  %2 = getelementptr inbounds i8, i8* %load, i32 0
  %load1 = load i8, i8* %2
  store i8 %load1, i8* %0
  br label %loop

loop:                                             ; preds = %exit, %entry
  %ind = phi i32 [ 0, %entry ], [ %nextind, %exit ]
  %3 = call float @printf(i8* %0)
  %4 = call float @printf(i8* getelementptr inbounds ([2 x i8], [2 x i8]* @strtmp.2, i32 0, i32 0))
  br i1 true, label %then, label %else

then:                                             ; preds = %loop
  %5 = call float @printf(i8* getelementptr inbounds ([11 x i8], [11 x i8]* @strtmp.3, i32 0, i32 0))
  br label %exit

else:                                             ; preds = %loop
  br label %exit

exit:                                             ; preds = %else, %then
  %nextind = add i32 %ind, 1
  %6 = getelementptr inbounds i8, i8* %load, i32 %nextind
  %load2 = load i8, i8* %6
  store i8 %load2, i8* %0
  %loopcond = icmp ne i32 14, %nextind
  br i1 %loopcond, label %loop, label %exitloop

exitloop:                                         ; preds = %exit
  call void @call()
  ret void
}
