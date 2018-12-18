#!/usr/bin/env python3

# pylint: disable=missing-docstring
import unittest

import purify_text  # pylint: disable=import-error


class TestPurifyTest(unittest.TestCase):

    def test_fix_lines(self):
        data = [
            # 2-tuples of (expected processed lines, original lines)

            # empty and nearly empty files
            ([], []),
            ([], ['\n']),
            ([], ['\r']),
            ([], ['\t']),
            ([], ['\v']),
            ([], [' \n']),
            ([], [' \t']),
            ([], [' \v']),
            ([], ['\v\t']),
            ([], [' \v\t']),

            # one-line files
            (['foo\n'], ['foo\n']),
            (['  foo\n'], ['  foo\n']),
            (['\tfoo\n'], ['\tfoo\n']),
            (['\vfoo\n'], ['\vfoo\n']),

            # trailing empty lines
            (['foo\n'], ['foo\n', '\n']),
            (['foo\n'], ['foo\n', ' \n']),
            (['foo\n'], ['foo\n', '\t\n']),
            (['foo\n'], ['foo\n', '\v\n']),

            # leading empty lines
            (['foo\n'], ['\n', 'foo\n']),
            (['foo\n'], [' \n', 'foo\n']),
            (['foo\n'], ['\t\n', 'foo\n']),
            (['foo\n'], ['\v\n', 'foo\n']),
        ]
        for expected, lines in data:
            actual = purify_text.fix_lines(lines)
            self.assertEqual(expected, actual)


if __name__ == '__main__':
    unittest.main()
