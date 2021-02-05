/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2014 Bart Kiers
 * Copyright (c) 2019 Robert Einhorn
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation
 * files (the "Software"), to deal in the Software without
 * restriction, including without limitation the rights to use,
 * copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following
 * conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
 * OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
 * HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
 * WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 * OTHER DEALINGS IN THE SOFTWARE.
 *
 * Project      : Python3-parser; an ANTLR4 grammar for Python 3
 *                https://github.com/bkiers/Python3-parser
 * Developed by : Bart Kiers, bart@big-o.nl
 *
 * Project      : an ANTLR4 grammar for Python 3 without actions
 *                https://github.com/antlr/grammars-v4/tree/master/python/python3-without-actions
 * Developed by : Robert Einhorn, robert.einhorn.hu@gmail.com
 */

// This file was modified to only parse import statements. The original file was
// extracted from the ANTLR4 grammar for Python 3 without actions.

grammar Python3;

/*
 * parser rules
 */

fileInput : ( NEWLINE | importStmt )*;

importStmt : importName | importFrom;

importName : 'import' dottedAsNames;

importFrom
 : 'from' ( ( '.' | '...' )* dottedName
        | ( '.' | '...' )+
        )
   'import' ( '*'
          | '(' importAsNames ')'
          | importAsNames
          )
 ;

importAsName : NAME ( 'as' NAME )?;

dottedAsName : dottedName ( 'as' NAME )?;

importAsNames : importAsName ( ',' importAsName )* ','?;

dottedAsNames : dottedAsName ( ',' dottedAsName )*;

dottedName : NAME ( '.' NAME )*;

/*
 * lexer rules
 */

NEWLINE : ( '\r'? '\n' | '\r' | '\f' ) SPACES?;

NAME : ID_START ID_CONTINUE*;

DOT : '.';

SKIP_ : ( SPACES | COMMENT | LINE_JOINING ) -> skip;

UNKNOWN_CHAR : .;

/*
 * fragments
 */

fragment SPACES : [ \t]+;

fragment COMMENT : '#' ~[\r\n\f]*;

fragment LINE_JOINING : '\\' SPACES? ( '\r'? '\n' | '\r' | '\f' );

fragment OTHER_ID_START : [\u2118\u212E\u309B\u309C];

fragment ID_START : '_' | [\p{Letter}\p{Letter_Number}] | OTHER_ID_START;

fragment ID_CONTINUE
 : ID_START
   | [\p{Nonspacing_Mark}\p{Spacing_Mark}\p{Decimal_Number}\p{Connector_Punctuation}\p{Format}]
 ;
