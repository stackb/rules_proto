# Check generated files are as expected
import pathlib
import sys

expected = {
    'objc_lib/Demo.pbobjc.m',
    'objc_lib/Demo.pbobjc.h',
    'objc_lib/folder/NestedDemo.pbobjc.m',
    'objc_lib/folder/NestedDemo.pbobjc.h',
}

files = set([str(p) for p in pathlib.Path('objc_lib').glob('**/*') if p.is_file()])

if files != expected:
    print('Files found do not match expected')
    print('Found {}'.format(files))
    print('Expected {}'.format(expected))
    sys.exit(1)
