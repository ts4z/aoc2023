#!/usr/bin/perl -w

my $sum = 0;

for my $line (<ARGV>) {
    $line =~ /Card +(\d+):([\d+ ]+)\s+\|\s+([\d ]+)/ or die "can't parse $line";
    my @wins=grep !/^$/, split(/\s+/,$2);
    my %wins;
    map { $wins{$_} = 1 } (@wins);
    my @picks=grep !/^$/, split(/\s+/, $3);
    print "wins ", join(",", @wins), "; picks ", join(",",@picks),"\n";
    my $n = scalar(grep { $wins{$_} } (@picks));
    if ($n > 0) {
        print "match $n\n";
        $sum += 2**($n - 1);
    }
}

print "sum $sum\n";
