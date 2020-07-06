import os
import sys
import timeit
import math
import statistics

if len(sys.argv) != 3:
    print("usage:\n\tpython3 script.py <entry=FILE> <repetition=int>")
    sys.exit(1)

entry = sys.argv[1]
try:
    repetition = int(sys.argv[2])
except ValueError:
    print("err: The repetition paramater must be an int.")
    sys.exit(1)

sample = []
execution = "go run " + entry

def printResults(sample, description):
        try:
            print(description)
            print("\tNumber of repetition: % s" % repetition)
            print("\tAverage of sample is % s " % (statistics.mean(sample)))
            print("\tStandard Deviation of sample is % s " % (statistics.stdev(sample)))
            print("\tStandard error of the mean is % s" % (statistics.stdev(sample) / math.sqrt(repetition)))
        except statistics.StatisticsError:
            print(description)
            print("\tNumber of repetition: % s" % repetition)
            print("\nTime: % s" % sample[0])

for _ in range(0, repetition):
    start_time = timeit.default_timer()
    os.system(execution)
    sample.append(timeit.default_timer() - start_time)
printResults(sample, "Experiment using " + entry)