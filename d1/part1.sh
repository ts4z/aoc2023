perl -ne 's/\D//g; $n = substr($_, 0, 1) . substr($_, -1, 1); $ttl += $n; END { print "$ttl\n" };'   < input
