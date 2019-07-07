; ModuleID = 'expectoroot'
source_filename = "expectoroot"

@strtmp = private unnamed_addr constant [11 x i8] c"that's it\0A\00", align 1
@strtmp.1 = private unnamed_addr constant [10 x i8] c"SOMEarray\00", align 1
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

define void @main() {
entry:
  br label %loop

loop:                                             ; preds = %loop, %entry
  %ind = phi i32 [ 0, %entry ], [ %nextind, %loop ]
  %0 = call float @printf(i8* getelementptr inbounds ([10 x i8], [10 x i8]* @strtmp.1, i32 0, i32 0))
  %1 = call float @printf(i8* getelementptr inbounds ([2 x i8], [2 x i8]* @strtmp.2, i32 0, i32 0))
  %nextind = add i32 %ind, 1
  %loopcond = icmp ne i32 5, %nextind
  br i1 %loopcond, label %loop, label %exitloop

exitloop:                                         ; preds = %loop
  call void @call()
  ret void
}
