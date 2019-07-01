; ModuleID = 'expectoroot'
source_filename = "expectoroot"

@strtmp = private unnamed_addr constant [15 x i8] c"Hello Expecto\0A\00", align 1
@str = private unnamed_addr constant [14 x i8] c"Hello Expecto\00", align 1

declare i32 @printf(i8* %x)

define float @add(float %a, float %b) {
entry:
  %addtmp = fadd float %a, %b
  ret float %addtmp
}

define float @main() {
entry:
  %puts = call i32 @puts(i8* getelementptr inbounds ([14 x i8], [14 x i8]* @str, i64 0, i64 0))
  %calltmp1 = call float @add(float 4.000000e+00, float 1.000000e+00)
  %addtmp = fadd float %calltmp1, 4.000000e+00
  %multmo = fmul float %addtmp, 2.000000e+00
  %addtmp2 = fadd float %multmo, 3.000000e+00
  ret float %addtmp2
}

; Function Attrs: nounwind
declare i32 @puts(i8* nocapture readonly) #0

attributes #0 = { nounwind }
