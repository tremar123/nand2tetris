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
// function Class1.set 0
(Class1.set)
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
// 	pop static 0
@SP
M=M-1
@SP
A=M
D=M
@Class1.0
M=D
// 	push argument 1
@ARG
D=M
@1
A=A+D
D=M
@SP
A=M
M=D
@SP
M=M+1
// 	pop static 1
@SP
M=M-1
@SP
A=M
D=M
@Class1.1
M=D
// 	push constant 0
@0
D=A
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
// function Class1.get 0
(Class1.get)
// 	push static 0
@Class1.0
D=M
@SP
A=M
M=D
@SP
M=M+1
// 	push static 1
@Class1.1
D=M
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
// function Class2.set 0
(Class2.set)
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
// 	pop static 0
@SP
M=M-1
@SP
A=M
D=M
@Class2.0
M=D
// 	push argument 1
@ARG
D=M
@1
A=A+D
D=M
@SP
A=M
M=D
@SP
M=M+1
// 	pop static 1
@SP
M=M-1
@SP
A=M
D=M
@Class2.1
M=D
// 	push constant 0
@0
D=A
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
// function Class2.get 0
(Class2.get)
// 	push static 0
@Class2.0
D=M
@SP
A=M
M=D
@SP
M=M+1
// 	push static 1
@Class2.1
D=M
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
// 	push constant 6
@6
D=A
@SP
A=M
M=D
@SP
M=M+1
// 	push constant 8
@8
D=A
@SP
A=M
M=D
@SP
M=M+1
// 	call Class1.set 2
@Sys.init$ret.2
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
@7
D=A
@SP
D=M-D
@ARG
M=D
@SP
D=M
@LCL
M=D
@Class1.set
0;JMP
(Sys.init$ret.2)
// 	pop temp 0 // dumps the return value
@SP
M=M-1
@SP
A=M
D=M
@5
M=D
// 	push constant 23
@23
D=A
@SP
A=M
M=D
@SP
M=M+1
// 	push constant 15
@15
D=A
@SP
A=M
M=D
@SP
M=M+1
// 	call Class2.set 2
@Sys.init$ret.3
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
@7
D=A
@SP
D=M-D
@ARG
M=D
@SP
D=M
@LCL
M=D
@Class2.set
0;JMP
(Sys.init$ret.3)
// 	pop temp 0 // dumps the return value
@SP
M=M-1
@SP
A=M
D=M
@5
M=D
// 	call Class1.get 0
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
@Class1.get
0;JMP
(Sys.init$ret.4)
// 	call Class2.get 0
@Sys.init$ret.5
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
@Class2.get
0;JMP
(Sys.init$ret.5)
// label END
(END)
// 	goto END
@END
0;JMP