@256
D=A
@SP
M=D
@TOP_LEVEL$ret.1
D=A
@SP
A=M
M=D
@SP
M=M+1
@LCL
D=M
@SP
A=M
M=D
@SP
M=M+1
@ARG
D=M
@SP
A=M
M=D
@SP
M=M+1
@THIS
D=M
@SP
A=M
M=D
@SP
M=M+1
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1
@5
D=A
@SP
D=M-D
@ARG
M=D
@SP
D=M
@LCL
M=D
@Sys.init
0;JMP
(TOP_LEVEL$ret.1)
// function Main.fibonacci 0
(Main.fibonacci)
// 	push argument 0
@ARG
D=M
@0
A=A+D
D=M
@SP
A=M
M=D
@SP
M=M+1
// 	push constant 2
@2
D=A
@SP
A=M
M=D
@SP
M=M+1
// 	lt                     
@SP
M=M-1
A=M
D=M
A=A-1
D=M-D
M=-1
@END_JLT_0
D;JLT
@SP
A=M-1
M=0
(END_JLT_0)
// 	if-goto N_LT_2        
@SP
M=M-1
@SP
A=M
D=M
@N_LT_2
D;JNE
// 	goto N_GE_2
@N_GE_2
0;JMP
// label N_LT_2               // if n < 2 returns n
(N_LT_2)
// 	push argument 0        
@ARG
D=M
@0
A=A+D
D=M
@SP
A=M
M=D
@SP
M=M+1
// 	return
@LCL
D=M
@R14
M=D
@5
A=D-A
D=M
@R15
M=D
@ARG
D=M
@0
D=D+A
@R13
M=D
@SP
M=M-1
@SP
A=M
D=M
@R13
A=M
M=D
@ARG
D=M+1
@SP
M=D
@R14
A=M-1
D=M
@THAT
M=D
@R14
D=M
@2
A=D-A
D=M
@THIS
M=D
@R14
D=M
@3
A=D-A
D=M
@ARG
M=D
@R14
D=M
@4
A=D-A
D=M
@LCL
M=D
@R15
A=M
0;JMP
// label N_GE_2               // if n >= 2 returns fib(n - 2) + fib(n - 1)
(N_GE_2)
// 	push argument 0
@ARG
D=M
@0
A=A+D
D=M
@SP
A=M
M=D
@SP
M=M+1
// 	push constant 2
@2
D=A
@SP
A=M
M=D
@SP
M=M+1
// 	sub
@SP
M=M-1
A=M
D=M
A=A-1
M=M-D
// 	call Main.fibonacci 1  // computes fib(n - 2)
@Main.fibonacci$ret.2
D=A
@SP
A=M
M=D
@SP
M=M+1
@LCL
D=M
@SP
A=M
M=D
@SP
M=M+1
@ARG
D=M
@SP
A=M
M=D
@SP
M=M+1
@THIS
D=M
@SP
A=M
M=D
@SP
M=M+1
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1
@6
D=A
@SP
D=M-D
@ARG
M=D
@SP
D=M
@LCL
M=D
@Main.fibonacci
0;JMP
(Main.fibonacci$ret.2)
// 	push argument 0
@ARG
D=M
@0
A=A+D
D=M
@SP
A=M
M=D
@SP
M=M+1
// 	push constant 1
@1
D=A
@SP
A=M
M=D
@SP
M=M+1
// 	sub
@SP
M=M-1
A=M
D=M
A=A-1
M=M-D
// 	call Main.fibonacci 1  // computes fib(n - 1)
@Main.fibonacci$ret.3
D=A
@SP
A=M
M=D
@SP
M=M+1
@LCL
D=M
@SP
A=M
M=D
@SP
M=M+1
@ARG
D=M
@SP
A=M
M=D
@SP
M=M+1
@THIS
D=M
@SP
A=M
M=D
@SP
M=M+1
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1
@6
D=A
@SP
D=M-D
@ARG
M=D
@SP
D=M
@LCL
M=D
@Main.fibonacci
0;JMP
(Main.fibonacci$ret.3)
// 	add                    // returns fib(n - 1) + fib(n - 2)
@SP
M=M-1
A=M
D=M
A=A-1
M=M+D
// 	return
@LCL
D=M
@R14
M=D
@5
A=D-A
D=M
@R15
M=D
@ARG
D=M
@0
D=D+A
@R13
M=D
@SP
M=M-1
@SP
A=M
D=M
@R13
A=M
M=D
@ARG
D=M+1
@SP
M=D
@R14
A=M-1
D=M
@THAT
M=D
@R14
D=M
@2
A=D-A
D=M
@THIS
M=D
@R14
D=M
@3
A=D-A
D=M
@ARG
M=D
@R14
D=M
@4
A=D-A
D=M
@LCL
M=D
@R15
A=M
0;JMP
// function Sys.init 0
(Sys.init)
// 	push constant 4
@4
D=A
@SP
A=M
M=D
@SP
M=M+1
// 	call Main.fibonacci 1   // computes the 4'th fibonacci element
@Sys.init$ret.4
D=A
@SP
A=M
M=D
@SP
M=M+1
@LCL
D=M
@SP
A=M
M=D
@SP
M=M+1
@ARG
D=M
@SP
A=M
M=D
@SP
M=M+1
@THIS
D=M
@SP
A=M
M=D
@SP
M=M+1
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1
@6
D=A
@SP
D=M-D
@ARG
M=D
@SP
D=M
@LCL
M=D
@Main.fibonacci
0;JMP
(Sys.init$ret.4)
// label END  
(END)
// 	goto END                // loops infinitely
@END
0;JMP
