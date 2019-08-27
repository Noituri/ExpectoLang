; ModuleID = 'expectoroot'
source_filename = "expectoroot"

@strtmp = private unnamed_addr constant [6 x i8] c"true\0A\00", align 1
@strtmp.1 = private unnamed_addr constant [7 x i8] c"false\0A\00", align 1
@strtmp.2 = private unnamed_addr constant [3 x i8] c"aa\00", align 1
@strtmp.3 = private unnamed_addr constant [3 x i8] c"bb\00", align 1
@strtmp.4 = private unnamed_addr constant [4 x i8] c"Hi\0A\00", align 1
@str = private unnamed_addr constant [6 x i8] c"false\00", align 1
@str.5 = private unnamed_addr constant [3 x i8] c"Hi\00", align 1
@str.6 = private unnamed_addr constant [5 x i8] c"true\00", align 1

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

define i1 @"binary$&&"(i1 %a, i1 %b) {
entry:
  %0 = xor i1 %a, %b
  %spec.select = select i1 %0, i1 false, i1 true
  ret i1 %spec.select
}

define i1 @"unary$!"(i1 %a) {
entry:
  %spec.select = select i1 %a, i1 false, i1 true
  ret i1 %spec.select
}

define void @main() {
entry:
  %0 = call i1 @"unary$!"(i1 true)
  %1 = call i1 @"unary$!"(i1 %0)
  %2 = call i1 @"unary$!"(i1 false)
  %3 = call i1 @"unary$!"(i1 false)
  %4 = call i1 @"binary$&&"(i1 %3, i1 true)
  br i1 %4, label %then, label %else

then:                                             ; preds = %entry
  %puts2 = call i32 @puts(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @str.6, i64 0, i64 0))
  br label %exit

else:                                             ; preds = %entry
  %puts = call i32 @puts(i8* getelementptr inbounds ([6 x i8], [6 x i8]* @str, i64 0, i64 0))
  br label %exit

exit:                                             ; preds = %else, %then
  %5 = call double @add(double 1.000000e+00, double 1.000000e-01)
  %puts1 = call i32 @puts(i8* getelementptr inbounds ([3 x i8], [3 x i8]* @str.5, i64 0, i64 0))
  ret void
}

define double @add(double %a, double %b) {
entry:
  %addtmp = fadd double %a, %b
  %addtmp1 = fadd double %addtmp, 2.000000e+00
  ret double %addtmp1
}

; Function Attrs: nounwind
declare i32 @puts(i8* nocapture readonly) #0

attributes #0 = { nounwind }
