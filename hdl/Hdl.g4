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

inputs: IN_LABEL pin_declarations SEMICOLON;
outputs: OUT_LABEL pin_declarations SEMICOLON;
pin_declarations: pin_declaration (',' pin_declaration)*;
pin_declaration: ID ('[' NUMBER ']')?;

parts: part_declaration+;
part_declaration: ID '(' connections ')' SEMICOLON;
connections: connection (',' connection)*;
connection: slice '=' slice;
slice: ID ('[' NUMBER ('..' NUMBER)? ']')?;
