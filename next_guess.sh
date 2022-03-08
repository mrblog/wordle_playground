if [ $# -ne 4 ] ; then
  echo "usage: next_guess guess clue input-words output-words" 1>&2
  exit 2
fi

GUESS="$1"
CLUE="$2"
INWORDS="$3"
OUTWORDS="$4"

echo ./wordle_guess "$GUESS" "$CLUE" "$INWORDS" 1>&2

./wordle_guess "$GUESS" "$CLUE" "$INWORDS" 2> result.txt

grep '^Guess:' result.txt | sort -n --key=4,8 | awk 'BEGIN { lo = 9999; hlo = lo } { if ($4 <= lo && $8 <= hlo) { lo = $4; hlo = $8; print } }' > next_guesses.txt
cat next_guesses.txt

grep '^Guess:' result.txt | awk '{print $2 }' > "$OUTWORDS"
echo "$(wc -l $OUTWORDS) possibilities"
exit 0
