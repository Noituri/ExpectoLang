; ModuleID = 'expectoroot'
source_filename = "expectoroot"

@strtmp = private unnamed_addr constant [5 x i8] c"sss\0A\00", align 1
@strtmp.1 = private unnamed_addr constant [6 x i8] c"ss2s\0A\00", align 1
@strtmp.2 = private unnamed_addr constant [5 x i8] c"sss\0A\00", align 1
@strtmp.3 = private unnamed_addr constant [6 x i8] c"ss2s\0A\00", align 1
@strtmp.4 = private unnamed_addr constant [5 x i8] c"sss\0A\00", align 1
@strtmp.5 = private unnamed_addr constant [5 x i8] c"sss\0A\00", align 1
@strtmp.6 = private unnamed_addr constant [6 x i8] c"ss2s\0A\00", align 1
@strtmp.7 = private unnamed_addr constant [5 x i8] c"sss\0A\00", align 1
@strtmp.8 = private unnamed_addr constant [5 x i8] c"sss\0A\00", align 1
@strtmp.9 = private unnamed_addr constant [9 x i8] c"I'm done\00", align 1
@str = private unnamed_addr constant [4 x i8] c"sss\00", align 1
@str.10 = private unnamed_addr constant [5 x i8] c"ss2s\00", align 1
@str.11 = private unnamed_addr constant [4 x i8] c"sss\00", align 1
@str.12 = private unnamed_addr constant [5 x i8] c"ss2s\00", align 1

declare float @printf(i8* %x)

define float @main() {
entry:
  %puts = call i32 @puts(i8* getelementptr inbounds ([4 x i8], [4 x i8]* @str, i64 0, i64 0))
  %puts16 = call i32 @puts(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @str.10, i64 0, i64 0))
  %puts17 = call i32 @puts(i8* getelementptr inbounds ([4 x i8], [4 x i8]* @str.11, i64 0, i64 0))
  %puts18 = call i32 @puts(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @str.12, i64 0, i64 0))
  %calltmp15 = call float @printf(i8* getelementptr inbounds ([9 x i8], [9 x i8]* @strtmp.9, i64 0, i64 0))
  ret float 9.000000e+00
}

; Function Attrs: nounwind
declare i32 @puts(i8* nocapture readonly) #0

attributes #0 = { nounwind }
