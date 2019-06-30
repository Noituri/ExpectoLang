; ModuleID = 'expectoroot'
source_filename = "expectoroot"

@strtmp = private unnamed_addr constant [15 x i8] c"Hello Expecto\0A\00", align 1

declare i32 @printf(i8* %x)

define float @main() {
entry:
  %calltmp = call i32 @printf(i8* getelementptr inbounds ([15 x i8], [15 x i8]* @strtmp, i32 0, i32 0))
  ret float 0.000000e+00
}
