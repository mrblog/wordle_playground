#!/usr/bin/env bash

rm -f oct*.txt
row=1
cp /dev/null solved.txt
nsolved=0
echo "Enter starting word:"
read guess
inf="words.txt"
unsolved=8
while [ $row -le 13 ] ; do
  pz=1
  echo "Guessing: " `echo $guess | tr [a-z] [A-Z]`
  if [ $unsolved -gt 0 ] ; then
    while [ $pz -le 8 ] ; do
      if [ ! -f oct${pz}_solved.txt ] ; then
        echo "Enter clue for puzzle ${pz} for guess ${guess}:"
        read clue
        if [ $row -gt 1 ] ; then
          prev_row=$(expr $row - 1)
          inf=oct${pz}_${prev_row}_words.txt
        fi
        ./wordle_guess $guess $clue ${inf} 2> oct${pz}_${row}.txt
        head -2 oct${pz}_${row}.txt
        grep "^Guess: " oct${pz}_${row}.txt | awk '{print $2 }' | sort > oct${pz}_${row}_words.txt
        n=`cat oct${pz}_${row}_words.txt | wc -l`
        if [ $n -eq 1 ] ; then
          answer=`cat oct${pz}_${row}_words.txt`
          if [ "$guess" != "$answer" ] ; then
            cat oct${pz}_${row}_words.txt solved.txt | sort -u > solved$$.txt
            mv solved$$.txt solved.txt
          fi
          cat oct${pz}_${row}_words.txt > oct${pz}_solved.txt
          echo "Puzzle ${pz} solved: " `cat oct${pz}_solved.txt`
          unsolved=$(expr $unsolved - 1)
        else
          ./wordle_guess -e oct${pz}_${row}_words.txt 2> oct${pz}_${row}_easy.txt
        fi
      fi
      pz=$(expr $pz + 1)
    done
    
    nsolved=$(cat solved.txt | wc -l)
    if [ $unsolved -gt 0 ] || [ $nsolved -gt 0 ] ; then
      if [ $nsolved -gt 0 ] ; then
        echo "checking $nsolved solved guesses"
        cp solved.txt next_guesses.txt
      else
        grep '^Guess:' oct?_${row}.txt | sort -n --key=4,8 | awk 'BEGIN { lo = 9999; hlo = lo } { if ($4 <= lo && $8 <= hlo) { lo = $4; hlo = $8; print $2 } }' > next_guesses.txt
        avg=`grep '^Guess:' oct?_${row}.txt | sort -n --key=4,8 | awk 'BEGIN { lo = 9999; hlo = lo } { if ($4 <= lo && $8 <= hlo) { lo = $4; hlo = $8; print }}' | head -1 | awk '{print $4}'`
        if (( $(echo "$avg > 1" | bc -l) ))  ; then
          grep '^Guess:' oct?_${row}_easy.txt | sort -n --key=4,8 | head
          grep '^Guess:' oct?_${row}_easy.txt | sort -n --key=4,8 | awk 'BEGIN { lo = 9999; hlo = lo } { if ($4 <= lo && $8 <= hlo) { lo = $4; hlo = $8; print $2 } }' > next_guesses.txt
        else
          grep '^Guess:' oct?_${row}.txt | sort -n --key=4,8 | head
        fi
        ng=$(cat next_guesses.txt | wc -l)
        echo "checking $ng next guesses"
      fi
      if [ $unsolved -gt 0 ] ; then
        cp /dev/null scored_guesses.txt
        while read next_guess  ; do
          grep "Guess: $next_guess" oct?_${row}_easy.txt | awk '{print $4}' | awk -f ~/sumstats.awk | grep mean | awk '{ print "'${next_guess}'",$NF }' >> scored_guesses.txt
        done < next_guesses.txt
        guess=`sort -n --key=2 scored_guesses.txt | head -1 | awk '{print $1}'`
        if [ $nsolved -gt 0 ] ; then
          cat solved.txt | grep -v "$guess" | sort -u > solved$$.txt
          mv solved$$.txt solved.txt
        fi
      else
        guess=`head -1 solved.txt`
        cat solved.txt | grep -v "$guess" | sort -u > solved$$.txt
        mv solved$$.txt solved.txt
      fi
    fi
  else
    nsolved=$(cat solved.txt | wc -l)
    if [ $nsolved -gt 0 ] ; then
      guess=`head -1 solved.txt`
      cat solved.txt | grep -v "$guess" | sort -u > solved$$.txt
      mv solved$$.txt solved.txt
    else
      echo "Success in ${row} tries!"
      exit 0
    fi
  fi
  row=$(expr $row + 1)
done

if [ $row -gt 13 ] ; then
  echo "Bummer"
fi

exit 0