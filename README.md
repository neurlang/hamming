# hamming
Golang Hamming(7,4) code with an additional parity bit (SECDED)

# purpose

This allows to correct single bit errors in any data, for instance when transmitting
information over a wire, that can become slightly corrupted. The cost is, that the
amount of information doubles.

# extended Hamming code

With the addition of an overall parity bit, it becomes the (8,4) extended Hamming code which is SECDED.
This package can both detect and correct single-bit errors and detect (but not correct) double-bit errors.


