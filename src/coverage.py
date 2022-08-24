import re;
import sys;


if (len(sys.argv) <= 1):
    threshold = 30
else:
    threshold = float(sys.argv[1])

print("Threshold: " + str(threshold))

#Get the last line
coverage_str = open('coverage.txt', 'r').readlines()[-1]
pattern =  r"^total\:\W+\w+\W+(?P<val>\d+\.\d)"
m = re.search(pattern, coverage_str)
cov = float(m.groups()[0])
print(cov)

if (cov < threshold):
    sys.exit("Code Coverage is too low: " + str(cov) + "%, expected 30%")

print("Code Coverage check is passed with " + str(cov) + "%, expected 30%")

