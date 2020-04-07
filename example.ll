; ModuleID = 'novumroot'
source_filename = "novumroot"

@strtmp = private unnamed_addr constant [11 x i8] c"looping...\00", align 1
@strtmp.1 = private unnamed_addr constant [13 x i8] c"Hello world!\00", align 1
@strtmp.2 = private unnamed_addr constant [15 x i8] c"Hello world 2!\00", align 1
@strtmp.3 = private unnamed_addr constant [6 x i8] c"hello\00", align 1
@strtmp.4 = private unnamed_addr constant [4 x i8] c"%s\0A\00", align 1

define void @test_loop(i1 %should_loop) {
entry:
  %loopcond = icmp ne i1 %should_loop, false
  br i1 %loopcond, label %loop, label %exitloop

loop:                                             ; preds = %loop, %entry
  %ind = phi i32 [ 0, %entry ], [ %nextind, %loop ]
  call void @writeln(i8* getelementptr inbounds ([11 x i8], [11 x i8]* @strtmp, i32 0, i32 0))
  %nextind = add i32 %ind, 1
  %loopcond1 = icmp ne i1 %should_loop, false
  br i1 %loopcond1, label %loop, label %exitloop

exitloop:                                         ; preds = %loop, %entry
  ret void
}

define void @main() {
entry:
  call void @writeln(i8* getelementptr inbounds ([13 x i8], [13 x i8]* @strtmp.1, i32 0, i32 0))
  call void @writeln(i8* getelementptr inbounds ([15 x i8], [15 x i8]* @strtmp.2, i32 0, i32 0))
  %0 = alloca i8
  %load = load i8, i8* getelementptr inbounds ([6 x i8], [6 x i8]* @strtmp.3, i32 0, i32 0)
  store i8 %load, i8* %0
  br label %loop

loop:                                             ; preds = %loop, %entry
  %ind = phi i32 [ 0, %entry ], [ %nextind, %loop ]
  call void @writeln(i8* %0)
  %nextind = add i32 %ind, 1
  %1 = getelementptr inbounds i8, i8* getelementptr inbounds ([6 x i8], [6 x i8]* @strtmp.3, i32 0, i32 0), i32 %nextind
  %load1 = load i8, i8* %1
  store i8 %load1, i8* %0
  %loopcond = icmp ne i32 6, %nextind
  br i1 %loopcond, label %loop, label %exitloop

exitloop:                                         ; preds = %loop
  call void @test_loop(i1 false)
  %2 = call i1 @"binary_&&"(i1 true, i1 true)
  %3 = call i1 @"unary_!"(i1 %2)
  %4 = call i1 @"binary_&&"(i1 false, i1 %3)
  br i1 %4, label %then, label %else

then:                                             ; preds = %exitloop
  ret void

else:                                             ; preds = %exitloop
  br label %exit

exit:                                             ; preds = %else
  ret void
}

declare i32 @printf(i8* %msg, i8* %format)

define i1 @"binary_&&"(i1 %val_a, i1 %val_b) {
entry:
  %cmptmp = icmp eq i1 %val_a, %val_b
  br i1 %cmptmp, label %then, label %else

then:                                             ; preds = %entry
  ret i1 true

else:                                             ; preds = %entry
  br label %exit

exit:                                             ; preds = %else
  ret i1 false
}

define i1 @"unary_!"(i1 %a) {
entry:
  br i1 %a, label %then, label %else

then:                                             ; preds = %entry
  ret i1 false

else:                                             ; preds = %entry
  br label %exit

exit:                                             ; preds = %else
  ret i1 true
}

define void @writeln(i8* %msg) {
entry:
  %0 = call i32 @printf(i8* getelementptr inbounds ([4 x i8], [4 x i8]* @strtmp.4, i32 0, i32 0), i8* %msg)
  ret void
}
