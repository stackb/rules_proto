import argparse
import glob
import os
import sys

parser = argparse.ArgumentParser(
    description = 'Create a node module index.js from a list of files in a directory.')

parser.add_argument('Dir',
                       metavar='srcdir',
                       type=str,
                       help='the dir to scan')
parser.add_argument('Out',
                       metavar='outfile',
                       type=str,
                       help='the filename to write')
args = parser.parse_args()

# print("dir: " + args.Dir)
# print("out: " + args.Out)
if os.path.dirname(args.Dir) != os.path.dirname(args.Out):
    print("invalid arguments: %s and %s must be in the same directory" % (args.Dir, args.Out))
    sys.exit(1)

lines = []
for r, d, f in os.walk(args.Dir):
    for file in f:
        abs = os.path.join(r, file)
        # print("abs: " + abs)
        rel = abs[len(args.Dir)+1:]
        # print("rel: " + rel)
        base = os.path.splitext(os.path.basename(rel))[0]
        lines.append(' "%s": require("./%s")' % (base, rel))

package_json = open(args.Out, "w")
package_json.write("module.exports = {\n%s\n}\n" % ",\n".join(lines))
package_json.close()

# print("WIP")
# sys.exit(1)
