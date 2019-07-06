; ModuleID = 'expectoroot'
source_filename = "expectoroot"

@strtmp = private unnamed_addr constant [11 x i8] c"that's it\0A\00", align 1
@strtmp.1 = private unnamed_addr constant [2 x i8] c"t\00", align 1
@strtmp.2 = private unnamed_addr constant [5 x i8] c"t123\00", align 1
@strtmp.3 = private unnamed_addr constant [5 x i8] c"YAY\0A\00", align 1
@strtmp.4 = private unnamed_addr constant [5 x i8] c"NAY\0A\00", align 1

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
  br i1 icmp eq (i8 ptrtoint ([2 x i8]* @strtmp.1 to i8), i8 ptrtoint ([5 x i8]* @strtmp.2 to i8)), label %then, label %else

then:                                             ; preds = %entry
  %0 = call float @printf(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @strtmp.3, i32 0, i32 0))
  br label %exit

else:                                             ; preds = %entry
  %1 = call float @printf(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @strtmp.4, i32 0, i32 0))
  br label %exit

exit:                                             ; preds = %else, %then
  call void @call()
  ret void
}
