; ModuleID = 'expectoroot'
source_filename = "expectoroot"

@strtmp = private unnamed_addr constant [5 x i8] c"test\00", align 1
@strtmp.1 = private unnamed_addr constant [5 x i8] c"lol\0A\00", align 1

declare i32 @printf(i8* %x)
; Function Attrs: nounwind readonly
declare dso_local i64 @strlen(i8*) #1

; Function Attrs: nounwind
declare dso_local noalias i8* @malloc(i64) #2
define dso_local i8* @stringMerge(i8*, i8*) #0 {
  %3 = alloca i8*, align 8
  %4 = alloca i8*, align 8
  %5 = alloca i64, align 8
  %6 = alloca i8*, align 8
  %7 = alloca i8*, align 8
  store i8* %0, i8** %3, align 8
  store i8* %1, i8** %4, align 8
  %8 = load i8*, i8** %3, align 8
  %9 = call i64 @strlen(i8* %8) #5
  store i64 %9, i64* %5, align 8
  %10 = load i64, i64* %5, align 8
  %11 = load i8*, i8** %4, align 8
  %12 = call i64 @strlen(i8* %11) #5
  %13 = add i64 %10, %12
  %14 = add i64 %13, 1
  %15 = call noalias i8* @malloc(i64 %14) #6
  store i8* %15, i8** %6, align 8
  %16 = load i8*, i8** %6, align 8
  %17 = icmp ne i8* %16, null
  br i1 %17, label %18, label %41

; <label>:18:                                     ; preds = %2
  %19 = load i8*, i8** %6, align 8
  store i8* %19, i8** %7, align 8
  br label %20

; <label>:20:                                     ; preds = %26, %18
  %21 = load i8*, i8** %3, align 8
  %22 = getelementptr inbounds i8, i8* %21, i32 1
  store i8* %22, i8** %3, align 8
  %23 = load i8, i8* %21, align 1
  %24 = load i8*, i8** %7, align 8
  store i8 %23, i8* %24, align 1
  %25 = icmp ne i8 %23, 0
  br i1 %25, label %26, label %29

; <label>:26:                                     ; preds = %20
  %27 = load i8*, i8** %7, align 8
  %28 = getelementptr inbounds i8, i8* %27, i32 1
  store i8* %28, i8** %7, align 8
  br label %20

; <label>:29:                                     ; preds = %20
  br label %30

; <label>:30:                                     ; preds = %35, %29
  %31 = load i8*, i8** %4, align 8
  %32 = getelementptr inbounds i8, i8* %31, i32 1
  store i8* %32, i8** %4, align 8
  %33 = load i8, i8* %31, align 1
  %34 = load i8*, i8** %7, align 8
  store i8 %33, i8* %34, align 1
  br label %35

; <label>:35:                                     ; preds = %30
  %36 = load i8*, i8** %7, align 8
  %37 = getelementptr inbounds i8, i8* %36, i32 1
  store i8* %37, i8** %7, align 8
  %38 = load i8, i8* %36, align 1
  %39 = icmp ne i8 %38, 0
  br i1 %39, label %30, label %40

; <label>:40:                                     ; preds = %35
  br label %41

; <label>:41:                                     ; preds = %40, %2
  %42 = load i8*, i8** %6, align 8
  ret i8* %42
}

define float @main() {
entry:
  %calltmp = call i8* @stringMerge(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @strtmp, i64 0, i64 0), i8* getelementptr inbounds ([5 x i8], [5 x i8]* @strtmp.1, i64 0, i64 0))
  %calltmp1 = call i32 @printf(i8* %calltmp)
  ret float 9.000000e+00
}
