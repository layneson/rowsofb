# Rows of B

RowsOfB is a command-line matrix calculator. It was created out of the developer's frustration with currently existing online matrix calculators. Firstly, not all of them have the desired functions. For instance, few have a simple operation to get the reduced row echelon form of a matrix; most do Gauss-Jordan elimination, but they either only reduce to row echelon form or spit out the parametric solution. Linear algebra homework can be tedious wih the amount of reductions needed, and none of the existing matrix calculation sites adequately deal with this frustration. In addition, few sites handle fractions for the duration of the operation; most spit out decimals... its 2016! That is not necessary! The interface of the web is clearly not ideal.

Inspiration for RowsOfB comes from the TI-84 calculator line. These graphing calculators have an exhaustive list of matrix operations (which include transposition, inverting, and reducing to row echelon AND reduced row echelon forms) and a simple method of entering matrix values. No one should have to press any keys other than tab to move between entries, and press enter to calculate the result.

RowsOfB is terminal-based because it is convenient to have on a computer. Many schools use online homework submission systems, so a student's computer is already open. Secondly, many problems which deal with matrices require a single matrix to be used multiple times. RowsOfB saves defined matrices, saving one time when he needs to, say, augment one and reduce it a second time.

## Supported operations

RowsOfB supports a number of matrix operations. Some are unary, and some are binary. Binary operations take either two matrices or a scalar and a matrix.

In the following documentation, `[A]` denotes the matrix named 'A' and `c` denotes a scalar.

### Unary operations
 - `inv [A]`: Calculates the inverse of 'A', or indicates that no such inverse exists.
 - `trans [A]`: Calculates the transpose of 'A'.
 - `ref [A]`: Uses Gauss-Jordan elimination to put 'A' into row echelon form.
 - `rref [A]`: Uses Gauss-Jordan elimination to put 'A' into reduced row echelon form.

### Binary operations
 - `add [A] [B]`: Adds matrix 'A' to 'B'.
 - `mul [A] [B]`: Multiplies matrix 'A' and 'B'.
 - `scl c [A]`: Multiplies matrix 'A' by the scalar 'c'.
 - `aug [A] [B]`: Augments 'A' with 'B'.

## Commands

RowsOfB operates as an interactive shell. Every operation listed above is also a command. The result of every command is stored in a hidden result slot. In addition to those operations, several commands exist:

 - `def [A]`: Opens an interactive process to define the matrix 'A'.
 - `set [A] [B]`: Sets matrix 'A' to matrix 'B'. In essence, it copies 'B' into 'A'.
 - `set [A]`: Sets matrix 'A' to the result of the last command.
 - `zero [A]`: Zeros matrix 'A'. Equivalent to `scl 0 [A]`.
 - `del [A]`: Deletes matrix 'A'. A deleted matrix has no size and no entries, and any operations using a deleted matrix raises an error.
 - `clr`: Deletes all matrices. This is the equivalent to restarting RowsOfB.