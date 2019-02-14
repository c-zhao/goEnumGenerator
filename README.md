# goEnumGenerator
Used to generate an enum type in golang from a plain text file

Some time we have a list of strings and we want to create a enum type for it in Go. Then there will be a plenty copy paste to do the job. This tiny tool can read from a plain text file and generated a Go source code file to define a enum type. You can then add more stuff from there to avoid tedious work. 

# File format
The first line is the enum type name
And the following lines are enum values.

# Error handling
This tool is implemented as a handy tool for developers. There is not lots of error checks, please do not try to fool it. :P
