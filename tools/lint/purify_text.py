#!/usr/bin/env python3

r'''Purify text files.


1.  In each line, remove trailing whitespace characters.
    Whitespace characters are:
            \t \v \n \r space

2.  In each file, remove leading and trailing empty lines, i.e., the whole line
    only consists of whitespace characters.


3.  Add a UNIX new-line character \n to EOF if the EOF does not have one
    already.  The only exception to this rule is when the file is actually
    empty.


4.  Convert to UNIX line ending characters.
    This means to convert \r\n to \n.
'''

import argparse
import re



__EMPTYLINE_REGEX__ = re.compile(r'^\s*$')


def fix_lines(lines):
    '''Purify a list of lines.

    Args:
        line: The list of lines to process.

    Returns:
        The purified lines, as a list of strings.
    '''
    retval = []

    # Fix each line.
    for line in lines:
        line = fix_line(line)
        retval.append(line)

    # Remove leading and trailing empty lines.
    while retval and __EMPTYLINE_REGEX__.match(retval[0]):
        retval.pop(0)
    while retval and __EMPTYLINE_REGEX__.match(retval[-1]):
        retval.pop(-1)

    return retval


def fix_line(line):
    '''Purify one line.

    Args:
        line: String. The line to process.

    Returns:
        The purified line.
    '''
    return line.rstrip() + '\n'


def purify_text_files(filelist):
    '''Purify text files in-place.

    The detailed specs of the purification are described in the module
    docstring.

    Args:
        filelist: A list of file paths.
    '''
    for path in filelist:
        with open(path, 'r') as f:
            old_lines = f.readlines()
        new_lines = fix_lines(old_lines)
        if old_lines != new_lines:
            with open(path, 'w') as f:
                for line in new_lines:
                    f.write(line)
            print('Fixing', path)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        'files', nargs=argparse.REMAINDER, help='Paths to the files to format.'
    )
    args = parser.parse_args()
    purify_text_files(args.files)


if __name__ == '__main__':
    main()
