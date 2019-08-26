; ModuleID = 'expectoroot'
source_filename = "expectoroot"

@strtmp = private unnamed_addr constant [15 x i8] c"IterThroughMe\0A\00", align 1
@strtmp.1 = private unnamed_addr constant [6 x i8] c"true\0A\00", align 1
@strtmp.2 = private unnamed_addr constant [7 x i8] c"false\0A\00", align 1
@strtmp.3 = private unnamed_addr constant [4 x i8] c"Hi\0A\00", align 1
@str = private unnamed_addr constant [6 x i8] c"false\00", align 1
@str.4 = private unnamed_addr constant [3 x i8] c"Hi\00", align 1
@str.5 = private unnamed_addr constant [5 x i8] c"true\00", align 1

declare i32 @printf(i8* %x)

declare void @exit(i32 %x)

define i1 @"binary$?"(double %a, double %b) {
entry:
  %cmptmp = fcmp oeq double %a, 1.000000e+00
  %cmptmp1 = fcmp oeq double %b, 0.000000e+00
  %. = select i1 %cmptmp1, i1 true, i1 false
  %merge = select i1 %cmptmp, i1 %., i1 false
  ret i1 %merge
}

define void @main() {
entry:
  %0 = alloca i8, align 1
  store i8 73, i8* %0, align 1
  br label %loop

loop:                                             ; preds = %loop, %entry
  %ind = phi i32 [ 0, %entry ], [ %nextind, %loop ]
  %1 = call i32 @printf(i8* nonnull %0)
  %nextind = add i32 %ind, 1
  %2 = sext i32 %nextind to i64
  %3 = getelementptr inbounds [15 x i8], [15 x i8]* @strtmp, i64 0, i64 %2
  %load1 = load i8, i8* %3, align 1
  store i8 %load1, i8* %0, align 1
  %loopcond = icmp eq i32 %nextind, 15
  br i1 %loopcond, label %exitloop, label %loop

exitloop:                                         ; preds = %loop
  %4 = call i1 @"binary$?"(double 1.000000e+00, double 0.000000e+00)
  br i1 %4, label %then, label %else

then:                                             ; preds = %exitloop
  %puts3 = call i32 @puts(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @str.5, i64 0, i64 0))
  br label %exit

else:                                             ; preds = %exitloop
  %puts = call i32 @puts(i8* getelementptr inbounds ([6 x i8], [6 x i8]* @str, i64 0, i64 0))
  br label %exit

exit:                                             ; preds = %else, %then
  %5 = call double @add(double 1.000000e+00, double 1.000000e-01)
  %puts2 = call i32 @puts(i8* getelementptr inbounds ([3 x i8], [3 x i8]* @str.4, i64 0, i64 0))
  ret void
}

define double @add(double %a, double %b) {
entry:
  %addtmp = fadd double %a, %b
  ret double %addtmp
}

; Function Attrs: nounwind
declare i32 @puts(i8* nocapture readonly) #0

attributes #0 = { nounwind }
