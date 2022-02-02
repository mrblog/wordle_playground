# Yet another Wordle Playground

Play around with the Wordle leaked word lists and guessing algorithms.

**BEWARE:** *Spoiler is here because all the Wordle answers are in this repo.
Don't look if you don't want to know.*

## Overview

This is not meant for cheating. That's no fun. 
The fun comes in comparing your Wordle solution to the algorithm just to compare
and contrast your approach vs. the simple math guesser.

The tool can also be useful for analyzing your guesses. I like to see how rapidly
guesses converge on the solution.

The basic *Go* app has two main functions, triggered by command-line args:

1. Generate a list of guess quality scores from a given input word list of possible answers (and possibly a separate guess list)
2. Check a guess against a target (answer) word and optionally produce an output set of possible scored next guesses

The code is nothing too pretty, but it gets the job done and is hopefully not too
difficult to work with and understand what it's doing.

## Build the code

```shell
$ go build -o wordle_guess main.go
```

## Running the code

Show the Wordle match for a given guess and answer:

```shell
$ ./wordle_guess study moist
```

Outputs:

```shell
STUDY
🟨🟨⬛⬛⬛
```

Show the Wordle match for a given guess and answer and generate a ranked next-guess list:

```shell
$ ./wordle_guess study moist answers.txt next_guesses.txt
```

Output is the same as above, but the code in this case reads the `answers.txt` input for 
possible 5-letter word answers and generates an output of ranked possible guesses,
in this case to the specified file `next_guesses.txt`. The ranked guesses are output
in a form like the following

```shell
Answers: 5
Accepted Guesses: 5
Guess: moist Averge: 2.20 low: 1 high: 3
Guess: ghost Averge: 2.20 low: 1 high: 3
Guess: hoist Averge: 2.20 low: 1 high: 3
Guess: joist Averge: 2.20 low: 1 high: 3
Guess: foist Averge: 2.20 low: 1 high: 3
```

Generate a ranked set of possible guesses given an input answer-set:

```shell
$ ./wordle_guess answers.txt
```

Generate a ranked set of starting words. This will take several hours to run.

```shell
$ ./wordle_guess
```

The output of this is pre-computed and included in the repo in the `best_first.txt` file.

## The guesser script

Forget all the above. You can simply use 
the `guess.sh` script to run the guesser on a specified Wordle of the day, 
basically to play Wordle:

```shell
$ sh guess.sh "Feb 02 2022"
```

This will produce output ultimately showing the guesses in the Wordle form:

```shell
Wordle 228 4/6

TASER
🟨⬛🟨⬛⬛
Possibilities: 44
STUDY
🟨🟨⬛⬛⬛
Possibilities: 5
HOIST
⬛🟩🟩🟩🟩
Possibilities: 3
MOIST
🟩🟩🟩🟩🟩
Possibilities: 1
```

In the output above, you can see that the first guess reduces the possibilities from
the initial 5000+ down to 44. The second guess brings it down to 5 possible words,
and so on, until the puzzle is solved or the guesser fails to solve the puzzle in 6 tries.

There is some randomization to make it more interesting so the output may be different 
with each run (there is a luck factor for the computer just like there is for humans).

## The guessing algorithm

The "best guesses" are determined based on reducing possibilities, 
using a basic brute-force average of all possible answer words, how many
next-guess possibilities does the word give for that answer. Then it assigns
the score as the average of those. When there is a tie, the preference is the 
guess with the lower maximum number of next-guess possibilities.

In the case of the starting word, it has a pre-calculated set of starting words
and picks one from random (top 25 ranked starting words).

For subsequent guesses, it ranks as described above and if there are multiple 
equivalently scored guesses, it simply picks one at random.

These random aspects of the guesser gives us the luck factor.


