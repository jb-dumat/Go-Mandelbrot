------------------------------------------------------------------------
                        Technical Documentation
------------------------------------------------------------------------

To launch the programs:
    go run filename.go

To launch the benchmarks
    python3 script.py filename=FILE repetition=INT

For example:
    python3 script.py mandelbrot_worker-pool.go 50

Outputs are printed on stdout. You can redirect output to file using ">"
For example:
    python3 script.py mandelbrot_worker-pool.go 50 > OUTPUT.txt