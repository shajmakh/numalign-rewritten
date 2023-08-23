# numalign-rewritten
rewrite https://github.com/ffromani/numalign for self-experiment.

 This is a tool to check whether resources (cpus, memory & pci devices) consumed by a process are aligned to a single numa node. 

* Apply the following commands to build and run the tool:

 **`make`**

 **`./build/numalign`** will run the tool for the same process (self) created by the command and print the result on the standard output.

 **`./build/numalign -h`** will display a help output on the supported flags.

 **`./build/numalign -o myNumaAlign.log -p 15422 -v`** will run the tool and check resources alignment for process with id=15422; will print the output into the file myNumaAlign.log with additional details as v (verbose) is specified. 

* Environment variables:
DEV_RESOURCES - specify the list of pci devices names separated by "," that are used by the specific process for which you want to check its resources' alignment.