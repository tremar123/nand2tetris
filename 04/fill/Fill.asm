// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen
// by writing 'black' in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen by writing
// 'white' in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// if KBD == 0 set status to white else black
(LOOP)
    @i
    M=0
    @KBD
    D=M
    @SETWHITE
    D;JEQ
    @status
    M=-1
    @SETSCREEN
    0;JMP

(SETWHITE)
    @status
    M=0
    @SETSCREEN
    0;JMP

(SETSCREEN)
    @i
    D=M
    @8192
    D=D-A
    @LOOP
    D;JEQ

    // get current addr and set value to status
    @i
    D=M
    @SCREEN
    D=D+A
    @addr
    M=D
    @status
    D=M
    @addr
    A=M
    M=D
    @i
    M=M+1

    @SETSCREEN
    0;JMP
