// https://tomassetti.me/antlr-mega-tutorial/

grammar Hdl;

fragment LETTER: [A-Za-z];
fragment DIGIT: [0-9];

SPACE: [ \n\r\t]+ -> skip;

CHIP_LABEL: 'CHIP';
PARTS_LABEL: 'PARTS';
IN_LABEL: 'IN';
OUT_LABEL: 'OUT';

ID: LETTER (LETTER | DIGIT)*;
NUMBER: DIGIT+;
SEMICOLON: ';';

chips: chip+ EOF;
chip: CHIP_LABEL ID '{' inputs outputs PARTS_LABEL ':' parts '}';

inputs: IN_LABEL pinDeclarations SEMICOLON;
outputs: OUT_LABEL pinDeclarations SEMICOLON;
pinDeclarations: pinDeclaration (',' pinDeclaration)*;
pinDeclaration: ID ('[' NUMBER ']')?;

parts: partDeclaration+;
partDeclaration: ID '(' connections ')' SEMICOLON;
connections: connection (',' connection)*;
connection: slice '=' slice;
slice: ID ('[' NUMBER ('..' NUMBER)? ']')?;
