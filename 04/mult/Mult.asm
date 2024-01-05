// Multiply RAM[0] and RAM[1] and store result in RAM[2]
// i = RAM[1]
// while i != 0
// RAM[2] += RAM[0]
// i--

// set i to RAM[1]
@R2
M=0
@R1
D=M
@i
M=D

(LOOP)
    // check i == 0 goto end if yes
    @i
    D=M
    @END
    D;JEQ

    // RAM[2]+=RAM[0]
    @R2
    D=M
    @R0
    D=D+M
    @R2
    M=D

    // i--
    @i
    M=M-1

    @LOOP
    0;JMP

(END)
    @END
    0;JMP
