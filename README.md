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

**UPDATE** Updated with the new NYT answers as of Mar 30, 2022.

## Overview

This is not meant for cheating. That's no fun. 
The fun comes in comparing your Wordle solution to the algorithm, to compare
and contrast your approach vs. the simple math guesser.

The tool can also be useful for analyzing your guesses. I like to see how rapidly
guesses converge on the solution.

The basic *Go* app has three main functions, triggered by command-line args:

1. Generate a list of guess quality scores from a given input word list of possible answers (and possibly a separate guess list)
2. Check a guess against a target (answer) word and optionally produce an output set of possible scored next guesses
3. Helper-mode, given a guess and a clue, and a possible words list, produce a ranked "next guess" list

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
The default plays hard-mode. You can provide the -e option to play easy mode
using the vocabulary of over 12,000 words that Wordle accepts.

```shell
./wordle_guess -e possibilities.txt
```

It will take much longer to generate this list.

Generate a ranked set of starting words. This will take several hours to run.

```shell
$ ./wordle_guess
```

The output of this is pre-computed and included in the repo in the `best_first.txt` file.

Add the -r option to "reverse" the operation, to generate a "difficulty" value for
the answer-words, the average words remaining for all possible starting words for that answer
wrd. Outpusin this case is of the form:

```shell
Answer: youth Averge: 187.47 low: 1 high: 1324 (amity/vivid)
```

Generate a ranked list of "next guesses" given a guess, a set of possible answers, and the "clue" returned by Wordle.

```shell
$ ./wordle_guess taser 🟨⬛🟨⬛⬛ words.txt
```

The output will be the ranked possible words, as above:

```shell
Answers: 2393
Vocabulary: 161
Guess: angst Averge: 1.93 low: 0 high: 43 (aback/dusty)
Guess: artsy Averge: 2.42 low: 0 high: 59 (aback/beset)
Guess: ascot Averge: 1.90 low: 0 high: 32 (aback/baste)
...
```

This mode is used by the `next_guess.sh` helper script described below.

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

RAISE ARISE IRATE AROSE SANER ALTER
LATER STARE SNARE SLATE LASER ALERT
CRATE TRACE STALE

*UPDATE March 24, 2022:* Doing some first+second analysis 
(see https://mrblog.org/2022/03/24/more-fun-with-wordle-second-guess-analysis/)
I've updated the top-15 starting words selection to these:

SLATE LEAST TRACE REACT CRATE STALE
STARE LEANT TRAIL TRIAL RAISE ALTER
LATER SANER CRANE

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

## The helper mode

This is your cheater mode.

Start with your fist guess, for example:

```shell
$ sh next_guess.sh taser 🟨⬛🟨⬛⬛ words.txt guess1.txt
```
You can either provide the unicode for the "clue" as shown above, or letters
'g' for GREEN, 'y' for YELLOW, and 'b' for BLACK, or e.g. ybybb as the equivalent to the above.

The output will be the recommended next-guess and the output
file (`guess1.txt` in the above example) will contain the list of
possible answer-words, to use as input for the next guess, e.g.

```shell
./wordle_guess taser 🟨⬛🟨⬛⬛ words.txt
Guess: hoist Averge: 3.41 low: 1 high: 6 (ghost/stock)
      49 guess1.txt possibilities
```

The guesser suggests the next guess should be "hoist" and that there
are 49 possible answers.

You provide that to Wordle and use the "clue" returned by Wordle to produce the next guess:
```shell
$ sh next_guess.sh hoist bgggg guess1.txt guess2.txt
./wordle_guess hoist bgggg guess1.txt
Guess: foist Averge: 1.67 low: 1 high: 2 (foist/joist)
Guess: joist Averge: 1.67 low: 1 high: 2 (joist/foist)
Guess: moist Averge: 1.67 low: 1 high: 2 (moist/foist)
       3 guess2.txt possibilities
```

In this example you can see there are three possible guesses with the same 
ranking. In these cases, you simply pick one at random, or using
your favorite divining method:
```shell
$ sh next_guess.sh joist bgggg guess2.txt guess3.txt
./wordle_guess joist bgggg guess2.txt
Guess: foist Averge: 1.00 low: 1 high: 1 (foist/foist)
Guess: moist Averge: 1.00 low: 1 high: 1 (foist/foist)
       2 guess3.txt possibilities
```