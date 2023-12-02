#!/usr/bin/perl -w

use strict;

my %limits = ('red' => 12, 'green' => 13, 'blue' => 14);
my $sum = 0;

foreach my $line (<ARGV>) {
    $line =~ /Game (\d+): (.*)/ or die "can't parse $_";
    my $gn = $1;
    my $draws = $2;
    my @splits = split(/\s*;\s*/, $draws);
    my $ok = 1;
    for my $split (@splits) {
        my @descs = split(/\s*,\s/, $split);
        for my $desc (@descs) {
            $desc =~ /(\d+) (\w+)/;
            my $n = $1;
            my $color = $2;

            if ($limits{$color} < $n) {
                print STDERR "$gn $color $n over limit\n";
                $ok = undef;
            }
        }
    }
    if ($ok) {
        print STDERR "game $gn ok\n";
        $sum += $gn;
    }
}

print "$sum\n";
