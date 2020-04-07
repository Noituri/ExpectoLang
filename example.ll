; ModuleID = 'novumroot'
source_filename = "novumroot"

@strtmp = private unnamed_addr constant [4 x i8] c"%s\0A\00", align 1
@strtmp.1 = private unnamed_addr constant [13 x i8] c"Hello world!\00", align 1
@strtmp.2 = private unnamed_addr constant [15 x i8] c"Hello world 2!\00", align 1

declare i32 @printf(i8* %msg, i8* %format)

define i1 @"binary_&&"(i1 %val_a, i1 %val_b) {
entry:
  %0 = xor i1 %val_a, %val_b
  %spec.select = select i1 %0, i1 false, i1 true
  ret i1 %spec.select
}

define i1 @"unary_!"(i1 %a) {
entry:
  %spec.select = select i1 %a, i1 false, i1 true
  ret i1 %spec.select
}

define void @writeln(i8* %msg) {
entry:
  %0 = call i32 @printf(i8* getelementptr inbounds ([4 x i8], [4 x i8]* @strtmp, i64 0, i64 0), i8* %msg)
  ret void
}

define void @main() {
entry:
  call void @writeln(i8* getelementptr inbounds ([13 x i8], [13 x i8]* @strtmp.1, i64 0, i64 0))
  call void @writeln(i8* getelementptr inbounds ([15 x i8], [15 x i8]* @strtmp.2, i64 0, i64 0))
  %0 = call i1 @"binary_&&"(i1 true, i1 true)
  %1 = call i1 @"unary_!"(i1 %0)
  %2 = call i1 @"binary_&&"(i1 false, i1 %1)
  ret void
}
