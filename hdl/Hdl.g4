// https://tomassetti.me/antlr-mega-tutorial/

grammar Hdl;

fragment LETTER: [A-Za-z];
fragment DIGIT: [0-9];

SPACE: [ \n]+ -> skip;
CHIP: 'CHIP';
LEFT_BRACE: '{';
RIGHT_BRACE: '}';
ID: LETTER (LETTER | DIGIT)*;

chip: CHIP id=ID LEFT_BRACE RIGHT_BRACE EOF;
