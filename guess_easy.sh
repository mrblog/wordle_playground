
if [ $# -lt 1 ] ; then
  echo "usage: guess match (e.g. date) [ starting-word ]" 1>&2
  exit 1
fi
MATCH="$1"
n=`cat -n ordered_answers.txt | grep -i "$MATCH" | awk '{print $1}'`
tail "+$n"  ordered_answers.txt | awk '{if (length($NF) == 5) print tolower($NF) }' > today.txt
today_word=`head -1 today.txt`
today_num=`tail +"$n"  ordered_answers.txt | head -1 | awk '{print $5}'`
echo "today_word: $today_word"
if [ $# -gt 1 ] ; then
  start_word="$2"
else
  start_index=`od -An -N1 -i /dev/random`
  start_index=`echo $start_index 25 | awk '{ slope = ($2 - 1) / 255.;  print sprintf("%.f", 1 + slope * $1) }'`
  start_word=`cat top25_starting.txt | tail +$start_index | head -1`
fi
echo "start_word: $start_word"
guess="$start_word"
## remove words that we have already guessed (small cheat)
np=`expr $n - 1`
if [ $np -gt 0 ] ; then
  head "-$np" ordered_answers.txt | awk '{if (length($NF) == 5) print tolower($NF) }' | tail | sort > previous.txt
  comm -13 previous.txt words.txt > wordset.txt
else
  cp words.txt wordset.txt
fi
echo > today_guesses.txt
guess_num=1
while [ true ] ; do
  echo "guess_num: $guess_num guess: $guess"
  guess_file="guess${guess_num}.txt"
  avg_file="avg${guess_num}.txt"
  avg_easy_file="avg_easy${guess_num}.txt"
  echo ./wordle_guess "$guess" "$today_word" wordset.txt $guess_file 1>&2
  ./wordle_guess "$guess" "$today_word" wordset.txt $guess_file 2>> today_guesses.txt
  if [ "$guess" == "$today_word" ] ; then
    echo "Wordle $today_num ${guess_num}/6"
    cat today_guesses.txt
    exit 0
  fi
  if [ $guess_num -ge 6 ] ; then
    echo "Fail"
    cat today_guesses.txt
    exit 0
  fi
  echo  ./wordle_guess $guess_file 1>&2
  n=$(cat $guess_file |wc -l)
  if [ $n -lt 200 -a $n -gt 2 ] ; then
    ./wordle_guess $guess_file 2> $avg_file
    grep '^Guess:' $avg_file | sort -n --key=4,8 | awk 'BEGIN { lo = 9999; hlo = lo } { if ($4 <= lo && $8 <= hlo) { lo = $4; hlo = $8; print } }' > next_guesses.txt
    avg=`head -1 next_guesses.txt | awk '{print $4}'`
    if (( $(echo "$avg > 1" | bc -l) ))  ; then
      ./wordle_guess -e $guess_file 2> $avg_easy_file
      grep '^Guess:' $avg_easy_file | sort -n --key=4,8 | awk 'BEGIN { lo = 9999; hlo = lo } { if ($4 <= lo && $8 <= hlo) { lo = $4; hlo = $8; print } }' > next_guesses.txt
      easy_avg=`head -1 next_guesses.txt | awk '{print $4}'`
      if (( $(echo "$easy_avg < $avg" | bc -l) )) ; then
        echo "using easy-mode guess: $easy_avg < $avg" 1>&2
        cp $avg_easy_file $avg_file
      fi
    fi
  else
    ./wordle_guess $guess_file 2> $avg_file
  fi
  grep '^Guess:' $avg_file | sort -n --key=4,8 | awk 'BEGIN { lo = 9999; hlo = lo } { if ($4 <= lo && $8 <= hlo) { lo = $4; hlo = $8; print } }' > next_guesses.txt
  n=$(cat next_guesses.txt |wc -l)
  if [ $n -gt 1 ] ; then
    guess_index=`od -An -N1 -i /dev/random`
    #echo "guess_index: $guess_index of $n $(awk '{print $2}' next_guesses.txt | fmt)" 1>&2
    guess_index=`echo $guess_index $n | awk '{ slope = ($2 - 1) / 255.;  print sprintf("%.f", 1 + slope * $1) }'`
    #echo "guess_index: $guess_index of $n" 1>&2
    guess=$(tail +$guess_index next_guesses.txt | head -1 | awk '{ print $2 }')
  else
    guess=$(head -1 next_guesses.txt | awk '{ print $2 }')
  fi
  cp $guess_file wordset.txt
  guess_num=$(expr $guess_num + 1)

done
