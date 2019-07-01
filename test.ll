; ModuleID = 'test.c'
source_filename = "test.c"
target datalayout = "e-m:e-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@__const.main.s1 = private unnamed_addr constant [7 x i8] c"Hello \00", align 1
@__const.main.s2 = private unnamed_addr constant [6 x i8] c"World\00", align 1

; Function Attrs: noinline nounwind optnone uwtable
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

; Function Attrs: nounwind readonly
declare dso_local i64 @strlen(i8*) #1

; Function Attrs: nounwind
declare dso_local noalias i8* @malloc(i64) #2

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 {
  %1 = alloca [7 x i8], align 1
  %2 = alloca [6 x i8], align 1
  %3 = alloca i8*, align 8
  %4 = bitcast [7 x i8]* %1 to i8*
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 1 %4, i8* align 1 getelementptr inbounds ([7 x i8], [7 x i8]* @__const.main.s1, i32 0, i32 0), i64 7, i1 false)
  %5 = bitcast [6 x i8]* %2 to i8*
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 1 %5, i8* align 1 getelementptr inbounds ([6 x i8], [6 x i8]* @__const.main.s2, i32 0, i32 0), i64 6, i1 false)
  %6 = getelementptr inbounds [7 x i8], [7 x i8]* %1, i32 0, i32 0
  %7 = getelementptr inbounds [6 x i8], [6 x i8]* %2, i32 0, i32 0
  %8 = call i8* @stringMerge(i8* %6, i8* %7)
  store i8* %8, i8** %3, align 8
  %9 = load i8*, i8** %3, align 8
  %10 = call i32 @puts(i8* %9)
  %11 = load i8*, i8** %3, align 8
  call void @free(i8* %11) #6
  ret i32 0
}

; Function Attrs: argmemonly nounwind
declare void @llvm.memcpy.p0i8.p0i8.i64(i8* nocapture writeonly, i8* nocapture readonly, i64, i1) #3

declare dso_local i32 @puts(i8*) #4

; Function Attrs: nounwind
declare dso_local void @free(i8*) #2

attributes #0 = { noinline nounwind optnone uwtable "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "min-legal-vector-width"="0" "no-frame-pointer-elim"="true" "no-frame-pointer-elim-non-leaf" "no-infs-fp-math"="false" "no-jump-tables"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+fxsr,+mmx,+sse,+sse2,+x87" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #1 = { nounwind readonly "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "no-frame-pointer-elim"="true" "no-frame-pointer-elim-non-leaf" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+fxsr,+mmx,+sse,+sse2,+x87" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #2 = { nounwind "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "no-frame-pointer-elim"="true" "no-frame-pointer-elim-non-leaf" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+fxsr,+mmx,+sse,+sse2,+x87" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #3 = { argmemonly nounwind }
attributes #4 = { "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "no-frame-pointer-elim"="true" "no-frame-pointer-elim-non-leaf" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+fxsr,+mmx,+sse,+sse2,+x87" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #5 = { nounwind readonly }
attributes #6 = { nounwind }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{!"clang version 8.0.0 (tags/RELEASE_800/final)"}
