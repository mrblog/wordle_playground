# Yet another Wordle Playground

Play around with the Wordle leaked word lists and guessing algorithms.

**BEWARE:** *Spoiler is here because all the Wordle answers are in this repo.
Don't look if you don't want to know.*

From Josh Wardle, the Wordle author:

> "people are saying it’s like if you’re doing the puzzle section in the newspaper, you can turn the puzzle upside down and read the answers. It’s more about like, who are you cheating if you’re doing that?"

*For research purposes only.*

**UPDATE** Updated with new NYT answers Feb 15 2022. Original Wordle answers are 
now in the file `ordered_answers_original.txt`. Likewise, the accepted guesses have been 
updated to match the new NYT set and the original list is now in the file `accepted_guesses_original.txt` 

## Overview

This is not meant for cheating. That's no fun. 
The fun comes in comparing your Wordle solution to the algorithm, to compare
and contrast your approach vs. the simple math guesser.

The tool can also be useful for analyzing your guesses. I like to see how rapidly
guesses converge on the solution.

The basic *Go* app has two main functions, triggered by command-line args:

1. Generate a list of guess quality scores from a given input word list of possible answers (and possibly a separate guess list)
2. Check a guess against a target (answer) word and optionally produce an output set of possible scored next guesses

The code is nothing too pretty. It is not forgiving. 
It does little bounds checking and will crash with bad inputs.
But it gets the job done and is hopefully not too
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

Show the Wordle match for a given guess and answer and generate a list of possible
next-guess words:

```shell
$ ./wordle_guess study moist answers.txt possibilities.txt
```

Output is the same as above, but the code in this case reads the `answers.txt` input for 
possible 5-letter word answers and generates an output of possible next-guess words, one per line,
in this case to the specified file `possibilities.txt`. 

Generate a ranked set of possible guesses given an input list of guesses / answers:

```shell
$ ./wordle_guess possibilities.txt
```

The ranked guesses are output in a form like the following

```shell
Answers: 5
Accepted Guesses: 5
Guess: moist Averge: 2.20 low: 1 high: 3
Guess: ghost Averge: 2.20 low: 1 high: 3
Guess: hoist Averge: 2.20 low: 1 high: 3
Guess: joist Averge: 2.20 low: 1 high: 3
Guess: foist Averge: 2.20 low: 1 high: 3
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

### Best starting words

The best starting words as determined by the above algorithm from the set of allowed
words can be found as follows:

```shell
$ grep Guess: best_first.txt | sort -n --key=4,8  | head -25
```

Here are some top starting words, based on this algorithm:

ROATE RAISE RAILE SOARE ARISE IRATE ORATE ARIEL AROSE RAINE ARTEL
TALER RATEL AESIR ARLES

## Try different starting words

You can force the guesser to use a specific starting word by providing one, as shown below:

```shell
$ sh guess.sh "Feb 02 2022" great
```
Output:

```shell
Wordle 228 6/6

GREAT
⬛⬛⬛⬛🟩
Possibilities: 41
SCOUT
🟨⬛🟨⬛🟩
Possibilities: 5
FOIST
⬛🟩🟩🟩🟩
Possibilities: 3
JOIST
⬛🟩🟩🟩🟩
Possibilities: 2
HOIST
⬛🟩🟩🟩🟩
Possibilities: 1
MOIST
🟩🟩🟩🟩🟩
Possibilities: 1
```

Be warned that if your starting word is weak, it could take a *long* time for the
guesser to run, or at least to pick the second guess, possibly several minutes 
depending on how terrible the supplied starting word is. It will generally recover
pretty well after the first guess and probably still solve the puzzle because it
will pick a much better second guess using the algorithm described above.
