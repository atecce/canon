# TODO format properly and prevent two passes over the data
fd -t f | while read f; do if grep "404 Not Found" $f; then echo $f; rm $f; fi; done
fd -t f | while read f; do cat "$f" | if grep -q -E "003c.*003e.*File.*Found.*003c.*003e.*003c.*"; then echo "$f"; cat "$f"; rm "$f"; fi; done
fd -t d -x rmdir
