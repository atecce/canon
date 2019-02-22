fd -t f | while read f; do cat $f | if grep -q -E "003c.*003e.*File.*Found.*003c.*003e.*003c.*"; then echo $f; cat $f; rm $f; fi; done
fd -t d -x rmdir