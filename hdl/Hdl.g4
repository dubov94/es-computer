grammar Hdl;

SPACE: (' ' | '\n')+ -> skip;
CHIP: 'CHIP';
LEFT_BRACE: '{';
RIGHT_BRACE: '}';
ID: LETTER (LETTER | DIGIT)*;

fragment LETTER: [A-Za-z];
fragment DIGIT: [0-9];

chip: CHIP ID LEFT_BRACE RIGHT_BRACE;
